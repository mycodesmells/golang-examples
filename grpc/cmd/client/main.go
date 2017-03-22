package main

import (
	"log"
	"context"

	"google.golang.org/grpc"

	pb "github.com/mycodesmells/golang-examples/grpc/proto/service"
)

func main() {
	conn, err := grpc.Dial("localhost:6000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleServerClient(conn)

	_, err = client.CreateUser(context.Background(), &pb.CreateUserRequest{User: &pb.User{Username: "slomek", Role: "joker"}})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Println("Created user!")

	resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Username: "slomek"})
	if err != nil {
		log.Fatalf("Failed to get created user: %v", err)
	}
	log.Printf("User exists: %v\n", resp)

	resp2, err := client.GreetUser(context.Background(), &pb.GreetUserRequest{Greeting: "howdy", Username: "slomek"})
	if err != nil {
		log.Fatalf("Failed to greet user: %v", err)
	}
	log.Printf("Greeting: %s\n", resp2.Greeting)
}
