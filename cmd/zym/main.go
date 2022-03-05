package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	"github.com/benjaminbartels/zymurgauge/internal/brewfather"
	"github.com/benjaminbartels/zymurgauge/internal/chamber"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/device/onewire"
	"github.com/benjaminbartels/zymurgauge/internal/device/tilt"
	"github.com/benjaminbartels/zymurgauge/internal/platform/bluetooth"
	c "github.com/benjaminbartels/zymurgauge/internal/platform/context"
	"github.com/benjaminbartels/zymurgauge/internal/settings"
	"github.com/benjaminbartels/zymurgauge/web"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"periph.io/x/host/v3"
)

const (
	build             = "development"
	dbFilePermissions = 0o600
	bboltReadTimeout  = 1 * time.Second
)

type config struct {
	Host            string        `default:":8080"`
	DebugHost       string        `default:":4000"`
	DBPath          string        `default:"zymurgaugedb"`
	ReadTimeout     time.Duration `default:"5s"`
	WriteTimeout    time.Duration `default:"10s"`
	IdleTimeout     time.Duration `default:"120s"`
	ShutdownTimeout time.Duration `default:"20s"`
	Debug           bool          `default:"false"`
}

func main() {
	logger := logrus.New()

	var cfg config

	if err := envconfig.Process("zym", &cfg); err != nil {
		logger.WithError(err).Error("could not process env vars")
	}

	if err := run(logger, cfg); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

//nolint:funlen,cyclop // TODO: Revisit later
func run(logger *logrus.Logger, cfg config) error {
	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	if _, err := host.Init(); err != nil {
		return errors.Wrap(err, "could not initialize gpio")
	}

	ctx, interruptCancel := c.WithInterrupt(context.Background())
	defer interruptCancel()

	errCh := make(chan error, 1)

	go func() {
		if err := http.ListenAndServe(cfg.DebugHost, handlers.DebugMux()); err != nil {
			logger.WithError(err).Errorf("Debug endpoint %s closed.", cfg.DebugHost)
		}
	}()

	monitor := tilt.NewMonitor(bluetooth.NewBLEScanner(), logger)

	go func() {
		errCh <- monitor.Run(ctx)
	}()

	db, err := bbolt.Open(cfg.DBPath, dbFilePermissions, &bbolt.Options{Timeout: bboltReadTimeout})
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		return errors.Wrap(err, "could not create chamber repo")
	}

	settingsRepo, err := database.NewSettingsRepo(db)
	if err != nil {
		return errors.Wrap(err, "could not create settings repo")
	}

	configurator := &chamber.DefaultConfigurator{
		TiltMonitor: monitor,
	}

	var logToBrewfather bool

	// TODO: Handle settings update without restart
	s, err := settingsRepo.Get()
	if err != nil {
		logger.WithError(err).Warn("could not get settings")
	}

	brewfatherClient := brewfather.New(s.BrewfatherAPIUserID, s.BrewfatherAPIKey, s.BrewfatherLogURL)

	if s.BrewfatherLogURL != "" {
		logToBrewfather = true
	}

	chamberManager, err := chamber.NewManager(ctx, chamberRepo, configurator, brewfatherClient, logToBrewfather, logger)
	if err != nil {
		logger.WithError(err).Warn("An error occurred while creating chamber manager")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// TODO: Rename DefaultDevicePath and make configurable based on OS

	settingsCh := make(chan settings.Settings)

	go func() {
		for {
			update := <-settingsCh
			brewfatherClient.UpdateSettings(update.BrewfatherAPIUserID, update.BrewfatherAPIKey, update.BrewfatherLogURL)

			var logToBrewfather bool // TODO: get rid of logToBrewfather

			if update.BrewfatherLogURL != "" {
				logToBrewfather = true
			}

			chamberManager.SetLogToBrewfather(logToBrewfather)
		}
	}()

	httpServer := &http.Server{
		Addr:         cfg.Host,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler: handlers.NewAPI(chamberManager, onewire.DefaultDevicePath, brewfatherClient, web.FS,
			settingsRepo, settingsCh, shutdown, logger),
	}

	go func() {
		logger.Infof("zymurgauge version %s started, listening at %s", build, cfg.Host)
		errCh <- httpServer.ListenAndServe()
	}()

	return wait(ctx, httpServer, errCh, cfg.ShutdownTimeout, logger)
}

func wait(ctx context.Context, server *http.Server, errCh chan error, timeout time.Duration,
	logger *logrus.Logger) error {
	select {
	case err := <-errCh:
		return errors.Wrap(err, "fatal error occurred")
	case <-ctx.Done():
		logger.Info("stopping zymurgauge")

		ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
		defer timeoutCancel()

		//nolint: contextcheck // https://github.com/sylvia7788/contextcheck/issues/2
		if err := server.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("Could not shutdown http server.")

			if err := server.Close(); err != nil {
				logger.Error(errors.Wrap(err, "could not close http server"))
			}
		}
	}

	logger.Info("zymurgauge stopped 👋!")

	return nil
}
