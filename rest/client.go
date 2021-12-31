package rest

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/go-resty/resty/v2"
)

var _ MyService = &myServiceClient{}

type myServiceClient struct {
	client *resty.Client
}

func NewAPIClient(baseURL string) *myServiceClient {
	return &myServiceClient{
		client: resty.New().SetBaseURL(baseURL),
	}
}

func (c *myServiceClient) SayHello(ctx context.Context, name string) (string, error) {
	var res string
	_, err := c.client.R().
		SetContext(ctx).
		SetResult(&res).
		SetQueryParam("name", name).
		ForceContentType("application/json").
		Get("/say-hello")
	return res, err
}

func (c *myServiceClient) StreamHello(ctx context.Context, name string, count int) (chan string, error) {
	r, err := c.client.R().
		SetContext(ctx).
		SetDoNotParseResponse(true).
		SetQueryParam("name", name).
		SetQueryParam("count", strconv.FormatInt(int64(count), 10)).
		Get("/stream-hello")
	if err != nil {
		return nil, err
	}
	ch := make(chan string)
	go func() {
		defer r.RawBody().Close()
		defer close(ch)

		_ = ParseSSE(ctx, r.RawBody(), func(event string, data []byte) error {
			var res string
			if err := json.Unmarshal(data, &res); err != nil {
				return err
			}
			ch <- res
			return nil
		})
	}()
	return ch, nil
}
