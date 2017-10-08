package models_test

import (
	"fmt"

	"github.com/mycodesmells/golang-examples/misc/embedding/models"
)

func ExampleMusicStar() {
	ms := models.MusicStar{
		ID:       "STAR",
		Nickname: "Starry",
		Singer: models.Singer{
			Person: models.Person{
				ID:   "person",
				Name: "Joe Star",
			},
			MusicGenre: "pop",
		},
	}
	ms.GreetCrowd("Oklahoma")
	ms.Sing("To the top")
	ms.Talk("Thank you!!")
	// output:
	// Joe Star (ID: ★-STAR) greets the people of Oklahoma!!
	// Joe Star (ID: S-person) sings To the top in the style of pop.
	// Joe Star (a person, ID: P-person) says "Thank you!!"
}

func ExampleMusicStar_ToJSON() {
	ms := models.MusicStar{
		ID:       "STAR",
		Nickname: "Starry",
		Singer: models.Singer{
			Person: models.Person{
				ID:   "person",
				Name: "Joe Star",
				DoB:  "01-02-1975",
			},
			MusicGenre: "pop",
		},
	}
	msJSON, _ := ms.ToJSON()
	fmt.Println(msJSON)
	// output: {"name":"Joe Star","music_genre":"pop","id":"STAR","nickname":"Starry"}
}
