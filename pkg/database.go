package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DatabaseURI = "mongodb+srv://Aphrodi:2000@cluster0.vkgghvi.mongodb.net/?retryWrites=true&w=majority"

var Client *mongo.Client

func Connect() error {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(DatabaseURI).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}

	Client = client
	return nil
}
