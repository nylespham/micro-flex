package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = "12800"
	rpcPort  = "11200"
	mongoURL = "mongodb://localhost:27017"
	grpcPort = "13500"
)

var client *mongo.Client

type Config struct {
}

func main() {
	//connect to mongodb
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	// disconnect from mongo
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options

	clientOptions := options.Client().ApplyURI(mongoURL)

	clientOptions.SetAuth(options.Credential{
		Username: "root",
		Password: "root",
	})
	// connect to mongo
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error connecting to mongo")
		return nil, err
	}

	return client, nil
}
