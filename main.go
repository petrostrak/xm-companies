package main

import (
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	var app Application

	err := app.New()
	if err != nil {
		os.Exit(1)
	}

	srv := &http.Server{
		Addr:        app.Config.ServerAddress,
		Handler:     app.Routes,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
	}

	app.Logger.Printf("starting development server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		app.Logger.Println(err)
		os.Exit(1)
		// gracefully shutdown.
	}
}
