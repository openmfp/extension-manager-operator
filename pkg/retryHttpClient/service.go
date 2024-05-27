package retryHttpClient

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
	"time"
)

// Service defines an interface for making HTTP requests
type Service interface {
	Do(url, method string, rawBody interface{}) ([]byte, error)
}

// service wraps a retryable HTTP client
type service struct {
	httpClient *retryablehttp.Client
}

// New creates a new service with a configured retryable HTTP client
// retryMax specifies the maximum number of retries
// retryWaitMin specifies the minimum wait time between retries
// retryWaitMax specifies the maximum wait time between retries
func New(retryMax int, retryWaitMin, retryWaitMax time.Duration) *service {
	// Initialize a new retryable HTTP client
	httpClient := retryablehttp.NewClient()
	httpClient.RetryMax = retryMax         // Set the maximum number of retries
	httpClient.RetryWaitMin = retryWaitMin // Set the minimum wait time between retries
	httpClient.RetryWaitMax = retryWaitMax // Set the maximum wait time between retries

	// Return a new service instance with the configured HTTP client
	return &service{
		httpClient: httpClient,
	}
}

// Do sends an HTTP request to the specified URL with the given method and body
// It uses the retryable HTTP client to handle retries
// Returns the response body as a byte slice or an error if the request fails
func (s *service) Do(url, method string, rawBody interface{}) ([]byte, error) {
	// Create a new retryable HTTP request
	req, err := retryablehttp.NewRequest(method, url, rawBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Execute the request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Return the response body as a byte slice
	return body, nil
}
