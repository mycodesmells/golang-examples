package models_test

import (
	"fmt"

	"github.com/mycodesmells/golang-examples/misc/embedding/models"
)

func ExamplePerson_Talk() {
	p := models.Person{ID: "123", Name: "John Doe"}
	p.Talk("Hi there!")
	// output: John Doe (a person, ID: P-123) says "Hi there!"
}

func ExamplePerson_ToJSON() {
	p := models.Person{
		ID:   "123",
		Name: "John Doe",
		DoB:  "01-02-1975",
	}
	pJSON, _ := p.ToJSON()
	fmt.Println(pJSON)
	// output: {"id":"123","name":"John Doe","dob":"01-02-1975"}
}
