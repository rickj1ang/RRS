package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	env  string
}

type application struct {
	config      config
	loggerInfo  *log.Logger
	loggerError *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	loggerInfo := log.New(os.Stdout, "RRS::info: ", log.Ldate|log.Ltime)
	loggerError := log.New(os.Stdout, "RRS::error: ", log.Ldate|log.Ltime)

	app := &application{
		config:      cfg,
		loggerInfo:  loggerInfo,
		loggerError: loggerError,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	app.loggerInfo.Printf("start serving %s server on %d port", app.config.env, app.config.port)
	err := srv.ListenAndServe()
	if err != nil {
		app.loggerError.Fatal(err)
	}

}
