package service

import (
	"time"

	"github.com/Aneeshie/loom/internal/nlq"
	pb "github.com/Aneeshie/loom/proto"
)

func buildLogFilter(intent *nlq.QueryIntent) *pb.LogFilter {
	filter := &pb.LogFilter{}

	if intent.Level != "" {
		filter.Level = &intent.Level
	}

	if intent.Service != "" {
		filter.ServiceName = &intent.Service
	}

	if intent.Host != "" {
		filter.Host = &intent.Host
	}

	if intent.Since != "" {
		duration, err := time.ParseDuration(intent.Since)
		if err == nil {
			start := time.Now().Add(-duration).Unix()
			filter.StartTime = &start
		}
	}

	return filter
}
