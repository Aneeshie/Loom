package grpc

import (
	"context"
	"log"

	pb "github.com/Aneeshie/loom/proto"
)

type LogService struct {
	pb.UnimplementedLogServiceServer
}

func (s *LogService) SendLog(ctx context.Context, req *pb.SendLogRequest) (*pb.SendLogResponse, error) {

	log.Printf(
		"[%s] %s",
		req.Level,
		req.Message,
	)
	return &pb.SendLogResponse{
		Message: "received",
	}, nil
}
