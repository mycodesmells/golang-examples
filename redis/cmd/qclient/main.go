package main

import (
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

	for {
		time.Sleep(5 * time.Second)

		task, err := client.BLPop(0, "queue-key").Result()
		if err != nil {
			log.Fatalf("Failed to get task from queue: %v\n", err)
		}

		log.Printf("Working on '%s' task...\n", task[1])
	}
}
