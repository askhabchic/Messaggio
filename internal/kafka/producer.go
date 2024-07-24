//package kafka
//
//import (
//	"Messaggio/internal/models"
//	"encoding/json"
//	"fmt"
//	"github.com/confluentinc/confluent-kafka-go/kafka"
//	"strconv"
//)
//
//type Producer struct {
//	producer *kafka.Producer
//}
//
//func NewProducer() (*Producer, error) {
//
//	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
//	if err != nil {
//		panic(err)
//	}
//
//	defer p.Close()
//
//	// Delivery report handler for produced messages
//	go func() {
//		for e := range p.Events() {
//			switch ev := e.(type) {
//			case *kafka.Message:
//				if ev.TopicPartition.Error != nil {
//					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
//				} else {
//					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
//				}
//			}
//		}
//	}()
//
//	// Produce messages to topic (asynchronously)
//	topic := "myTopic"
//	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
//		p.Produce(&kafka.Message{
//			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
//			Value:          []byte(word),
//		}, nil)
//	}
//
//	// Wait for message deliveries before shutting down
//	p.Flush(15 * 1000)
//}
//
//func (p *Producer) SendMessage(msg *models.Message) error {
//	const fn = "kafka.producer.SendMessage()"
//
//	msgBytes, err := json.Marshal(msg)
//	if err != nil {
//		return fmt.Errorf("%s: %w", fn, err)
//	}
//
//	kafkaMsg := kafka.Message{
//		Value:          msgBytes,
//		Key:            []byte(strconv.Itoa(msg.ID)),
//	}
//	return kafka.Producer{}.
//}
