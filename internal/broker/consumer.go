package broker

import (
	"Messaggio/internal/storage"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type Consumer struct {
	consumer *kafka.Consumer
	storage  *storage.Storage
}

//func NewConsumer(cfg)

func ReadMessage() error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	err = c.SubscribeTopics([]string{"myTopic"}, nil)

	if err != nil {
		return err
	}

	// A signal handler or similar could be used to set this to false to break the loop.
	for true {
		var ii interface{}
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		}
		if err := json.Unmarshal(msg.Value, &ii); err != nil {
			return err
		}

	}

	c.Close()
	return nil
}
