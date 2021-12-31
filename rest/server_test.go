package rest

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnaryCall(t *testing.T) {
	server := httptest.NewServer(NewAPIHandler())
	defer server.Close()
	client := NewAPIClient(server.URL)

	res, err := client.SayHello(context.TODO(), "World")
	require.NoError(t, err)

	assert.Equal(t, "Hello World", res)
}

func TestStreamCall(t *testing.T) {
	server := httptest.NewServer(NewAPIHandler())
	defer server.Close()
	client := NewAPIClient(server.URL)

	res, err := client.StreamHello(context.TODO(), "World", 2)
	require.NoError(t, err)

	var responses []string
	for next := range res {
		responses = append(responses, next)
	}

	assert.ElementsMatch(t, []string{"Hello World 0", "Hello World 1"}, responses)
}
