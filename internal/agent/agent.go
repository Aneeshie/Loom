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
	cfg    *config.AgentConfig
}

func New(client pb.LogServiceClient, cfg *config.AgentConfig) *Agent {
	return &Agent{
		client: client,
		cfg:    cfg,
	}
}

func (a *Agent) Run() {
	file, err := os.Open(a.cfg.Source.Path)

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
			ServiceName: a.cfg.Service.Name,
			Host:        a.cfg.Service.Host,
			Level:       level,
			Message:     message,
			Timestamp:   time.Now().Unix(),
		})

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(resp.Message)

		logResp, err := a.client.GetLogs(
			context.Background(),
			&pb.GetLogsRequest{},
		)

		if err != nil {
			log.Fatal(err)
		}

		for _, logEntry := range logResp.Logs {
			fmt.Printf(
				"[%s] %s\n",
				logEntry.Level,
				logEntry.Message,
			)
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
