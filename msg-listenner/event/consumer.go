package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type PayLoad struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()

	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()

	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consummer *Consumer) Listen(topics []string) error {
	channel, err := consummer.conn.Channel()

	if err != nil {
		return err
	}

	defer channel.Close()

	queue, err := declareRandomQueue(channel)

	if err != nil {
		return err
	}

	for _, s := range topics {
		channel.QueueBind(
			queue.Name,
			s,
			"logs_topics",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			var payLoad PayLoad
			_ = json.Unmarshal(d.Body, &payLoad)
			go handlePayLoad(payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topics, %s]\n", queue.Name)
	<-forever
	return nil

}

func handlePayLoad(payload PayLoad) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		err := authEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}

}

func logEvent(entry PayLoad) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceUrl := "http://localhost:12800/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
