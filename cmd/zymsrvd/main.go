package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/benjaminbartels/zymurgauge/cmd/zymsrvd/handlers"
	_ "github.com/benjaminbartels/zymurgauge/cmd/zymsrvd/statik"
	"github.com/benjaminbartels/zymurgauge/internal/database"
	"github.com/benjaminbartels/zymurgauge/internal/middleware"
	"github.com/benjaminbartels/zymurgauge/internal/platform/app"
	"github.com/benjaminbartels/zymurgauge/internal/platform/pubsub"
	"github.com/boltdb/bolt"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
)

func main() {

	// Setup graceful exit
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		os.Exit(1)
	}()

	db, err := bolt.Open("zymurgaugedb", 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err) //ToDo: Implement logger
	}
	defer db.Close()

	chamberRepo, err := database.NewChamberRepo(db)
	if err != nil {
		panic(err) //ToDo: Implement logger
	}

	beerRepo, err := database.NewBeerRepo(db)
	if err != nil {
		panic(err) //ToDo: Implement logger
	}

	fermentationRepo, err := database.NewFermentationRepo(db)
	if err != nil {
		panic(err) //ToDo: Implement logger
	}

	statikFS, err := fs.New()
	if err != nil {
		panic(err) //ToDo: Implement logger
	}

	routes := []app.Route{
		app.Route{Path: "chambers", Handler: handlers.NewChamberHandler(chamberRepo, pubsub.New())},
		app.Route{Path: "beers", Handler: handlers.NewBeerHandler(beerRepo)},
		app.Route{Path: "fermentations", Handler: handlers.NewFermentationHandler(fermentationRepo)},
	}

	app := app.New(routes, statikFS)

	options := cors.Options{
		AllowedOrigins: []string{"*"},
	}

	server := http.Server{
		Addr:    ":3000",
		Handler: app.Handler(middleware.RequestLogger, cors.New(options).Handler),
	}

	fmt.Println("Listening.....", server.Addr)
	panic(server.ListenAndServe())

}
