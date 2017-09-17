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

	app := &handlers.App{
		API: &handlers.API{
			ChamberHandler:      handlers.NewChamberHandler(chamberRepo, pubsub.New()),
			BeerHandler:         handlers.NewBeerHandler(beerRepo),
			FermentationHandler: handlers.NewFermentationHandler(fermentationRepo),
		},
		Web: http.StripPrefix("/", http.FileServer(statikFS)),
	}

	options := cors.Options{
		AllowedOrigins: []string{"*"},
	}

	c := cors.New(options)

	corsHandler := c.Handler(app)

	server := http.Server{
		Addr:    ":3000",
		Handler: corsHandler,
	}

	fmt.Println("Listening.....", server.Addr)
	panic(server.ListenAndServe())

}
