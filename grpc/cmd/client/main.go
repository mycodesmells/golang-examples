package main

import (
	"context"
	"log"

	"google.golang.org/grpc"

	pb "github.com/mycodesmells/golang-examples/grpc/proto/service"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("cmd/server/server-cert.pem", "")
	if err != nil {
		log.Fatalf("cert load error: %s", err)
	}

	conn, err := grpc.Dial("localhost:6000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}
	defer conn.Close()

	client := pb.NewSimpleServerClient(conn)

	md := metadata.Pairs("token", "valid-token")
	ctx := metadata.NewContext(context.Background(), md)

	_, err = client.CreateUser(ctx, &pb.CreateUserRequest{User: &pb.User{Username: "slomek", Role: "joker"}})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Println("Created user!")

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Username: "slomek"})
	if err != nil {
		log.Fatalf("Failed to get created user: %v", err)
	}
	log.Printf("User exists: %v\n", resp)

	resp2, err := client.GreetUser(ctx, &pb.GreetUserRequest{Greeting: "howdy", Username: "slomek"})
	if err != nil {
		log.Fatalf("Failed to greet user: %v", err)
	}
	log.Printf("Greeting: %s\n", resp2.Greeting)
}
