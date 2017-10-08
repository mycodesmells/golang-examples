package models

import (
	"encoding/json"
	"fmt"
)

type Singer struct {
	Person
	MusicGenre string `json:"music_genre,omitempty"`
}

func (s Singer) Id() string {
	return fmt.Sprintf("S-%s", s.ID)
}

func (s Singer) Sing(title string) {
	fmt.Printf("%s (ID: %s) sings %s in the style of %s.\n", s.Name, s.Id(), title, s.MusicGenre)
}

func (s Singer) ToJSON() (string, error) {
	bs, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
