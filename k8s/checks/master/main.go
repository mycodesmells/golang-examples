package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type workerResponse struct {
	Result float64 `json:"result,omitempty"`
}

func main() {
	addr := os.Getenv("ADDR")
	workerAddr := os.Getenv("WORKER_ADDR")

	// Usage:
	// curl -X GET 'http://localhost:8000/work?base=5&power=2'
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		base := req.FormValue("base")
		power := req.FormValue("power")

		url := fmt.Sprintf("%s/work?base=%s&power=%s", workerAddr, base, power)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to calculate power: %v", err)
			http.Error(rw, "failed to calculate power", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var jsonResp workerResponse
		if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
			log.Printf("Failed to decode power result: %v", err)
			http.Error(rw, "failed to decode power result", http.StatusInternalServerError)
			return
		}

		log.Printf("base=%s power=%s result=%f", base, power, jsonResp.Result)
		fmt.Fprint(rw, jsonResp.Result)
	})
	fmt.Println(http.ListenAndServe(addr, nil))
}
