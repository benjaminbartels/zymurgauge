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
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
	"periph.io/x/host/v3"
)

const (
	build             = "development"
	dbFilePermissions = 0o600
)

type config struct {
	Host                string        `default:":8080"`
	DebugHost           string        `default:":4000"`
	ReadTimeout         time.Duration `default:"5s"`
	WriteTimeout        time.Duration `default:"10s"`
	IdleTimeout         time.Duration `default:"120s"`
	ShutdownTimeout     time.Duration `default:"20s"`
	BrewfatherAPIUserID string        `required:"true"`
	BrewfatherAPIKey    string        `required:"true"`
	BrewfatherLogURL    string        `required:"false"`
	BleScannerTimeout   time.Duration
	Debug               bool `default:"false"`
}

func main() {
	logger := logrus.New()
	if err := run(logger); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func run(logger *logrus.Logger) error {
	var cfg config

	if err := envconfig.Process("zym", &cfg); err != nil {
		return errors.Wrap(err, "could not process env vars")
	}

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

	scanner := bluetooth.NewBLEScanner()
	monitor := tilt.NewMonitor(scanner, logger)

	go func() {
		errCh <- monitor.Run(ctx)
	}()

	db, err := bbolt.Open("zymurgaugedb", dbFilePermissions, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		return errors.Wrap(err, "could not create chamber repo")
	}

	configurator := &chamber.DefaultConfigurator{
		TiltMonitor: monitor,
	}

	var opts []brewfather.OptionsFunc
	if cfg.BrewfatherLogURL != "" {
		logger.Infof("Brewfather Log URL is set to %s", cfg.BrewfatherLogURL)
		opts = append(opts, brewfather.SetTiltURL(cfg.BrewfatherLogURL))
	}

	brewfatherService := brewfather.New(cfg.BrewfatherAPIUserID, cfg.BrewfatherAPIKey, opts...)

	chamberManager, err := chamber.NewManager(ctx, chamberRepo, configurator, brewfatherService, logger)
	if err != nil {
		logger.WithError(err).Warn("An error occurred while creating chamber manager")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	httpServer := &http.Server{
		Addr:         cfg.Host,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handlers.NewAPI(chamberManager, onewire.DefaultDevicePath, brewfatherService, shutdown, logger),
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
