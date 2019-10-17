package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/slomek/mappy"
)

type Person struct {
	FirstName string `map:"first_name"`
	LastName  string `map:"last_name"`
	Number    int    `map:"number"`
}

func main() {
	data := map[string]string{
		"first_name": "Tim",
		"last_name":  "Duncan",
		"number":     "21",
	}

	var p Person
	if err := mappy.Unmarshal(data, &p); err != nil {
		if errors.Is(err, mappy.ErrMapUnmarshal) {
			fmt.Println("I'm so sorry, it's my fault - bad input data!")
			os.Exit(1)
		}

		fmt.Printf("slomek/mappy is crappy! Good input data, but returns error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Person: %+v\n", p)
}
