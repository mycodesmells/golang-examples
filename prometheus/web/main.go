package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mycodesmells/golang-examples/prometheus/web/metrics"
)

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	handler.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		delay := int32(r.Float64() * 10000)
		time.Sleep(time.Millisecond * time.Duration(delay))

		fmt.Printf("Responed with 'Hello' in %dms\n", delay)
		rw.Write([]byte("Hello!"))
	})

	withMetrics := metrics.Middleware(handler)

	log.Fatal(http.ListenAndServe(":3000", withMetrics))
}
