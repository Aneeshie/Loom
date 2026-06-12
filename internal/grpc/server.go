package grpc

import (
	"log"
	"net"

	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
}

func NewServer() *GRPCServer {

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLogServiceServer(grpcServer, &LogService{})

	return &GRPCServer{
		server: grpcServer,
	}

}

func (g *GRPCServer) Run() {

	//TODO: get ip addr from cfg
	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC server listening on :8080")

	if err := g.server.Serve(list); err != nil {
		log.Fatal(err)
	}
}
