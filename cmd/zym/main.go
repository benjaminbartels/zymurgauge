package main

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alecthomas/kong"
	"github.com/alexcesaro/statsd"
	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/auth"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/device/onewire"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth/ibeacon"
	"github.com/benjaminbartels/zymurgauge/internal/platform/debug"
	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/benjaminbartels/zymurgauge/ui"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"golang.org/x/crypto/bcrypt"
	"periph.io/x/host/v3"
)

const (
	build             = "development"
	dbFilePermissions = 0o600
	bboltReadTimeout  = 1 * time.Second
)

type config struct {
	Host                   string        `default:":8080"`
	DebugHost              string        `default:":4000"`
	DBPath                 string        `default:"data/zymurgaugedb"`
	ReadTimeout            time.Duration `default:"5s"`
	WriteTimeout           time.Duration `default:"10s"`
	IdleTimeout            time.Duration `default:"120s"`
	ShutdownTimeout        time.Duration `default:"20s"`
	ReadingsUpdateInterval time.Duration `default:"1m"`
	Debug                  bool          `default:"false"`
}

type cli struct {
	Run  struct{} `kong:"cmd,help:'Run zymurgauge service.'"`
	Init struct {
		AdminUsername string `kong:"arg,help:'Admin username.'"`
		AdminPassword string `kong:"arg,help:'Admin password.'"`
	} `kong:"cmd,help:'Initialize admin credentials.'"`
}

func main() {
	logger := logrus.New()

	var cfg config

	if err := envconfig.Process("zym", &cfg); err != nil {
		logger.WithError(err).Error("could not process env vars")
	}

	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	cli := cli{}
	ctx := kong.Parse(&cli,
		kong.Name("zym"),
		kong.Description("Zymurgauge Brewery Manager"),
		kong.UsageOnError(),
	)

	switch ctx.Command() {
	case "run":
		if err := run(logger, cfg); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	case "init <admin-username> <admin-password>":
		_, settingsRepo, err := createRepos(cfg.DBPath)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		if err := checkAndInitSettings(cli.Init.AdminUsername, cli.Init.AdminPassword, settingsRepo,
			logger); err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	}
}

