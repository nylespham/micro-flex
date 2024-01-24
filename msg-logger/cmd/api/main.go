package main

import (
	"context"
	"fmt"
	"log"
	"msg-logger/data"
	"net"
	"net/http"
	"net/rpc"
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
	Models data.Models
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

	app := Config{
		Models: data.New(client),
	}
	// Register RPC server

	err = rpc.Register(new(RPCServer))

	go app.rpcListen()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port: ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options

	clientOptions := options.Client().ApplyURI(mongoURL)

	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "P@ssw0rd",
	})
	// connect to mongo
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error connecting to mongo")
		return nil, err
	}
	log.Println("Connected to mongo")
	return client, nil
}
