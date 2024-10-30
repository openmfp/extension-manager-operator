package server

import (
	"encoding/json"
	"net/http"

	"github.com/openmfp/golang-commons/logger"

	"github.com/openmfp/extension-content-operator/pkg/validation"
)

type requestValidate struct {
	ContentConfiguration string `json:"contentConfiguration"`
}

type responseValid struct {
	ParsedConfiguration string `json:"parsedConfiguration"`
}

type responseError struct {
	ValidationErrors []validationError `json:"validationErrors"`
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

func (h *HttpValidateHandler) HandlerValidate(w http.ResponseWriter, r *http.Request) {
	var request requestValidate
	var rValid responseValid

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&request)

	if err != nil {
		h.log.Error().Err(err).Msg("Reading request body failed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	parsedConfig, err := h.validator.Validate([]byte(request.ContentConfiguration), "json")
	if err != nil {
		var responseErr responseError
		responseErr.ValidationErrors = []validationError{{
			Message: err.Error(),
		}}
		w.WriteHeader(http.StatusOK)
		responseBytes, err := json.Marshal(responseErr)
		if err != nil {
			h.log.Error().Err(err).Msg("Marshalling response failed")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Marshalling response failed"))
			return
		}
		w.Write(responseBytes)
		return
	}

	rValid.ParsedConfiguration = parsedConfig
	responseBytes, err := json.Marshal(rValid)
	if err != nil {
		h.log.Error().Err(err).Msg("Marshalling response failed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Marshalling response failed"))
		return
	}
	w.Write(responseBytes)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
