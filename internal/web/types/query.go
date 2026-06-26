package types

import pb "github.com/Aneeshie/loom/proto"

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Intent IntentResponse `json:"intent"`
	Logs   []LogResponse  `json:"logs"`
}

type IntentResponse struct {
	Query   string `json:"query"`
	Level   string `json:"level"`
	Service string `json:"service"`
	Host    string `json:"host"`
	Since   string `json:"since"`
}

type LogResponse struct {
	ID          int64  `json:"id"`
	ServiceName string `json:"service_name"`
	Host        string `json:"host"`
	Level       string `json:"level"`
	Message     string `json:"message"`
	Timestamp   int64  `json:"timestamp"`
}

func NewLogResponse(log *pb.Log) LogResponse {
	return LogResponse{
		ID:          log.Id,
		ServiceName: log.ServiceName,
		Host:        log.Host,
		Level:       log.Level,
		Message:     log.Message,
		Timestamp:   log.Timestamp,
	}
}