func run(logger *logrus.Logger, cfg config) error {
	if _, err := host.Init(); err != nil {
		return errors.Wrap(err, "could not initialize gpio")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	chamberRepo, settingsRepo, err := createRepos(cfg.DBPath)
	if err != nil {
		return errors.Wrap(err, "could not create databases")
	}

	// TODO: Handle settings update without restart
	s, err := settingsRepo.Get()
	if err != nil {
		logger.WithError(err).Warn("could not get settings")
	}

	if s == nil {
		return errors.New("Settings are not initialized. Please run 'zym init'")
	}

	errCh := make(chan error, 1)

	monitor, err := createTiltMonitor(ctx, logger, errCh)
	if err != nil {
		return errors.Wrap(err, "could not create tilt monitor")
	}

	startDebugEndpoint(cfg.DebugHost, logger)

	configurator := &chamber.DefaultConfigurator{
		TiltMonitor: monitor,
	}

	var statsdClient *statsd.Client

	if s.StatsDAddress != "" {
		statsdClient, err = statsd.New(statsd.Address(s.StatsDAddress))
		if err != nil {
			logger.WithError(err).Warn("Could not create statsd client.")
		}
	}

	brewfatherClient := brewfather.New(s.BrewfatherAPIUserID, s.BrewfatherAPIKey, s.BrewfatherLogURL)

	chamberManager, err := chamber.NewManager(ctx, chamberRepo, configurator, brewfatherClient, logger, statsdClient,
		cfg.ReadingsUpdateInterval)
	if err != nil {
		logger.WithError(err).Warn("An error occurred while creating chamber manager")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Rename DefaultDevicePath and make configurable based on OS

	settingsCh := startUpdateSettingsChannel(brewfatherClient)

	app, err := handlers.NewApp(chamberManager, onewire.DefaultDevicePath, brewfatherClient, settingsRepo, settingsCh,
		ui.FS, shutdown, logger)
	if err != nil {
		return errors.Wrap(err, "could not create new app")
	}

	httpServer := &http.Server{
		Addr:         cfg.Host,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      app,
	}

	go func() {
		logger.Infof("zymurgauge version %s started, listening at %s", build, cfg.Host)
		errCh <- httpServer.ListenAndServe()
	}()

	return wait(ctx, httpServer, errCh, cfg.ShutdownTimeout, logger)
}

func createTiltMonitor(ctx context.Context, logger *logrus.Logger, errCh chan error) (*tilt.Monitor, error) {
	discoverer, err := ibeacon.NewDiscoverer(logger)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new ibeacon discoverer")
	}

	monitor := tilt.NewMonitor(discoverer, logger)

	go func() {
		errCh <- monitor.Run(ctx)
	}()

	return monitor, nil
}

func startDebugEndpoint(host string, logger *logrus.Logger) {
	go func() {
		if err := http.ListenAndServe(host, debug.Mux()); err != nil {
			logger.WithError(err).Errorf("Debug endpoint %s closed.", host)
		}
	}()
}

func createRepos(path string) (chamberRepo *database.ChamberRepo, settingsRepo *database.SettingsRepo, err error) {
	db, err := bbolt.Open(path, dbFilePermissions, &bbolt.Options{Timeout: bboltReadTimeout})
	if err != nil {
		err = errors.Wrap(err, "could not open database")

		return
	}

	chamberRepo, err = database.NewChamberRepo(db)
	if err != nil {
		err = errors.Wrap(err, "could not create chamber repo")

		return
	}

	settingsRepo, err = database.NewSettingsRepo(db)
	if err != nil {
		err = errors.Wrap(err, "could not create settings repo")

		return
	}

	return
}

func checkAndInitSettings(username, password string, settingsRepo *database.SettingsRepo, logger *logrus.Logger) error {
	user, err := settingsRepo.Get()
	if err != nil {
		return errors.Wrap(err, "could not get settings")
	}

	// settings exist
	if user != nil {
		if username != "" && password != "" {
			logger.Warn("Admin credentials have already been set.")
		}

		return nil
	}

	// settings do not exist

	if username == "" || password == "" {
		return errors.New("initial credentials have not been set, please set them by running 'zym init'")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "could not generate password hash")
	}

	letters := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	secretKeyLength := 64

	b := make([]byte, secretKeyLength)

	for i := 0; i < secretKeyLength; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}

	s := &settings.Settings{
		AppSettings: settings.AppSettings{
			AuthSecret:       string(b),
			TemperatureUnits: "Celsius",
		},
		Credentials: auth.Credentials{
			Username: username,
			Password: string(hash),
		},
		ModTime: time.Now(),
	}

	if err := settingsRepo.Save(s); err != nil {
		return errors.Wrap(err, "could not save initial settings")
	}

	return nil
}

func startUpdateSettingsChannel(brewfatherClient *brewfather.ServiceClient) chan settings.Settings {
	settingsCh := make(chan settings.Settings)

	go func() {
		for {
			update := <-settingsCh
			brewfatherClient.UpdateSettings(update.BrewfatherAPIUserID, update.BrewfatherAPIKey, update.BrewfatherLogURL)
		}
	}()

	return settingsCh
}

func wait(ctx context.Context, server *http.Server, errCh chan error, timeout time.Duration,
	logger *logrus.Logger,
) error {
	select {
	case err := <-errCh:
		return errors.Wrap(err, "fatal error occurred")
	case <-ctx.Done():
		logger.Info("stopping zymurgauge")

		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
		defer timeoutCancel()

		//nolint: contextcheck // https://github.com/sylvia7788/contextcheck/issues/2
		if err := server.Shutdown(timeoutCtx); err != nil {
			logger.WithError(err).Error("Could not shutdown http server.")

			if err := server.Close(); err != nil {
				logger.Error(errors.Wrap(err, "could not close http server"))
			}
		}
	}

	logger.Info("zymurgauge stopped 👋!")

	return nil
}
