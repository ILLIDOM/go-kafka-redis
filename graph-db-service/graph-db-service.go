package main

import (
	"log"

	"github.com/Shopify/sarama"
)

func createConsumer() sarama.Consumer {
	broker := []string{"localhost:9092"}
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(broker, config)
	if err != nil {
		panic(err)
	} else {
		log.Println("New Kafka Consumer created")
	}
	return consumer
}

func main() {
	consumer := createConsumer()
	partitionConsumer, err := consumer.ConsumePartition("test", 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed Message offset %d\n", msg.Offset)
			log.Printf("Message was: %s", string(msg.Value))
			// Check if Cache has data else get it from arangodb
		}
	}

}
