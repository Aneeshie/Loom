package grpc

import pb "github.com/Aneeshie/loom/proto"

type Subscriber struct {
	ch chan *pb.Log

	filter *pb.LogFilter
}
