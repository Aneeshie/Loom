package agent

import (
	"log"

	"github.com/Aneeshie/loom/internal/config"
	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient() (pb.LogServiceClient, error) {

	cfg, err := config.LoadServerConfig("../../configs/server.yaml")
	if err != nil {
		log.Fatalf("Could not load the server config, %v", err)
	}

	conn, err := grpc.NewClient(
		cfg.GRPC.Addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		return nil, err
	}

	return pb.NewLogServiceClient(conn), nil
}
