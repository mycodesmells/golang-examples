package main

import (
	"fmt"

	"github.com/mycodesmells/golang-examples/grpc/proto/message"
)

func main() {
	p := message.Person{
		FirstName:    "John",
		LastName:     "Doe",
		DateOfBirth:  "1960-10-17T0:00:00Z",
		Cool:         true,
		ArgumentsWon: 7,
		Hobbies: []*message.Hobby{
			{
				Name:        "Running",
				Description: "Occasionally, about 10km a week",
			}, {
				Name:        "Computer games",
				Description: "Flappy bird, mostly",
			},
		},
	}

	fmt.Printf("Person created for .proto structure: %v\n", p)

	fmt.Printf("Full name (custom fn): %s\n", p.FullName())
}
