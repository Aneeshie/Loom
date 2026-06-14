package grpc

import (
	"log"
	"net"

	"github.com/Aneeshie/loom/internal/config"
	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewServer(logService *LogService) *GRPCServer {

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLogServiceServer(grpcServer, logService)

	return &GRPCServer{
		server: grpcServer,
	}

}

func (g *GRPCServer) Run() {

	cfg, err := config.LoadServerConfig("../../configs/server.yaml")

	if err != nil {
		log.Fatalf("Could not load the server config %v", err)
	}

	list, err := net.Listen("tcp", cfg.GRPC.Addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC server listening on :8080")

	if err := g.server.Serve(list); err != nil {
		log.Fatal(err)
	}
}
