package testutils

import (
	"context"
	"github.com/maxime-joseph/Jobros/jobros-service/internal/app"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

// NewMockAppContext creates a new AppContext instance with mock MongoDB components for testing
func NewMockAppContext(t *testing.T) (*app.AppContext, error) {
	mockClient, err := NewMockMongoClient(t)
	if err != nil {
		return nil, err
	}

	mockDB := NewMockDatabase(mockClient, "test_db")

	appCtx := &app.AppContext{
		Config: app.AppConfig{
			Mongo: app.MongoConfig{
				URI:      "mongodb://mock:27017",
				Database: "test_db",
			},
			Logging: app.LoggingConfig{
				Level: 2,
			},
		},
		MongoClient: mockClient,
		Database:    mockDB,
	}

	return appCtx, nil
}

// NewMockMongoClient creates a new mock MongoDB client for testing
func NewMockMongoClient(t *testing.T) (*mongo.Client, error) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	return mt.Client, nil
}

// NewMockDatabase creates a mock database instance
func NewMockDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

// CleanupMockClient handles cleanup of mock client resources
func CleanupMockClient(ctx context.Context, client *mongo.Client) error {
	return client.Disconnect(ctx)
}
