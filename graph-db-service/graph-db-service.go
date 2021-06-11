package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

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

func GetRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

func main() {
	redisClient := GetRedisClient()
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
			val, getErr := redisClient.Get(ctx, string(msg.Value)).Result()
			if getErr == redis.Nil {
				//key does not exit -> fetch from arango db and write it into cache
				setErr := redisClient.Set(ctx, string(msg.Value), "", 0).Err()
				if setErr != nil {
					panic(setErr)
				}
			} else if getErr != nil {
				panic(getErr)
			} else {
				fmt.Println("key", val)
			}
		}
	}

}
