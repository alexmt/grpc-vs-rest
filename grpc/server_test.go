package grpc

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func runTestServer(ctx context.Context) (MyServiceClient, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}
	server := NewGRPCServer()
	go func() {
		if err := server.Serve(ln); err != nil {
			panic(err)
		}
	}()
	port := ln.Addr().(*net.TCPAddr).Port
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	return NewMyServiceClient(conn), nil
}

func TestUnaryCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := runTestServer(ctx)
	require.NoError(t, err)

	resp, err := client.SayHello(ctx, &HelloRequest{Name: "world"})
	require.NoError(t, err)

	assert.Equal(t, "Hello world", resp.Message)
}

func TestStreamCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := runTestServer(ctx)
	require.NoError(t, err)

	resp, err := client.StreamHello(ctx, &StreamHelloRequest{Name: "world", Count: 2})
	require.NoError(t, err)

	var responses []string
	for {
		r, err := resp.Recv()
		if err != nil {
			break
		}
		responses = append(responses, r.Message)
	}

	assert.ElementsMatch(t, []string{"Hello world 0", "Hello world 1"}, responses)
}
