package embedding_test

import (
	"fmt"

	"github.com/mycodesmells/golang-examples/misc/embedding"
)

func ExamplePerson_Talk() {
	p := embedding.Person{Name: "John Doe"}
	p.Talk("Hi there!")
	// output: John Doe (type=PERSON) says "Hi there!"
}

func ExamplePerson_ToJSON() {
	p := embedding.Person{
		Name: "John Doe",
		DoB:  "01-02-1975",
	}
	pJSON, _ := p.ToJSON()
	fmt.Println(pJSON)
	// output: {"name":"John Doe","dob":"01-02-1975"}
}
