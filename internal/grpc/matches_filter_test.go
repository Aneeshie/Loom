package grpc

import (
	"testing"

	pb "github.com/Aneeshie/loom/proto"
)

func strPtr(s string) *string {
	return &s
}

func TestMatchesFilter(t *testing.T) {
	tests := []struct {
		name   string
		log    *pb.Log
		filter *pb.LogFilter
		want   bool
	}{
		{
			name: "nil filter matches everything",
			log: &pb.Log{
				Level:       "INFO",
				ServiceName: "api",
				Host:        "localhost",
			},
			filter: nil,
			want:   true,
		},
		{
			name: "level matches",
			log: &pb.Log{
				Level: "INFO",
			},
			filter: &pb.LogFilter{
				Level: strPtr("INFO"),
			},
			want: true,
		},
		{
			name: "level mismatch",
			log: &pb.Log{
				Level: "INFO",
			},
			filter: &pb.LogFilter{
				Level: strPtr("ERROR"),
			},
			want: false,
		},
		{
			name: "service matches",
			log: &pb.Log{
				ServiceName: "auth",
			},
			filter: &pb.LogFilter{
				ServiceName: strPtr("auth"),
			},
			want: true,
		},
		{
			name: "multiple filters match",
			log: &pb.Log{
				Level:       "ERROR",
				ServiceName: "auth",
				Host:        "localhost",
			},
			filter: &pb.LogFilter{
				Level:       strPtr("ERROR"),
				ServiceName: strPtr("auth"),
				Host:        strPtr("localhost"),
			},
			want: true,
		},
		{
			name: "multiple filters one mismatch",
			log: &pb.Log{
				Level:       "ERROR",
				ServiceName: "auth",
				Host:        "localhost",
			},
			filter: &pb.LogFilter{
				Level:       strPtr("ERROR"),
				ServiceName: strPtr("billing"),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesFilter(tt.log, tt.filter)

			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
