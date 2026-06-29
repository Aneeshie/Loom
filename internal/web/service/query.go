package service

import (
	"context"
	"github.com/Aneeshie/loom/internal/nlq"
	"github.com/Aneeshie/loom/internal/web/types"
	pb "github.com/Aneeshie/loom/proto"
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

func (s *QueryService) Query(ctx context.Context, req types.QueryRequest) (*types.QueryResponse, error) {
	var intent *nlq.QueryIntent
	var err error

	hasFilters := req.Level != "" || req.Service != "" || req.Host != "" || req.Since != ""

	if !hasFilters && req.Query != "" {
		intent, err = s.parser.Parse(ctx, req.Query)
		if err != nil {
			return nil, err
		}
	} else {
		intent = &nlq.QueryIntent{
			Query:   req.Query,
			Level:   req.Level,
			Service: req.Service,
			Host:    req.Host,
			Since:   req.Since,
		}
	}

	filter := buildLogFilter(intent)

	var logs []*pb.Log

	if intent.Query != "" && intent.Query != " " {
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
