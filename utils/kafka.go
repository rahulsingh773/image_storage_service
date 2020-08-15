package utils

import (
	"fmt"
	"image_server/config"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var producer *kafka.Producer
var topic = config.Config["kafka_topic"]

func init() {
	fmt.Printf("creating kafka producer\n")

	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.Config["kafka_host"]})
	if err != nil {
		fmt.Printf("error while creating producer: %v\n", err)
		panic(err)
	}

	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to %v\n", ev)
				}
			}
		}
	}()
}

func PublishKafkaEvent(event string) error {
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(event),
	}, nil)

	return nil
}

func Close() {
	producer.Close()
}
