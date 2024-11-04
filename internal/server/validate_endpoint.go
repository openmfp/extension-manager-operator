package server

import (
	"encoding/json"
	"net/http"

	"github.com/openmfp/golang-commons/logger"

	"github.com/openmfp/extension-content-operator/pkg/validation"

	"github.com/openmfp/golang-commons/sentry"
)

type requestValidate struct {
	ContentType          string `json:"contentType,omitempty"`
	ContentConfiguration string `json:"contentConfiguration"`
}

type Response struct {
	ParsedConfiguration string            `json:"parsedConfiguration,omitempty"`
	ValidationErrors    []validationError `json:"validationErrors,omitempty"`
}

type validationError struct {
	Message string `json:"message"`
}

func NewHttpValidateHandler(log *logger.Logger, validator validation.ExtensionConfiguration) *HttpValidateHandler {
	return &HttpValidateHandler{
		validator: validator,
		log:       log,
	}
}

type HttpValidateHandler struct {
	validator validation.ExtensionConfiguration
	log       *logger.Logger
}

func (h *HttpValidateHandler) HandlerHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK")) // nolint: errcheck
}

func (h *HttpValidateHandler) HandlerValidate(w http.ResponseWriter, r *http.Request) {
	// decode request
	var request requestValidate
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)
	defer r.Body.Close()
	if err != nil {
		h.log.Error().Err(err).Msg("Reading request body failed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error())) // nolint: errcheck
		sentry.CaptureError(err, sentry.Tags{"error": "Writing response failed"}, sentry.Extras{"data": r.Body})
		return
	}

	// validation
	parsedConfig, merr := h.validator.Validate([]byte(request.ContentConfiguration), request.ContentType)
	if merr.Len() > 0 {
		var responseErr Response
		for _, e := range merr.Errors {
			responseErr.ValidationErrors = append(responseErr.ValidationErrors, validationError{
				Message: e.Error(),
			})
		}

		responseBytes, _ := json.Marshal(responseErr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes) // nolint: errcheck
		return
	}

	// send response
	var rValid Response
	rValid.ParsedConfiguration = parsedConfig
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&rValid)
}
