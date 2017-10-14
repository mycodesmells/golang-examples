package examples

import "fmt"

type User struct {
	Name string
}

func (u User) Hi() {
	fmt.Printf("Hi, my name is %s!", u.Name)
}
