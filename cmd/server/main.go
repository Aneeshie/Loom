package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedLogServiceServer
}

func (s *Server) SendLog(ctx context.Context, req *pb.SendLogRequest) (*pb.SendLogResponse, error) {

	log.Printf(
		"[%s] %s",
		req.Level,
		req.Message,
	)
	return &pb.SendLogResponse{
		Message: "received",
	}, nil
}

func main() {
	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLogServiceServer(grpcServer, &Server{})

	log.Println("gRPC server listening on :8080")

	if err := grpcServer.Serve(list); err != nil {
		log.Fatal(err)
	}

}
