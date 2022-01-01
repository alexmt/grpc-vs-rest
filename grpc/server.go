package grpc

import "google.golang.org/grpc"

func NewGRPCServer() *grpc.Server {
	server := grpc.NewServer()
	RegisterMyServiceServer(server, &myServiceImpl{})
	return server
}
