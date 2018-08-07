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
	"github.com/rakyll/statik/fs"

	"github.com/benjaminbartels/zymurgauge/internal/platform/safeclose"
	"github.com/benjaminbartels/zymurgauge/internal/platform/web"
)

func main() {

	logger := log.New(os.Stderr, "", log.LstdFlags)

	db, err := bolt.Open("zymurgaugedb", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Fatal(err)
	}
	defer safeclose.Close(db, &err)

	beerRepo, err := database.NewBeerRepo(db)
	if err != nil {
		logger.Fatal(err)
	}

	// chamberRepo, err := database.NewChamberRepo(db)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// fermentationRepo, err := database.NewFermentationRepo(db)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	// temperatureRepo, err := database.NewTemperatureChangeRepo(db)
	// if err != nil {
	// 	logger.Fatal(err)
	// }

	beerHandler := handlers.NewBeerHandler(beerRepo)
	// chamberHandler := handlers.NewChamberHandler(chamberRepo, pubsub.New(), logger)
	// fermentationHandler := handlers.NewFermentationHandler(fermentationRepo)
	// temperatureHandler := handlers.NewTemperatureChangeHandler(temperatureRepo)

	requestLogger := middleware.NewRequestLogger(logger)
	errorHandler := middleware.NewErrorHandler(logger)

	uiFS, err := fs.New()
	if err != nil {
		logger.Fatal(err)
	}

	api := web.NewAPI("v1", logger, requestLogger.Log, errorHandler.HandleError)
	api.Register("beers", beerHandler.Handle)
	api.Register("chambers", chamberHandler.Handle)

	//app := web.NewApp(logger, uiFS, requestLogger.Log, errorHandler.HandleError)

	app := web.NewApp(api, uiFS)

	// app.Register("GET", "/api/v1/beers", beerHandler.GetAll)
	// app.Register("GET", "/api/v1/beers/:id", beerHandler.GetOne)
	// app.Register("POST", "/api/v1/beers", beerHandler.Post)
	// app.Register("DELETE", "/api/v1/beers", beerHandler.Delete)
	// app.Register("GET", "/api/v1/chambers", chamberHandler.GetAll)
	// app.Register("GET", "/api/v1/chambers/:mac", chamberHandler.GetOne)
	// app.Register("POST", "/api/v1/chambers", chamberHandler.Post)
	// app.Register("DELETE", "/api/v1/chambers", chamberHandler.Delete)
	// app.Register("GET", "/api/v1/fermentations", fermentationHandler.GetAll)
	// app.Register("GET", "/api/v1/fermentations/:id", fermentationHandler.GetOne)
	// app.Register("GET", "/api/v1/fermentations/:id/temperaturechanges", temperatureHandler.GetRange)
	// app.Register("POST", "/api/v1/fermentations/:id/temperaturechanges", temperatureHandler.Post)

	server := http.Server{
		Addr:    ":3000",
		Handler: app,
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
