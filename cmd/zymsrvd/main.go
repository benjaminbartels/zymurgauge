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
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/boltdb/bolt"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
)

func main() {

	logger := log.New(os.Stderr, "", log.LstdFlags)

	db, err := bolt.Open("zymurgaugedb", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Fatal(err)
	}
	defer safeclose.Close(db, &err)

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	beerRepo, err := database.NewBeerRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	fermentationRepo, err := database.NewFermentationRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	statikFS, err := fs.New()
	if err != nil {
		logger.Fatal(err)
	}

	routes := []app.Route{
		app.Route{Path: "chambers", Handler: handlers.NewChamberHandler(chamberRepo, pubsub.New(), logger)},
		app.Route{Path: "beers", Handler: handlers.NewBeerHandler(beerRepo, logger)},
		app.Route{Path: "fermentations", Handler: handlers.NewFermentationHandler(fermentationRepo, logger)},
	}

	app := app.New(routes, statikFS, logger)

	options := cors.Options{
		AllowedOrigins: []string{"*"},
	}

	requestLogger := middleware.NewRequestLogger(logger)

	server := http.Server{
		Addr:         ":3000",
		Handler:      app.Handler(requestLogger.Handler, cors.New(options).Handler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	var wg sync.WaitGroup
	wg.Add(1)

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
	logger.Println("Bye!")
}
