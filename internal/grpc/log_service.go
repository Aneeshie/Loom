package grpc

import (
	"context"
	"strings"
	"sync"

	"github.com/Aneeshie/loom/internal/embeddings"
	"github.com/Aneeshie/loom/internal/storage"
	pb "github.com/Aneeshie/loom/proto"
	"google.golang.org/grpc"
)

type LogService struct {
	pb.UnimplementedLogServiceServer
	store *storage.Store

	mu          sync.RWMutex
	subscribers map[*Subscriber]struct{}

	embedder embeddings.Embedder
}

func NewLogService(store *storage.Store, embedder embeddings.Embedder) *LogService {
	return &LogService{
		store: store,

		subscribers: make(map[*Subscriber]struct{}),

		embedder: embedder,
	}
}

func (s *LogService) SendLog(ctx context.Context, req *pb.SendLogRequest) (*pb.SendLogResponse, error) {

	embedding, err := s.embedder.Embed(ctx, req.Message)
	if err != nil {
		return nil, err
	}

	err = s.store.InsertLog(ctx, req, embedding)

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
		ch:     make(chan *pb.Log, 100), // HARDCODING FOR NOW
		filter: req.Filter,
	}

	s.mu.Lock()
	s.subscribers[sub] = struct{}{}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.subscribers, sub)
		s.mu.Unlock()
	}()

	ctx := stream.Context()

	for {
		select {
		case <-ctx.Done():
			return nil
		case log := <-sub.ch:
			if err := stream.Send(log); err != nil {
				return err
			}
		}
	}

}

func (s *LogService) broadcast(log *pb.Log) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for sub := range s.subscribers {
		if matchesFilter(log, sub.filter) {
			sub.ch <- log
		}
	}
}

func matchesFilter(log *pb.Log, filter *pb.LogFilter) bool {
	if filter == nil {
		return true
	}

	if filter.Level != nil &&
		log.Level != *filter.Level {
		return false
	}

	if filter.ServiceName != nil &&
		log.ServiceName != *filter.ServiceName {
		return false
	}

	if filter.Host != nil &&
		log.Host != *filter.Host {
		return false
	}

	if filter.Search != nil && !strings.Contains(strings.ToLower(log.Message), strings.ToLower(*filter.Search)) {
		return false
	}

	return true

}

func (s *LogService) SimilarLogs(ctx context.Context, req *pb.SimilarLogsRequest) (*pb.GetLogsResponse, error) {
	embedding, err := s.embedder.Embed(
		ctx,
		req.Query,
	)

	if err != nil {
		return nil, err
	}

	logs, err := s.store.SimilarLogs(
		ctx,
		embedding,
		req.Limit,
		req.Filter,
	)

	if err != nil {
		return nil, err
	}

	return &pb.GetLogsResponse{
		Logs: logs,
	}, nil
}
