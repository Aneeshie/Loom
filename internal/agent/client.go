package agent

import (
	"log"

	"github.com/Aneeshie/loom/internal/config"
	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient() (pb.LogServiceClient, error) {

	cfg, err := config.LoadAgentConfig("../../configs/agent.yaml")
	if err != nil {
		log.Fatalf("Could not load the agent config, %v", err)
	}

	conn, err := grpc.NewClient(
		cfg.Server.Addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		return nil, err
	}

	return pb.NewLogServiceClient(conn), nil
}
