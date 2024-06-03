package validation

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

const (
	ErrorEmptyInput       = "empty input provided"
	ErrorNoValidator      = "no validator found for content type"
	ErrorMarshalJSON      = "error marshaling input to JSON"
	ErrorValidatingJSON   = "error validating JSON data"
	ErrorDocumentInvalid  = "The document is not valid:\n%s"
	ErrorRequiredField    = "field '%s' is required"
	ErrorInvalidFieldType = "field '%s' is invalid, got '%s', expected '%s'"
)

var ContentConfigurationSchema = []byte(`{
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
}`)

type contentConfiguration struct{}

func NewContentConfiguration() ContentConfigurationInterface {
	return &contentConfiguration{}
}

func (cC *contentConfiguration) Validate(input []byte, contentType string) (string, error) {
	if len(input) == 0 {
		return "", errors.New(ErrorEmptyInput)
	}

	switch contentType {
	case "json":
		return validateJSON(input)
	case "yaml":
		return validateYAML(input)
	default:

		return "", errors.New(ErrorNoValidator)
	}
}

func validateJSON(input []byte) (string, error) {
	var config ContentConfiguration
	if err := json.Unmarshal(input, &config); err != nil {
		return "", err
	}
	return validateSchema(config, ContentConfigurationSchema)
}

func validateYAML(input []byte) (string, error) {
	var config ContentConfiguration
	if err := yaml.Unmarshal(input, &config); err != nil {
		return "", err
	}

	return validateSchema(config, ContentConfigurationSchema)
}

func validateSchema(input ContentConfiguration, schema []byte) (string, error) {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "", errors.New(ErrorMarshalJSON)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return "", errors.New(ErrorValidatingJSON)
	}

	if !result.Valid() {
		var errorsAccumulator []string
		for _, desc := range result.Errors() {
			switch desc.Type() {
			case "required":
				errorsAccumulator = append(errorsAccumulator, fmt.Sprintf(ErrorRequiredField, desc.Field()))
			case "invalid_type":
				errorsAccumulator = append(errorsAccumulator, fmt.Sprintf(ErrorInvalidFieldType, desc.Field(), desc.Details()["type"], desc.Details()["expected"]))
			default:
				errorsAccumulator = append(errorsAccumulator, desc.String())
			}
		}
		return "", errors.Errorf(ErrorDocumentInvalid, fmt.Sprint(errorsAccumulator))
	}

	return string(jsonBytes), nil
}
