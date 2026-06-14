package main

import (
	"log"

	"github.com/Aneeshie/loom/internal/agent"
	"github.com/Aneeshie/loom/internal/config"
)

func main() {
	cfg, err := config.LoadAgentConfig("../../configs/agent.yaml")
	if err != nil {
		log.Fatalf("Could not load the agent config, %v", err)
	}

	logServiceClient, err := agent.NewClient(cfg.Server.Addr)

	if err != nil {
		log.Fatal(err)
	}

	a := agent.New(logServiceClient, cfg)
	a.Run()

}
