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
	c "github.com/benjaminbartels/zymurgauge/internal/platform/context"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

const (
	build             = "development"
	dbFilePermissions = 0o600
)

type config struct {
	Host            string        `default:":8080"`
	DebugHost       string        `default:":4000"`
	ReadTimeout     time.Duration `default:"5s"`
	WriteTimeout    time.Duration `default:"10s"`
	IdleTimeout     time.Duration `default:"120s"`
	ShutdownTimeout time.Duration `default:"20s"`
	APIUserID       string        `required:"true"`
	APIKey          string        `required:"true"`
	Debug           bool          `default:"false"`
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

	go func() {
		if err := http.ListenAndServe(cfg.DebugHost, handlers.DebugMux()); err != nil {
			logger.WithError(err).Errorf("Debug endpoint %s closed.", cfg.DebugHost)
		}
	}()

	db, err := bbolt.Open("zymurgaugedb", dbFilePermissions, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		return errors.Wrap(err, "could not create chamber repo")
	}

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
		logger.Infof("fermmon started version %s, listening at %s", build, cfg.Host)
		httpServerErrors <- httpServer.ListenAndServe()
	}()

	if err := startThermostatTest(chamberRepo, logger); err != nil {
		return errors.Wrap(err, "could not start start pid test")
	}

	return wait(ctx, httpServer, httpServerErrors, cfg.ShutdownTimeout, logger)
}

func startThermostatTest(chamberRepo chamber.Repo, logger *logrus.Logger) error {
	chambers, err := chamberRepo.GetAll()
	if err != nil {
		return errors.Wrap(err, "could not get all chambers")
	}

	if len(chambers) > 0 {
		chamber := chambers[0]

		logger.Infof("Using Chamber %s", chamber.ID)

		createFunc := CreateThermostat

		pid, err := createFunc(chamber.ThermometerAddress, chamber.ChillerPIN, chamber.HeaterPIN, chamber.ChillerKp,
			chamber.ChillerKi, chamber.ChillerKd, chamber.HeaterKp, chamber.HeaterKi, chamber.HeaterKd, logger)
		if err != nil {
			return errors.Wrap(err, "could not create pid")
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			if err := pid.Run(ctx, 55); err != nil { //nolint:gomnd
				logger.Error(errors.Wrap(err, "error occurred with pid temperature contoller"))
			}
		}()

		<-time.After(10 * time.Second) //nolint:gomnd
		cancel()
	} else {
		logger.Info("No chambers found")
	}

	return nil
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
			logger.WithError(err).Error("Could not shutdown http server.")

			if err := server.Close(); err != nil {
				logger.Error(errors.Wrap(err, "could not close http server"))
			}
		}
	}

	logger.Info("fermmon stopped")

	return nil
}
