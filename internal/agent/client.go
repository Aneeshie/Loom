package agent

import (
	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient() (pb.LogServiceClient, error) {
	conn, err := grpc.NewClient(
		":8080",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		return nil, err
	}

	return pb.NewLogServiceClient(conn), nil
}
