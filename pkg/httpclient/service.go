package httpclient

import (
	"fmt"
	"io"
	"net/http"
)

// Service interface to abstract the Do method for making HTTP requests.
type Service interface {
	Do(method, url string, body io.Reader) (res []byte, err error, retry bool)
}

// service to handle HTTP requests.
type service struct {
	client *http.Client
}

// NewService creates a new service with the given HTTP client.
func NewService() *service {
	return &service{
		client: &http.Client{},
	}
}

// Do makes an HTTP request to the specified URL.
func (s *service) Do(method, url string, requestBody io.Reader) (res []byte, err error, retry bool) {
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err), false
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err), false
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode != http.StatusOK {
		// Give the caller signal to retry if we have 5xx status codes
		if resp.StatusCode >= http.StatusInternalServerError && resp.StatusCode < 600 {
			return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode), true
		}

		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode), false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err), false
	}

	return body, nil, false
}
