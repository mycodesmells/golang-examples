package rdb

import (
	"github.com/go-redis/redis"
)

func Connect(addr, pass string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       1,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}
