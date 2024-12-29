package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rickj1ang/RRS/internal/data"
	"github.com/rickj1ang/RRS/internal/jsonlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MONGO_URI"), "MongoDB DSN")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	client, err := openMongoClient(cfg)
	if err != nil {
		// we can not use app.logger.error Now, before we declare it
		// it not beatiful but it is what it is
		logger.PrintFatal(err, nil)
	}

	defer client.Disconnect(context.TODO())

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(client),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
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
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{primitive.E{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}
