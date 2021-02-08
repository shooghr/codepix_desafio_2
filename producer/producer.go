package producer

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

func InitProducer() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	err := godotenv.Load(basepath + "/../.env")
	if err != nil {
		log.Fatalf("Error loading .env files")
	}

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
	}

	producer, err := ckafka.NewProducer(configMap)
	if err != nil {
		panic(err)
	}

	topic := os.Getenv("kafkaTopic")
	msg := "Message chellange-two"
	deliveryChan := make(chan ckafka.Event)

	message := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{Topic: &topic, Partition: ckafka.PartitionAny},
		Value:          []byte(msg),
	}

	err = producer.Produce(message, deliveryChan)

	if err != nil {
		panic(err)
	}

	for e := range deliveryChan {
		switch ev := e.(type) {
		case *ckafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Delivery failed:", ev.TopicPartition)
			} else {
				fmt.Println("Delivered message:", ev.TopicPartition)
			}
		}
	}
}
