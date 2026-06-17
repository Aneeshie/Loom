package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Aneeshie/loom/internal/agent"

	pb "github.com/Aneeshie/loom/proto"
)

var (
	server  = flag.String("server", "localhost:8080", "server address")
	level   = flag.String("level", "", "filter by log level")
	service = flag.String("service", "", "filter by service name")
	host    = flag.String("host", "", "filter by host")
	limit   = flag.Int64("limit", 20, "max results")
	follow  = flag.Bool("follow", false, "Follow logs in realtime")
	search  = flag.String("search", "", "search logs by message")
	since   = flag.Duration("since", 0, "show logs since duration (e.g. 1h, 24h, 7m)")
	similar = flag.String("similar", "", "semantic search query")
)

func main() {
	flag.Parse()

	client, err := agent.NewClient(*server)
	if err != nil {
		log.Fatal(err)
	}

	if *similar != "" {
		runSemanticSearch(client, *similar, *limit)
		return
	}

	var filter *pb.LogFilter

	if *level != "" || *service != "" || *host != "" || *search != "" {

		filter = &pb.LogFilter{}

		if *level != "" {

			filter.Level = level // already a *string

		}

		if *service != "" {

			filter.ServiceName = service

		}

		if *host != "" {

			filter.Host = host

		}

		if *search != "" {
			filter.Search = search
		}

		if *since != 0 {
			filter.StartTime = new(time.Now().Add(-*since).Unix())
		}

	}

	if *follow {
		runFollow(client, filter)
		return
	}

	resp, err := client.GetLogs(context.Background(), &pb.GetLogsRequest{
		Limit:  *limit,
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

func runFollow(client pb.LogServiceClient, filter *pb.LogFilter) {
	stream, err := client.StreamLogs(context.Background(), &pb.StreamLogsRequest{
		Filter: filter,
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		logEntry, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"[%s] %s - %s\n",
			logEntry.Level,
			logEntry.ServiceName,
			logEntry.Message,
		)
	}
}

func runSemanticSearch(client pb.LogServiceClient, query string, limit int64) {
	resp, err := client.SimilarLogs(context.Background(), &pb.SimilarLogsRequest{Query: query, Limit: limit})
	if err != nil {
		log.Fatal(err)
	}

	for _, logEntry := range resp.Logs {
		fmt.Printf(
			"[%s] %s - %s\n",
			logEntry.Level,
			logEntry.ServiceName,
			logEntry.Message,
		)
	}
}
