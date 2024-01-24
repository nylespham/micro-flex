package main

import (
	"context"
	"log"
	"msg-logger/data"
	"time"
)

type RPCServer struct {
}

type RPCPayLoad struct {
	Name string
	Data string
}

func (s *RPCServer) LogInfo(payload RPCPayLoad, response *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println(err)
		return err
	}

	*response = "Proceed payload via RPC: " + payload.Name

	return nil
}
