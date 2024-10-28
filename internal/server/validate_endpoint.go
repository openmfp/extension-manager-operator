package server

import (
	"encoding/json"
	"net/http"

	"github.com/openmfp/golang-commons/logger"

	"github.com/openmfp/extension-content-operator/pkg/validation"
)

type requestValidate struct {
	ConfigurationContentType string `json:"configurationContentType"` // Enum: "json" or "yaml"
	ContentConfiguration     string `json:"contentConfiguration"`
}

type responseValidate struct {
	ParsedConfiguration string            `json:"parsedConfiguration"`
	ValidationErrors    []validationError `json:"validationErrors"`
}

type validationError struct {
	Message string `json:"message"`
}

type validationHandler struct {
	validator validation.ExtensionConfiguration
	log       *logger.Logger
}

func (h *validationHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request requestValidate
		var response responseValidate

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if request.ConfigurationContentType == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing or empty configurationContentType value in request body"))
			return
		}

		if request.ConfigurationContentType != "json" && request.ConfigurationContentType != "yaml" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unexpected configuration content type"))
			return
		}

		if request.ContentConfiguration == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing or empty contentConfiguration value in request body"))
			return
		}

		parsedConfig, err := h.validator.Validate([]byte(request.ContentConfiguration), request.ConfigurationContentType)
		if err != nil {
			response.ValidationErrors = []validationError{{
				Message: err.Error(),
			}}
		} else {
			response.ParsedConfiguration = parsedConfig
		}

		rspBytes, err := json.Marshal(&response)
		if err != nil {
			h.log.Error().Err(err).Msg("Marshaling json failed")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(rspBytes)
	}
}
