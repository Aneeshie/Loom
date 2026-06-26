package service

import (
	"context"
	"github.com/Aneeshie/loom/internal/nlq"
	"github.com/Aneeshie/loom/internal/web/types"
	pb "github.com/Aneeshie/loom/proto"
	"log"
)

type QueryService struct {
	parser nlq.IntentParser
	client pb.LogServiceClient
}

func NewQueryService(parser nlq.IntentParser, client pb.LogServiceClient) *QueryService {
	return &QueryService{
		parser: parser,
		client: client,
	}
}

func (s *QueryService) Query(ctx context.Context, query string) (*types.QueryResponse, error) {
	// parse the ollama
	intent, err := s.parser.Parse(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	// build a LogFilter

	filter := buildLogFilter(intent)

	// decide between SimilarLogs or GetLogs
	var logs []*pb.Log

	if intent.Query != " " {
		resp, err := s.client.SimilarLogs(ctx, &pb.SimilarLogsRequest{
			Query:  intent.Query,
			Limit:  20,
			Filter: filter,
		})

		if err != nil {
			return nil, err
		}

		logs = resp.Logs
	} else {
		resp, err := s.client.GetLogs(ctx, &pb.GetLogsRequest{
			Limit:  20,
			Filter: filter,
		})

		if err != nil {
			return nil, err
		}

		logs = resp.Logs
	}
	// return a query Response
	return &types.QueryResponse{
		Intent: mapIntent(intent),
		Logs:   mapLogs(logs),
	}, nil
}

func mapIntent(intent *nlq.QueryIntent) types.IntentResponse {
	return types.IntentResponse{
		Query:   intent.Query,
		Level:   intent.Level,
		Service: intent.Service,
		Host:    intent.Host,
		Since:   intent.Since,
	}
}

func mapLogs(logs []*pb.Log) []types.LogResponse {
	resp := make([]types.LogResponse, 0, len(logs))

	for _, log := range logs {
		resp = append(resp, types.LogResponse{
			ID:          log.Id,
			ServiceName: log.ServiceName,
			Host:        log.Host,
			Level:       log.Level,
			Message:     log.Message,
			Timestamp:   log.Timestamp,
		})
	}

	return resp
}
