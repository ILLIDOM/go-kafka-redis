package main

import (
	"context"
	"log"
	"net"

	pb "example.com/prototype/proto"
	"github.com/Shopify/sarama"
	"google.golang.org/grpc"
)

const (
	GRPC_PORT = ":50500"
)

type server struct {
	pb.UnimplementedInterfaceServiceServer
}

func CreateProducer() sarama.AsyncProducer {
	var brokers = []string{"127.0.0.1:9092"}
	config := sarama.NewConfig()
	producer, err := sarama.NewAsyncProducer(brokers, config)

	if err != nil {
		panic(err)
	} else {
		log.Println("Kafka AsyncProducer up and running!")
	}

	return producer
}

func produceMessage(producer sarama.AsyncProducer, message string) {
	kafkaMessage := &sarama.ProducerMessage{Topic: "test", Value: sarama.StringEncoder(message)}
	select {
	case producer.Input() <- kafkaMessage:
		log.Printf("Message produced")
	case err := <-producer.Errors():
		log.Fatalf("Failed to produce message: %v", err)
	}
}

func (s *server) GetReceivedBytes(ctx context.Context, in *pb.InterfaceRequest) (*pb.InterfaceReply, error) {
	log.Printf("GRPC: Received Interface Name : %v and RouterName %v", in.InterfaceName, in.RouterName)

	msg := "Request from SR-APP"
	log.Printf("Creating KafkaProducer and writing message \"%s\" into topic", msg)
	produceMessage(CreateProducer(), msg)

	return &pb.InterfaceReply{ReceivedBytes: 100}, nil
}

func main() {
	lis, err := net.Listen("tcp", GRPC_PORT)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterInterfaceServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
