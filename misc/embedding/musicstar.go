package models

import (
	"encoding/json"
	"fmt"
)

type MusicStar struct {
	Singer
	ID       string `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	DoB      string `json:"dob,omitempty"`
}

func (p MusicStar) Id() string {
	return fmt.Sprintf("★-%s", p.ID)
}

func (ms MusicStar) GreetCrowd(city string) {
	fmt.Printf("%s (ID: %s) greets the people of %s!!\n", ms.Name, ms.Id(), city)
}

func (ms MusicStar) ToJSON() (string, error) {
	bs, err := json.Marshal(ms)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
