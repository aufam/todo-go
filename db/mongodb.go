package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	DB     *mongo.Database
	Client *mongo.Client
}

func MongoDBOpen(ctx context.Context, uri string) (Database, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	database := client.Database("todo")
	col := database.Collection("users")

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = col.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		client.Disconnect(ctx)
		return nil, err
	}

	return &MongoDB{database, client}, nil
}

func (d *MongoDB) Close(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}
