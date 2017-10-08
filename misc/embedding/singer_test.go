package embedding_test

import (
	"fmt"

	"github.com/mycodesmells/golang-examples/misc/embedding"
)

func ExampleSinger_Sing() {
	s := embedding.Singer{
		Person: embedding.Person{
			Name: "John Singer",
		},
		MusicGenre: "pop",
	}
	s.Sing("La la lake")
	s.Talk("Hi!")

	s2 := embedding.Singer{MusicGenre: "rock"}
	s2.Name = "Johny Singerra"
	s2.Sing("Welcome to the forest")
	s2.Talk("Hello!")

	// output:
	// John Singer (type=SINGER) sings La la lake in the style of pop.
	// John Singer (type=PERSON) says "Hi!"
	// Johny Singerra (type=SINGER) sings Welcome to the forest in the style of rock.
	// Johny Singerra (type=PERSON) says "Hello!"
}

func ExampleSinger_ToJSON() {
	s := embedding.Singer{
		Person: embedding.Person{
			Name: "John Singer",
			DoB:  "01-02-1975",
		},
		MusicGenre: "pop",
	}
	sJSON, _ := s.ToJSON()
	fmt.Println(sJSON)
	// output: {"name":"John Singer","dob":"01-02-1975","music_genre":"pop"}
}
