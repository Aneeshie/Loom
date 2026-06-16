package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Aneeshie/loom/internal/agent"

	pb "github.com/Aneeshie/loom/proto"
)

var (
	server = flag.String("server", "localhost:8080", "server address")
	level = flag.String("level", "", "filter by log level")
	service = flag.String("service", "", "filter by service name")
	host = flag.String("host", "", "filter by host")
	limit = flag.Int64("limit", 20, "max results")
)

func main() {
	flag.Parse()

	client, err := agent.NewClient(*server)
	if err != nil {
		log.Fatal(err)
	}

	var filter *pb.LogFilter


	if *level != "" || *service != "" || *host != "" {

		filter = &pb.LogFilter{}

		if *level != "" {

			filter.Level = level  // already a *string

		}

		if *service != "" {

			filter.ServiceName = service

		}

		if *host != "" {

			filter.Host = host

		}

	}

	resp, err := client.GetLogs(context.Background(),&pb.GetLogsRequest{
		Limit: *limit,
		Filter: filter,
	})

	if err != nil {
		log.Fatal(err)
	}

	if resp.Logs == nil {

		fmt.Println("no logs found")

		return

	}

	for _, log := range resp.Logs {
		fmt.Printf("[%s] %s - %s - %d\n", log.Level, log.ServiceName, log.Message, log.Timestamp)
	}
}
