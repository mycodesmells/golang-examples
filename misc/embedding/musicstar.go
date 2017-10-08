package embedding

import (
	"encoding/json"
	"fmt"
)

type MusicStar struct {
	Singer
	Nickname string `json:"nickname,omitempty"`
	DoB      string `json:"dob,omitempty"`
}

func (p MusicStar) Type() string {
	return "MUSIC★"
}

func (ms MusicStar) GreetCrowd(city string) {
	fmt.Printf("%s (type=%s) greets the people of %s!!\n", ms.Name, ms.Type(), city)
}

func (ms MusicStar) ToJSON() (string, error) {
	bs, err := json.Marshal(ms)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
