package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	addr := os.Getenv("ADDR")
	redisAddr := os.Getenv("REDIS_ADDR")

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Usage:
	// curl -X GET 'http://localhost:8000/work?base=5&power=2'
	http.HandleFunc("/work", func(rw http.ResponseWriter, req *http.Request) {
		baseStr := req.FormValue("base")
		base, err := strconv.ParseFloat(baseStr, 64)
		if err != nil {
			http.Error(rw, "invalid value for 'base'", http.StatusBadRequest)
			return
		}

		powerStr := req.FormValue("power")
		power, err := strconv.ParseFloat(powerStr, 64)
		if err != nil {
			http.Error(rw, "invalid value for 'power'", http.StatusBadRequest)
			return
		}

		result, err := redisClient.Get(cacheKey(base, power)).Float64()
		if err != nil {
			if err != redis.Nil {
				log.Printf("Error connecting to Redis: %v", err)
			} else {
				log.Printf("Cache miss for base=%f power=%f", base, power)
			}
			result = math.Pow(base, power)
		}

		if err := redisClient.Set(cacheKey(base, power), result, time.Minute).Err(); err != nil {
			log.Printf("Failed to cache result: %v", err)
		}

		log.Printf("base=%f power=%f result=%f", base, power, result)
		fmt.Fprintf(rw, `{"base": %f, "power": %f, "result": %f}`, base, power, result)
	})

	// Liveness check
	http.HandleFunc("/checks/liveness", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, "OK")
	})
	// Readiness check
	http.HandleFunc("/checks/readiness", func(rw http.ResponseWriter, req *http.Request) {
		if err := redisClient.Ping().Err(); err != nil {
			http.Error(rw, "FAIL", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(rw, "OK")
	})

	fmt.Println(http.ListenAndServe(addr, nil))
}

func cacheKey(base, power float64) string {
	return fmt.Sprintf("%f:%f", base, power)
}
