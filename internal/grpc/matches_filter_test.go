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
		{
			name: "search match",
			log: &pb.Log{
				Message: "database connection established",
			},
			filter: &pb.LogFilter{
				Search: strPtr("database"),
			},
			want: true,
		},
		{
			name: "search mismatch",
			log: &pb.Log{
				Message: "user logged in",
			},
			filter: &pb.LogFilter{
				Search: strPtr("database"),
			},
			want: false,
		},
		{
			name: "search case insensitive",
			log: &pb.Log{
				Message: "Database connection established",
			},
			filter: &pb.LogFilter{
				Search: strPtr("database"),
			},
			want: true,
		},
		{
			name: "search and level match",
			log: &pb.Log{
				Message: "database timeout",
				Level:   "ERROR",
			},
			filter: &pb.LogFilter{
				Search: strPtr("database"),
				Level:  strPtr("ERROR"),
			},
			want: true,
		},
		{
			name: "search matches but level mismatch",
			log: &pb.Log{
				Message: "database timeout",
				Level:   "INFO",
			},
			filter: &pb.LogFilter{
				Search: strPtr("database"),
				Level:  strPtr("ERROR"),
			},
			want: false,
		},
		{
			name: "search matches but service mismatch",
			log: &pb.Log{
				Message:     "database timeout",
				ServiceName: "auth",
			},
			filter: &pb.LogFilter{
				Search:      strPtr("database"),
				ServiceName: strPtr("billing"),
			},
			want: false,
		},
		{
			name: "search substring match",
			log: &pb.Log{
				Message: "database timeout while connecting",
			},
			filter: &pb.LogFilter{
				Search: strPtr("timeout"),
			},
			want: true,
		},
		{
			name: "empty search string matches",
			log: &pb.Log{
				Message: "database timeout",
			},
			filter: &pb.LogFilter{
				Search: strPtr(""),
			},
			want: true,
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
