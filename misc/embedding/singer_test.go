package models_test

import (
	"fmt"

	"github.com/mycodesmells/golang-examples/misc/embedding/models"
)

func ExampleSinger_Sing() {
	s := models.Singer{
		Person: models.Person{
			ID:   "687",
			Name: "John Singer",
		},
		MusicGenre: "pop",
	}
	s.Sing("La la lake")

	s2 := models.Singer{MusicGenre: "rock"}
	s2.Name = "Johny Singerra"
	s2.Sing("Welcome to the forest")

	// output:
	// John Singer (ID: S-687) sings La la lake in the style of pop.
	// Johny Singerra (ID: S-) sings Welcome to the forest in the style of rock.
}

func ExampleSinger_ToJSON() {
	s := models.Singer{
		Person: models.Person{
			ID:   "687",
			Name: "John Singer",
			DoB:  "01-02-1975",
		},
		MusicGenre: "pop",
	}
	sJSON, _ := s.ToJSON()
	fmt.Println(sJSON)
	// output: {"id":"687","name":"John Singer","dob":"01-02-1975","music_genre":"pop"}
}
