package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const URI = "mongodb://articles:articles_pass@localhost:27017/articles"

type MongoDBConn struct {
	Client *mongo.Client
}

func NewMongoDBConn() *MongoDBConn {
	return &MongoDBConn{}
}

func (m *MongoDBConn) Connect(ctx context.Context) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))
	if err != nil {
		return fmt.Errorf("error creating mongodb client with error: %w", err)
	}

	if err := client.Connect(ctx); err != nil {
		return fmt.Errorf("error creating mongodb client with error: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("error pinging MongoDBConn primary with error: %w", err)
	}

	m.Client = client

	return nil
}

func (m *MongoDBConn) Disconnect(ctx context.Context) {
	_ = m.Client.Disconnect(ctx)
}
