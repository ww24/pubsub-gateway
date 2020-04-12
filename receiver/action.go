package receiver

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	requestTimeout = 20 * time.Second
)

var (
	httpClient = &http.Client{
		Timeout: requestTimeout,
	}
)

// Executable provides Exec method for action.
type Executable interface {
	Exec(context.Context, []byte) error
}

// HTTPAction implements action for HTTP.
type HTTPAction struct {
	header http.Header
	method string
	url    string
}

// NewHTTPAction returns a new http action.
func NewHTTPAction(header http.Header, method, url string) *HTTPAction {
	return &HTTPAction{
		header: header,
		method: method,
		url:    url,
	}
}

// Exec executes pubsub action.
func (a *HTTPAction) Exec(ctx context.Context, payload []byte) error {
	var body io.Reader
	if len(payload) > 0 {
		body = bytes.NewBuffer(payload)
	}
	req, err := http.NewRequest(a.method, a.url, body)
	if err != nil {
		return err
	}
	req.Header = a.header
	resp, err := httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	log.Printf("Sent, status: %v\n", resp.Status)
	return nil
}
