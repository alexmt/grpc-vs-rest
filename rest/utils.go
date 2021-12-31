package rest

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// JSONHandler returns a handler that writes the result of action as a JSON
func JSONHandler(action func(*http.Request) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []byte
		resp, err := action(r)
		if err == nil {
			data, err = json.Marshal(resp)
		}
		if err != nil {
			_, _ = fmt.Fprintf(w, "Internal error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(data)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
		}
	}
}

// SSEJSONHandler returns a handler that writes the result of action as a Server-Sent Event
func SSEJSONHandler(action func(*http.Request) (func() (interface{}, bool), error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := action(r)
		if err != nil {
			_, _ = fmt.Fprintf(w, "Internal error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		for next, done := res(); !done; next, done = res() {
			data, err := json.Marshal(next)
			if err != nil {
				_, _ = fmt.Fprintf(w, "Internal error: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, _ = fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
}

// ParseSSE parses SSE data stream and executes callback for each server event
func ParseSSE(ctx context.Context, r io.Reader, callback func(event string, data []byte) error) error {
	br := bufio.NewReader(r)

	delim := []byte{':'}

	var currentEvent string

	for ctx.Err() == nil {
		bs, err := br.ReadBytes('\n')

		if err != nil && err != io.EOF {
			return err
		}

		if len(bs) < 2 {
			if err == io.EOF {
				break
			}
			continue
		}

		spl := bytes.SplitN(bs, delim, 2)

		if len(spl) < 2 {
			continue
		}

		switch string(spl[0]) {
		case "event":
			currentEvent = string(bytes.TrimSpace(spl[1]))
		case "data":
			if err := callback(currentEvent, bytes.TrimSpace(spl[1])); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}
