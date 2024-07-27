package broker

import (
	"Messaggio/internal/models"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"strconv"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer() (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	return &Producer{producer: p}, nil
}

func (p *Producer) SendMessage(msg models.Message) error {
	topic := "topicMessage"

	if err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(msg.Content),
		Key:            []byte(strconv.Itoa(msg.ID)),
	}, nil); err != nil {
		log.Fatal(err)
		return err
	}

	p.producer.Flush(15 * 1000)
	return nil
}

func (p *Producer) Close() {
	p.producer.Close()
}
