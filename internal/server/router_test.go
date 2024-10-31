package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openmfp/extension-content-operator/internal/config"
	"github.com/openmfp/extension-content-operator/pkg/validation"
	"github.com/openmfp/golang-commons/logger"
)

func initLog() *logger.Logger {
	logConfig := logger.DefaultConfig()
	logConfig.Name = "router_test"
	logConfig.Level = "DEBUG"
	logConfig.NoJSON = false
	log, _ := logger.New(logConfig)
	return log
}

func TestCreateRouter(t *testing.T) {
	tests := []struct {
		name       string
		isLocal    bool
		method     string
		path       string
		expectCode int
		expectCORS bool
	}{
		{
			name:       "validate endpoint exists",
			isLocal:    false,
			method:     http.MethodPost,
			path:       "/validate",
			expectCode: http.StatusInternalServerError,
		},
		{
			name:       "wrong method not allowed",
			isLocal:    false,
			method:     http.MethodGet,
			path:       "/validate",
			expectCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "local setup OK",
			isLocal:    true,
			method:     http.MethodPost,
			path:       "/validate",
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{IsLocal: tt.isLocal}
			log := initLog()

			validator := validation.NewContentConfiguration()

			router := CreateRouter(cfg, log, validator)
			assert.NotNil(t, router)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectCode, rr.Code)

		})
	}
}
