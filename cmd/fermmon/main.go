package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/brewfather"
	"github.com/benjaminbartels/zymurgauge/cmd/fermmon/handlers"
	c "github.com/benjaminbartels/zymurgauge/internal/platform/context"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type config struct {
	Host          string  `default:":8080"`
	APIUserID     string  `required:"true"`
	APIKey        string  `required:"true"`
	ThermometerID string  `required:"true"`
	ChillerPIN    string  `required:"true"`
	HeaterPIN     string  `required:"true"`
	ChillerKp     float64 `required:"true"`
	ChillerKi     float64 `required:"true"`
	ChillerKd     float64 `required:"true"`
	HeaterKp      float64 `required:"true"`
	HeaterKi      float64 `required:"true"`
	HeaterKd      float64 `required:"true"`
	Debug         bool    `default:"false"`
}

func main() {
	var cfg config

	if err := envconfig.Process("fermmon", &cfg); err != nil {
		fmt.Println("Could not process env vars:", err)
		os.Exit(1)
	}

	logger := logrus.New()

	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	bf := brewfather.New(brewfather.APIURL, cfg.APIUserID, cfg.APIKey)

	api := handlers.NewAPI(bf)

	httpServer := &http.Server{
		Addr:    cfg.Host,
		Handler: api, // TODO: add more settings
	}

	httpServerErrors := make(chan error, 1)

	go func() {
		logger.Infof("fermmon started, listening at %s", cfg.Host)
		httpServerErrors <- httpServer.ListenAndServe()
	}()

	ctx, interruptCancel := c.WithInterrupt(context.Background())
	defer interruptCancel()

	select {
	case err := <-httpServerErrors:
		logger.Error(errors.Wrap(err, "fatal http server error occurred"))
	case <-ctx.Done():
		logger.Info("Stopping fermmon")

		ctx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint:gomnd
		defer timeoutCancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Error(errors.Wrap(err, "could not shutdown http server"))

			if err := httpServer.Close(); err != nil {
				logger.Error(errors.Wrap(err, "could not close http server"))
			}
		}
	}

	logger.Info("fermmon stopped")
}
