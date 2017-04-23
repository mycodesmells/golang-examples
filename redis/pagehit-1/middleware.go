package pagehit

import (
	"fmt"
	"log"
	"net/http"

	redis "gopkg.in/redis.v4"

	"github.com/labstack/echo"
)

type Stats map[string]int

type Store interface {
	GetStats() (Stats, error)
	Hit(page string) error
}

func Middleware(hs Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			hs.Hit(ctx.Path())
			return next(ctx)
		}
	}
}

func Handler(hs Store) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		stats, err := hs.GetStats()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		return ctx.JSON(http.StatusOK, stats)
	}
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore() Store {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}

	return &redisStore{
		client: client,
	}
}

func (r redisStore) GetStats() (Stats, error) {
	var pages []string
	pages, err := r.client.SMembers("pages").Result()
	if err != nil {
		return Stats{}, err
	}

	stats := make(map[string]int)
	for _, p := range pages {
		count, err := r.client.Get(p).Int64()
		if err == nil {
			// return stats, err
			stats[p] = int(count)
		}
	}

	return stats, nil
}

func (r redisStore) Hit(url string) error {
	i, err := r.client.SAdd("pages", url).Result()
	if err != nil {
		return err
	}

	if i == 1 {
		fmt.Printf("First time hit to %s/n", url)
	}

	_, err = r.client.Incr(url).Result()

	return err
}
