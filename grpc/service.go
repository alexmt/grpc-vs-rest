package grpc

import (
	"context"
	"fmt"
)

type myServiceImpl struct {
	UnimplementedMyServiceServer
}

var _ MyServiceServer = &myServiceImpl{}

func (m *myServiceImpl) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello " + req.Name}, nil
}

func (m *myServiceImpl) StreamHello(req *StreamHelloRequest, stream MyService_StreamHelloServer) error {
	for i := int32(0); i < req.Count; i++ {
		if err := stream.Send(&HelloReply{Message: fmt.Sprintf("Hello %s %d", req.Name, i)}); err != nil {
			return err
		}
	}
	return nil
}
