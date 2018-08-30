package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/zymsrvd/handlers"
	_ "github.com/benjaminbartels/zymurgauge/cmd/zymsrvd/statik"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/boltdb/bolt"
	"github.com/kelseyhightower/envconfig"
	"github.com/rakyll/statik/fs"

	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

type config struct {
	HostAddress  string        `default:"0.0.0.0:3000"`
	ReadTimeout  time.Duration `default:"5s"`
	WriteTimeout time.Duration `default:"5s"`
	AuthSecret   string        `required:"true"`
}

func main() {

	logger := log.New(os.Stderr, "", log.LstdFlags)

	// Process env variables
	var appCfg config
	err := envconfig.Process("fermmond", &appCfg)
	if err != nil {
		logger.Fatal(err.Error())
	}

	db, err := bolt.Open("zymurgaugedb", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Fatal(err)
	}
	defer safeclose.Close(db, &err)

	beerRepo, err := database.NewBeerRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	fermentationRepo, err := database.NewFermentationRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	temperatureChangeRepo, err := database.NewTemperatureChangeRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	beerHandler := handlers.NewBeerHandler(beerRepo)
	chamberHandler := handlers.NewChamberHandler(chamberRepo, pubsub.New(), logger)
	fermentationHandler := handlers.NewFermentationHandler(fermentationRepo, temperatureChangeRepo)

	requestLogger := middleware.NewRequestLogger(logger)
	errorHandler := middleware.NewErrorHandler(logger)
	authorizer := middleware.NewAuthorizer(appCfg.AuthSecret, logger)

	uiFS, err := fs.New()
	if err != nil {
		logger.Fatal(err)
	}

	api := web.NewAPI("v1", logger, requestLogger.Log, errorHandler.HandleError, authorizer.Authorize)
	api.Register("beers", beerHandler.Handle, true)
	api.Register("chambers", chamberHandler.Handle, true)
	api.Register("fermentations", fermentationHandler.Handle, true)

	app := web.NewApp(api, uiFS, logger)

	startServer(app, appCfg, logger)

	logger.Println("Bye!")
}

func startServer(handler http.Handler, cfg config, logger *log.Logger) {

	var wg sync.WaitGroup
	wg.Add(1)

	server := http.Server{
		Addr:         cfg.HostAddress,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      handler,
	}

	go func() {
		logger.Printf("Listening at %s", server.Addr)
		logger.Printf("Listener at %s closed: %v", server.Addr, server.ListenAndServe())
		wg.Done()
	}()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Printf("Graceful shutdown did not complete in %v : %v", timeout, err)
		if err := server.Close(); err != nil {
			logger.Printf("Error killing server : %v", err)
		}
	}

	wg.Wait()
}
