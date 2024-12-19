package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rickj1ang/RRS/internal/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type cusLog struct {
	info  *log.Logger
	error *log.Logger
}

type application struct {
	config config
	logger cusLog
	models data.Models
}

func main() {
	var cfg config
	var myLog cusLog

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MONGO_URI"), "MongoDB DSN")
	flag.Parse()

	loggerInfo := log.New(os.Stdout, "RRS::info: ", log.Ldate|log.Ltime)
	loggerError := log.New(os.Stdout, "RRS::error: ", log.Ldate|log.Ltime)
	myLog.info = loggerInfo
	myLog.error = loggerError

	client, err := openMongoClient(cfg)
	if err != nil {
		// we can not use app.logger.error Now, before we declare it
		// it not beatiful but it is what it is
		myLog.error.Fatal(err)
	}

	defer client.Disconnect(context.TODO())

	app := &application{
		config: cfg,
		logger: myLog,
		models: data.NewModels(client),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	app.logger.info.Printf("start serving %s server on %d port", app.config.env, app.config.port)
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.error.Fatal(err)
	}

}

func openMongoClient(cfg config) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(cfg.db.dsn).SetServerAPIOptions(serverAPI)
	// poolsize set to 10 MongoDB driver's default is 100
	opts.SetMaxPoolSize(10)
	opts.SetMaxConnIdleTime(time.Second * 600)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}
