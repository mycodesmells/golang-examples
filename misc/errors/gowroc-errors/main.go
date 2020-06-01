package main

import (
	"errors"
	"fmt"
)

var errSome = errors.New("some error")
var errSome2 = errors.New("some error2")

type Person struct {
	FirstName string `map:"first_name"`
	LastName  string `map:"last_name"`
	Number    int    `map:"number"`
}

func main() {
	// pMap := map[string]string{
	// 	"first_name": "Shaquille",
	// 	"last_name":  "O'Neal",
	// 	"number":     "34",
	// }

	// var p Person
	// err := mappy.Unmarshal(pMap, &p)
	// if err != nil {
	// 	if errors.Is(err, mappy.ErrMapUnmarshal) {
	// 		fmt.Println("my bad, my input data was wrong")
	// 	} else {
	// 		fmt.Println("slomek/mappy is crappy")
	// 	}
	// }

	err := fmt.Errorf("errors %w %w", errSome, errSome2)
	fmt.Println(errors.Is(err, errSome))
	fmt.Println(errors.Is(err, errSome2))
}
