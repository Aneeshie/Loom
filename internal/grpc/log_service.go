package grpc

import (
	"context"

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

	err := s.store.InsertLog(ctx, req)

	if err != nil {
		return nil, err
	}

	return &pb.SendLogResponse{
		Message: "received",
	}, nil
}
