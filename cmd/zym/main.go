package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/zym/brewfather"
	"github.com/benjaminbartels/zymurgauge/cmd/zym/handlers"
	c "github.com/benjaminbartels/zymurgauge/internal/platform/context"
	"github.com/benjaminbartels/zymurgauge/internal/storage"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

const dbFilePermissions = 0o600

type config struct {
	Host            string        `default:":8080"`
	ReadTimeout     time.Duration `default:"5s"`
	WriteTimeout    time.Duration `default:"10s"`
	IdleTimeout     time.Duration `default:"120s"`
	ShutdownTimeout time.Duration `default:"20s"`
	APIUserID       string        `required:"true"`
	APIKey          string        `required:"true"`
	// ThermometerID   string        `required:"true"`
	// ChillerPIN      string        `required:"true"`
	// HeaterPIN       string        `required:"true"`
	// ChillerKp       float64       `required:"true"`
	// ChillerKi       float64       `required:"true"`
	// ChillerKd       float64       `required:"true"`
	// HeaterKp        float64       `required:"true"`
	// HeaterKi        float64       `required:"true"`
	// HeaterKd        float64       `required:"true"`
	Debug bool `default:"false"`
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

	db, err := bbolt.Open("zymurgaugedb", dbFilePermissions, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}

	chamberRepo, err := storage.NewChamberRepo(db)
	if err != nil {
		return errors.Wrap(err, "could not create chamber repo")
	}

	// createFunc := CreateThermostat

	// thermostat, err := createFunc(cfg.ThermometerID, cfg.ChillerPIN, cfg.HeaterPIN, cfg.ChillerKp, cfg.ChillerKi,
	// 	cfg.ChillerKd, cfg.HeaterKp, cfg.HeaterKi, cfg.HeaterKd, logger)
	// if err != nil {
	// 	return errors.Wrap(err, "could not create thermostat")
	// }

	brewfather := brewfather.New(brewfather.APIURL, cfg.APIUserID, cfg.APIKey)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	httpServer := &http.Server{
		Addr:         cfg.Host,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      handlers.NewAPI(chamberRepo, brewfather, shutdown, logger),
	}

	httpServerErrors := make(chan error, 1)

	ctx, interruptCancel := c.WithInterrupt(context.Background())
	defer interruptCancel()

	go func() {
		logger.Infof("fermmon started, listening at %s", cfg.Host)
		httpServerErrors <- httpServer.ListenAndServe()
	}()

	return wait(ctx, httpServer, httpServerErrors, cfg.ShutdownTimeout, logger)
}

func wait(ctx context.Context, server *http.Server, serverErrors chan error, timeout time.Duration,
	logger *logrus.Logger) error {
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "fatal http server error occurred")
	case <-ctx.Done():
		logger.Info("stopping fermmon")

		ctx, timeoutCancel := context.WithTimeout(context.Background(), timeout)
		defer timeoutCancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.WithError(err).Error("could not shutdown http server")

			if err := server.Close(); err != nil {
				logger.Error(errors.Wrap(err, "could not close http server"))
			}
		}
	}

	logger.Info("fermmon stopped")

	return nil
}
