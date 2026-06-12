package agent

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/Aneeshie/loom/proto"
)

const (
	serviceName = "test"
	host        = "ubuntu"
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
	//TODO: add test later to cfg path
	file, err := os.Open("../../testdata/test.log")

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
			ServiceName: serviceName,
			Host:        host,
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
