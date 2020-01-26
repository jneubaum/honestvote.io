package main

import (
	"fmt"

	"github.com/jneubaum/honestvote/core/core-database/database"
)

func main() {
	vote := database.Vote{
		Type:     "Vote",
		Election: "Chester",
		Receiver: map[string]string{"cool": "beans"},
	}

	voteHeaders := vote.Type + vote.Election
	for key, value := range vote.Receiver {
		voteHeaders += key + value
	}

	fmt.Println([]byte(voteHeaders))
}
