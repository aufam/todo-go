package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DB *mongo.Database
}

func MongoDBOpen(uri string) (Database, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return MongoDB{client.Database("todo-go")}, nil
}

func (d MongoDB) Close() error {
	return d.DB.Client().Disconnect(context.Background())
}

func getIDFromHex(id string) {
}
