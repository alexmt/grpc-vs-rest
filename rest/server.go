package rest

import (
	"net/http"
	"strconv"
)

// NewAPIHandler returns a new http.Handler that handles API requests.
func NewAPIHandler() http.Handler {
	mux := http.NewServeMux()
	svc := &myServiceImpl{}
	mux.HandleFunc("/say-hello", JSONHandler(func(r *http.Request) (interface{}, error) {
		return svc.SayHello(r.Context(), r.URL.Query().Get("name"))
	}))
	mux.HandleFunc("/stream-hello", SSEJSONHandler(func(r *http.Request) (func() (interface{}, bool), error) {
		count, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil {
			return nil, err
		}
		res, err := svc.StreamHello(r.Context(), r.URL.Query().Get("name"), count)
		if err != nil {
			return nil, err
		}
		return func() (interface{}, bool) {
			next := <-res
			return next, next == ""
		}, nil
	}))
	return mux
}
