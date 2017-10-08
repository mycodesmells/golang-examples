package embedding

import (
	"encoding/json"
	"fmt"
)

type Singer struct {
	Person
	MusicGenre string `json:"music_genre,omitempty"`
}

func (s Singer) Type() string {
	return "SINGER"
}

func (s Singer) Sing(title string) {
	fmt.Printf("%s (type=%s) sings %s in the style of %s.\n", s.Name, s.Type(), title, s.MusicGenre)
}

func (s Singer) ToJSON() (string, error) {
	bs, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
