package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const port = "8300"

type Config struct {
	Rabbit amqp.Connection
}

func main() {
	// connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	app := Config{
		Rabbit: *rabbitConn,
	}

	log.Printf("Starting server on port %s\n", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}
	err = server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://admin:P@ssw0rd@rabbitmq:5672")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(backOff), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
	}
	log.Println("RabbitMQ is connected")
	return connection, nil
}
