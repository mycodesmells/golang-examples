package embedding_test

import (
	"fmt"

	"github.com/mycodesmells/golang-examples/misc/embedding"
)

func ExampleMusicStar() {
	ms := embedding.MusicStar{
		Nickname: "Starry",
		Singer: embedding.Singer{
			Person: embedding.Person{
				Name: "Joe Star",
			},
			MusicGenre: "pop",
		},
	}
	ms.GreetCrowd("Oklahoma")
	ms.Sing("To the top")
	ms.Talk("Thank you!!")
	// output:
	// Joe Star (type=MUSIC★) greets the people of Oklahoma!!
	// Joe Star (type=SINGER) sings To the top in the style of pop.
	// Joe Star (type=PERSON) says "Thank you!!"
}

func ExampleMusicStar_ToJSON() {
	ms := embedding.MusicStar{
		Nickname: "Starry",
		Singer: embedding.Singer{
			Person: embedding.Person{
				Name: "Joe Star",
				DoB:  "01-02-1975",
			},
			MusicGenre: "pop",
		},
	}
	msJSON, _ := ms.ToJSON()
	fmt.Println(msJSON)
	// output: {"name":"Joe Star","music_genre":"pop","nickname":"Starry"}
}
