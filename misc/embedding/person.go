package embedding

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name,omitempty"`
	DoB  string `json:"dob,omitempty"`
}

func (p Person) Type() string {
	return "PERSON"
}

func (p Person) Talk(message string) {
	fmt.Printf("%s (type=%s) says \"%s\"\n", p.Name, p.Type(), message)
}

func (p Person) ToJSON() (string, error) {
	bs, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
