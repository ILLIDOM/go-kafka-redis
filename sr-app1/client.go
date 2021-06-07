package main

import (
	"context"
	"log"
	"time"

	pb "example.com/prototype/proto"
	"google.golang.org/grpc"
)

const (
	ADDRESS = "localhost:50500"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(ADDRESS, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewInterfaceServiceClient(conn)

	// Contact server and print out response
	interfaceName := "Gi0/0/0/1"
	routerName := "XR-1"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetReceivedBytes(ctx, &pb.InterfaceRequest{InterfaceName: interfaceName, RouterName: routerName})
	if err != nil {
		log.Fatalf("No Bytes received: %v", err)
	}
	log.Printf("Bytes Received are: %d", r.GetReceivedBytes())
}
