package main

import (
	"log"

	"github.com/Aneeshie/loom/internal/agent"
)

func main() {
	client, err := agent.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	a := agent.New(client)
	a.Run()

}
