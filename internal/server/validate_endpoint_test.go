package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openmfp/extension-content-operator/pkg/validation"
	"github.com/openmfp/golang-commons/logger"
	"github.com/stretchr/testify/assert"
)

type responseError struct {
	ValidationErrors []validationError `json:"validationErrors,omitempty"`
}

type responseSuccess struct {
	ParsedConfiguration string `json:"parsedConfiguration,omitempty"`
}

func TestHandlerValidate_Error(t *testing.T) {

	logcfg := logger.DefaultConfig()
	log, _ := logger.New(logcfg)

	handler := NewHttpValidateHandler(log, validation.NewContentConfiguration())

	reqBody := `{"contentConfiguration":"{\"luigiConfigFragment2\": {\"data\": {\"nodeDefaults\": {\"entityType\": \"global\",\"isolateView\": true},\"nodes\": [{\"entityType\": \"global\",\"icon\": \"home\",\"label\": \"Overview\",\"pathSegment\": \"home\"}],\"texts\": [{\"locale\": \"de\",\"textDictionary\": {\"hello\": \"Hallo\"}}]}},\"name\": \"overview\"}"}"`
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handler.HandlerValidate(w, req)

	resp := w.Result()

	r := &responseError{}
	err := json.NewDecoder(resp.Body).Decode(r)
	defer resp.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.GreaterOrEqual(t, len(r.ValidationErrors), 1)
}

func TestHandlerValidate_Success(t *testing.T) {

	logcfg := logger.DefaultConfig()
	log, _ := logger.New(logcfg)

	handler := NewHttpValidateHandler(log, validation.NewContentConfiguration())

	reqBody := `{"contentConfiguration":"{\"luigiConfigFragment\": {\"data\": {\"nodeDefaults\": {\"entityType\": \"global\",\"isolateView\": true},\"nodes\": [{\"entityType\": \"global\",\"icon\": \"home\",\"label\": \"Overview\",\"pathSegment\": \"home\"}],\"texts\": [{\"locale\": \"de\",\"textDictionary\": {\"hello\": \"Hallo\"}}]}},\"name\": \"overview\"}"}"`
	req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	handler.HandlerValidate(w, req)

	resp := w.Result()

	r := &responseSuccess{}
	err := json.NewDecoder(resp.Body).Decode(r)
	defer resp.Body.Close()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.GreaterOrEqual(t, len(r.ParsedConfiguration), 0)
}
