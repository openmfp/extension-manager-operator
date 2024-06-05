package subroutines

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/openmfp/extension-content-operator/pkg/validation"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/errors"
	"github.com/openmfp/golang-commons/logger"
)

const (
	ContentConfigurationSubroutineName = "ContentConfigurationSubroutine"
)

type ContentConfigurationSubroutine struct {
	client    *http.Client
	validator validation.ContentConfigurationInterface
}

func NewContentConfigurationSubroutine() *ContentConfigurationSubroutine {
	return &ContentConfigurationSubroutine{
		client:    http.DefaultClient,
		validator: validation.NewContentConfiguration(),
	}
}

func (r *ContentConfigurationSubroutine) WithClient(client *http.Client) {
	r.client = client
}

func (r *ContentConfigurationSubroutine) WithValidator(validator validation.ContentConfigurationInterface) {
	r.validator = validator
}

func (r *ContentConfigurationSubroutine) GetName() string {
	return ContentConfigurationSubroutineName
}

func (r *ContentConfigurationSubroutine) Finalize(
	ctx context.Context,
	runtimeObj lifecycle.RuntimeObject) (ctrl.Result, errors.OperatorError) {
	return ctrl.Result{}, nil
}

func (r *ContentConfigurationSubroutine) Finalizers() []string {
	return []string{}
}

func (r *ContentConfigurationSubroutine) Process(
	ctx context.Context, runtimeObj lifecycle.RuntimeObject,
) (ctrl.Result, errors.OperatorError) {
	log := logger.LoadLoggerFromContext(ctx)

	instance := runtimeObj.(*v1alpha1.ContentConfiguration)

	var contentType string
	var rawConfig []byte
	// InlineConfiguration has higher priority than RemoteConfiguration
	switch {
	case instance.Spec.InlineConfiguration.Content != "":
		contentType = instance.Spec.InlineConfiguration.ContentType
		rawConfig = []byte(instance.Spec.InlineConfiguration.Content)
	case instance.Spec.RemoteConfiguration.URL != "":
		bytes, err, retry := r.getRemoteConfig(instance.Spec.RemoteConfiguration.URL)
		if err != nil {
			log.Err(err).Msg("failed to fetch remote configuration")

			return ctrl.Result{}, errors.NewOperatorError(err, retry, true)
		}
		contentType = instance.Spec.RemoteConfiguration.ContentType
		rawConfig = bytes
	default:
		return ctrl.Result{}, errors.NewOperatorError(errors.New("no configuration provided"), false, true)
	}

	validatedConfig, err := r.validator.Validate(getSchema(), rawConfig, contentType)
	if err != nil {
		log.Err(err).Msg("failed to validate configuration")

		return ctrl.Result{}, errors.NewOperatorError(err, false, true)
	}

	instance.Status.ConfigurationResult = validatedConfig

	return ctrl.Result{}, nil
}

// Do makes an HTTP request to the specified URL.
func (r *ContentConfigurationSubroutine) getRemoteConfig(url string) (res []byte, err error, retry bool) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err), false
	}

	resp, err := r.client.Do(req)
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

	// TODO
	// we need to check the size of the received body before loading it to memory.
	// In case it exceeds a certain size we should reject it.
	// https://github.com/openmfp/extension-content-operator/pull/23#discussion_r1622598363

	return body, nil, false
}

func getSchema() []byte {
	s := `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://github.com/openmfp/extension-content-operator/pkg/validation/content-configuration",
  "$defs": {
    "LuigiConfigData": {
      "properties": {
        "nodes": {
          "items": {
            "$ref": "#/$defs/Node"
          },
          "type": "array"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": ["nodes"]
    },
    "LuigiConfigFragment": {
      "properties": {
        "data": {
          "$ref": "#/$defs/LuigiConfigData"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": ["data"]
    },
    "Node": {
      "properties": {
        "entityType": {
          "type": "string"
        },
        "pathSegment": {
          "type": "string"
        },
        "label": {
          "type": "string"
        },
        "icon": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": ["entityType", "pathSegment", "label", "icon"]
    }
  },
  "properties": {
    "name": {
      "type": "string"
    },
    "luigiConfigFragment": {
      "items": {
        "$ref": "#/$defs/LuigiConfigFragment"
      },
      "type": "array"
    }
  },
  "additionalProperties": false,
  "type": "object",
  "required": ["name", "luigiConfigFragment"]
}`

	return []byte(s)
}
