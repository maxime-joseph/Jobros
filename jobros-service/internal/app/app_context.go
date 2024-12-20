package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

type AppContext struct {
	Config      AppConfig
	MongoClient *mongo.Client
	Database    *mongo.Database
}

// NewAppContext initializes and returns a new AppContext.
func NewAppContext() (*AppContext, error) {
	var config AppConfig

	// Load config from file if provided, otherwise use env vars
	configPath := os.Getenv("CONFIG_PATH")
	if configPath != "" {
		configFile, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	}

	// Process environment variables to override config file or defaults
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, fmt.Errorf("failed to process env vars: %w", err)
	}

	// Initialize glog with the log level from config
	flag.Set("v", fmt.Sprintf("%d", config.Logging.Level))
	flag.Parse()
	defer glog.Flush()

	glog.V(2).Info("Starting application with config:", config)

	mongoClient, database, err := initMongo(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize mongo: %w", err)
	}

	appCtx := &AppContext{
		Config:      config,
		MongoClient: mongoClient,
		Database:    database,
	}

	return appCtx, nil
}

// initMongo initializes the MongoDB client and database
func initMongo(config AppConfig) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.URI))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to ping mongo: %w", err)
	}

	database := client.Database(config.Mongo.Database)
	log.Println("Connected to MongoDB!")
	glog.V(2).Info("Connected to MongoDB!")
	return client, database, nil
}
