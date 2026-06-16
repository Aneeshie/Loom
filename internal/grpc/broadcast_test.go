package grpc

import (
	"testing"

	pb "github.com/Aneeshie/loom/proto"
)

func TestBroadcast(t *testing.T) {
	tests := []struct {
		name     string
		filter   *pb.LogFilter
		log      *pb.Log
		received bool
	}{
		{
			name:     "nil filter receives",
			filter:   nil,
			log:      &pb.Log{Level: "INFO"},
			received: true,
		},
		{
			name: "matching level receives",
			filter: &pb.LogFilter{
				Level: strPtr("INFO"),
			},
			log: &pb.Log{
				Level: "INFO",
			},
			received: true,
		},
		{
			name: "mismatching level does not receive",
			filter: &pb.LogFilter{
				Level: strPtr("ERROR"),
			},
			log: &pb.Log{
				Level: "INFO",
			},
			received: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sub := &Subscriber{
				ch:     make(chan *pb.Log, 1),
				filter: tt.filter,
			}

			service := &LogService{
				subscribers: map[*Subscriber]struct{}{
					sub: {},
				},
			}

			service.broadcast(tt.log)

			select {
			case <-sub.ch:
				if !tt.received {
					t.Fatal("unexpected log received")
				}
			default:
				if tt.received {
					t.Fatal("expected log to be received")
				}
			}
		})
	}
}
