package message

import "fmt"

func (p Person) FullName() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}
