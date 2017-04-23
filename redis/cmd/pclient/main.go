package main

import (
	"fmt"
	"log"
	"time"

	redis "gopkg.in/redis.v4"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	pubsub, err := client.Subscribe("pubsub-key")
	if err != nil {
		log.Fatalf("Failed to get task from queue: %v\n", err)
	}

	for {
		msgi, err := pubsub.ReceiveTimeout(5 * time.Second)
		if err != nil {
			break
		}

		switch msg := msgi.(type) {
		case *redis.Subscription:
			fmt.Println("subscribed to", msg.Channel)
			time.Sleep(time.Second * 10)
		case *redis.Message:
			fmt.Println("received", msg.Payload, "from", msg.Channel)
		default:
			panic(fmt.Errorf("unknown message: %#v", msgi))
		}
	}

	// log.Printf("Working on '%s' task...\n", task[1])
}
