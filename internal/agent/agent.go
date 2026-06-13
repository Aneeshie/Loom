package agent

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Aneeshie/loom/internal/config"
	pb "github.com/Aneeshie/loom/proto"
)

type Agent struct {
	client pb.LogServiceClient
}

func New(client pb.LogServiceClient) *Agent {
	return &Agent{
		client: client,
	}
}

func (a *Agent) Run() {
	cfg, err := config.LoadAgentConfig("../../configs/agent.yaml")

	if err != nil {
		log.Fatalf("Cannot load agent config %v", err)
	}

	file, err := os.Open(cfg.Source.Path)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		level, message, err := ParseLogLine(line)
		if err != nil {
			log.Printf("failed to send log: %v", err)
			continue
		}

		resp, err := a.client.SendLog(context.Background(), &pb.SendLogRequest{
			ServiceName: cfg.Service.Name,
			Host:        cfg.Service.Host,
			Level:       level,
			Message:     message,
			Timestamp:   time.Now().Unix(),
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp.Message)

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
