package httpclient

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// TestService_Do tests the Do method of the service.
func TestService_Do(t *testing.T) {
	service := NewService()

	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    io.Reader
		mockResponse   string
		mockStatusCode int
		mockError      error
		expectedBody   string
		expectError    bool
	}{
		{
			name:           "successful GET request",
			method:         http.MethodGet,
			url:            "https://example.com/success",
			mockResponse:   `{"message": "success"}`,
			mockStatusCode: http.StatusOK,
			expectedBody:   `{"message": "success"}`,
			expectError:    false,
		},
		{
			name:           "non-200 status code",
			method:         http.MethodGet,
			url:            "https://example.com/error",
			mockResponse:   `{"message": "error"}`,
			mockStatusCode: http.StatusInternalServerError,
			expectedBody:   "",
			expectError:    true,
		},
		{
			name:         "request creation error",
			method:       "INVALID_METHOD",
			url:          "https://example.com/invalid",
			expectedBody: "",
			expectError:  true,
		},
		{
			name:           "successful POST request with body",
			method:         http.MethodPost,
			url:            "https://example.com/post",
			requestBody:    strings.NewReader(`{"key": "value"}`),
			mockResponse:   `{"status": "created"}`,
			mockStatusCode: http.StatusOK,
			expectedBody:   `{"status": "created"}`,
			expectError:    false,
		},
		{
			name:         "network error",
			method:       http.MethodGet,
			url:          "https://example.com/network-error",
			mockError:    errors.New("network error"),
			expectedBody: "",
			expectError:  true,
		},
		{
			name:           "read body error",
			method:         http.MethodGet,
			url:            "https://example.com/read-error",
			mockResponse:   string([]byte{0x81, 0xff, 0xff, 0xff}),
			mockStatusCode: http.StatusOK,
			expectedBody:   "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tt.mockError != nil {
				httpmock.RegisterResponder(tt.method, tt.url,
					httpmock.NewErrorResponder(tt.mockError))
			} else {
				httpmock.RegisterResponder(tt.method, tt.url,
					httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse))
			}

			body, err := service.Do(tt.method, tt.url, tt.requestBody)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}
