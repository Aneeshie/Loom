package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Aneeshie/loom/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	client := pb.NewLogServiceClient(conn)
	resp, err := client.SendLog(context.Background(), &pb.SendLogRequest{
		ServiceName: "test",
		Host:        "ubuntu",
		Level:       "INFO",
		Message:     "A test Log",
		Timestamp:   12321,
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Message)

}
