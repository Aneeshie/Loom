package main

import "github.com/Aneeshie/loom/internal/grpc"

func main() {
	server := grpc.NewServer()
	server.Run()
}
