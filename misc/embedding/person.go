package models

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	DoB  string `json:"dob,omitempty"`
}

func (p Person) Id() string {
	return fmt.Sprintf("P-%s", p.ID)
}

func (p Person) Talk(message string) {
	fmt.Printf("%s (a person, ID: %s) says \"%s\"\n", p.Name, p.Id(), message)
}

func (p Person) ToJSON() (string, error) {
	bs, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
