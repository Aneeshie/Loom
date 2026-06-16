package grpc

import (
	"context"
	"sync"

	"github.com/Aneeshie/loom/internal/storage"
	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
)

type LogService struct {
	pb.UnimplementedLogServiceServer
	store *storage.Store

	mu          sync.RWMutex
	subscribers map[*Subscriber]struct{}
}

func NewLogService(store *storage.Store) *LogService {
	return &LogService{
		store: store,

		subscribers: make(map[*Subscriber]struct{}),
	}
}

func (s *LogService) SendLog(ctx context.Context, req *pb.SendLogRequest) (*pb.SendLogResponse, error) {

	err := s.store.InsertLog(ctx, req)

	if err != nil {
		return nil, err
	}

	s.broadcast(&pb.Log{
		ServiceName: req.ServiceName,
		Host:        req.Host,
		Level:       req.Level,
		Message:     req.Message,
		Timestamp:   req.Timestamp,
	})

	return &pb.SendLogResponse{
		Message: "received",
	}, nil
}

func (s *LogService) GetLogs(ctx context.Context, req *pb.GetLogsRequest) (*pb.GetLogsResponse, error) {
	logs, err := s.store.GetLogs(ctx, req.Limit, req.Filter)

	if err != nil {
		return nil, err
	}

	return &pb.GetLogsResponse{
		Logs: logs,
	}, nil

}

func (s *LogService) StreamLogs(req *pb.StreamLogsRequest, stream grpc.ServerStreamingServer[pb.Log]) error {
	sub := &Subscriber{
		ch: make(chan *pb.Log),
	}

	s.mu.Lock()
	s.subscribers[sub] = struct{}{}
	s.mu.Unlock()

	s.mu.Lock()

	defer delete(s.subscribers, sub)

	s.mu.Unlock()

	for {
		log := <-sub.ch

		err := stream.Send(log)

		if err != nil {
			return err
		}
	}

}

func (s *LogService) broadcast(log *pb.Log) {
	s.mu.Lock()
	for sub := range s.subscribers {
		sub.ch <- log
	}

	s.mu.RUnlock()
}
