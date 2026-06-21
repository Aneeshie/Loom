package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Aneeshie/loom/internal/agent"
	"github.com/Aneeshie/loom/internal/nlq"

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
	ask     = flag.String("ask", "", "ask a natural language question about your logs")
)

func main() {
	flag.Parse()

	client, err := agent.NewClient(*server)
	if err != nil {
		log.Fatal(err)
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

	if *similar != "" {
		runSemanticSearch(client, *similar, *limit, filter)
		return
	}

	if *follow {
		runFollow(client, filter)
		return
	}

	if *ask != "" {
		runNLQ(client, *ask, *limit)
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

func runSemanticSearch(client pb.LogServiceClient, query string, limit int64, filter *pb.LogFilter) {
	resp, err := client.SimilarLogs(context.Background(), &pb.SimilarLogsRequest{Query: query, Limit: limit, Filter: filter})
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

func runNLQ(
	client pb.LogServiceClient,
	query string,
	limit int64,
) {

	parser := nlq.NewOllamaParser("llama3.2:3b")

	intent, err := parser.Parse(
		context.Background(),
		query,
	)

	if err != nil {
		log.Fatal(err)
	}

	filter := &pb.LogFilter{}

	if intent.Level != "" {
		filter.Level = &intent.Level
	}

	if intent.Service != "" {
		filter.ServiceName = &intent.Service
	}

	if intent.Host != "" {
		filter.Host = &intent.Host
	}

	if intent.Since != "" {

		duration, err := time.ParseDuration(
			intent.Since,
		)

		if err == nil {

			start := time.Now().
				Add(-duration).
				Unix()

			filter.StartTime = &start
		}
	}

	// semantic search path
	if intent.Query != "" {

		runSemanticSearch(
			client,
			intent.Query,
			limit,
			filter,
		)

		return
	}

	// filter-only path
	resp, err := client.GetLogs(
		context.Background(),
		&pb.GetLogsRequest{
			Limit:  limit,
			Filter: filter,
		},
	)

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
