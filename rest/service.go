package rest

import (
	"context"
	"fmt"
)

// MyService defines methods exposed by the service.
type MyService interface {
	SayHello(ctx context.Context, name string) (string, error)
	StreamHello(ctx context.Context, name string, count int) (chan string, error)
}

// myServiceImpl is an implementation of MyService.
type myServiceImpl struct{}

var _ MyService = &myServiceImpl{}

// SayHello is a simple method that returns a string.
func (m *myServiceImpl) SayHello(_ context.Context, name string) (string, error) {
	return "Hello " + name, nil
}

// StreamHello is a streaming method that returns a channel of strings.
func (m *myServiceImpl) StreamHello(_ context.Context, name string, count int) (chan string, error) {
	ch := make(chan string)
	go func() {
		for i := 0; i < count; i++ {
			ch <- fmt.Sprintf("Hello %s %d", name, i)
		}
		close(ch)
	}()
	return ch, nil
}
