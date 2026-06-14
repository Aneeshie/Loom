package grpc

import (
	"context"
	"log"

	"github.com/Aneeshie/loom/internal/storage"
	pb "github.com/Aneeshie/loom/proto"
)

type LogService struct {
	pb.UnimplementedLogServiceServer
	store *storage.Store
}

func NewLogService(store *storage.Store) *LogService {
	return &LogService{
		store: store,
	}
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
