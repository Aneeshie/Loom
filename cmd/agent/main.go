package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/Aneeshie/loom/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// for now
const (
	serviceName = "test_service"
	host        = "localhost"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))

	conn, err := grpc.NewClient(":8080", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//TODO: move log file path to config
	file, err := os.Open("../../testdata/test.log")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	client := pb.NewLogServiceClient(conn)

	for scanner.Scan() {
		line := scanner.Text()

		level, message, err := parseLogLine(line)
		if err != nil {
			log.Printf("failed to send log: %v", err)
			continue
		}

		resp, err := client.SendLog(context.Background(), &pb.SendLogRequest{
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

func parseLogLine(logLine string) (string, string, error) {

	start := strings.Index(logLine, "[")
	end := strings.Index(logLine, "]")

	if start != -1 && end != -1 && end > start {
		level := logLine[start+1 : end]

		message := strings.TrimSpace(logLine[end+1:])

		return level, message, nil

	}
	return "", "", fmt.Errorf("invalid format")

}
