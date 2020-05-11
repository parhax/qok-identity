package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDBCollection() (*mongo.Collection, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	url := os.Getenv("MONGO_URL")
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	collection := client.Database("GoLogin").Collection("users")
	return collection, nil
}
