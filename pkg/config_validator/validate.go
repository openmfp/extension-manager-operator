package config_validator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const (
	ERROR_EMPTY_INPUT        = "empty input provided"
	ERROR_NO_VALIDATOR       = "no validator found for content type"
	ERROR_MARSHAL_JSON       = "error marshaling input to JSON"
	ERROR_VALIDATING_JSON    = "error validating JSON data"
	ERROR_DOCUMENT_INVALID   = "The document is not valid:\n%s"
	ERROR_REQUIRED_FIELD     = "field '%s' is required"
	ERROR_INVALID_FIELD_TYPE = "field '%s' is invalid, got '%s', expected '%s'"
)

var ContentConfigurationSchema = []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration","$defs":{"LuigiConfigData":{"properties":{"nodes":{"items":{"$ref":"#/$defs/Node"},"type":"array"}},"additionalProperties":false,"type":"object","required":["nodes"]},"LuigiConfigFragment":{"properties":{"data":{"$ref":"#/$defs/LuigiConfigData"}},"additionalProperties":false,"type":"object","required":["data"]},"Node":{"properties":{"entityType":{"type":"string"},"pathSegment":{"type":"string"},"label":{"type":"string"},"icon":{"type":"string"}},"additionalProperties":false,"type":"object","required":["entityType","pathSegment","label","icon"]}},"properties":{"name":{"type":"string"},"luigiConfigFragment":{"items":{"$ref":"#/$defs/LuigiConfigFragment"},"type":"array"}},"additionalProperties":false,"type":"object","required":["name","luigiConfigFragment"]}`)

type contentConfiguration struct{}

func NewContentConfiguration() ContentConfigurationInterface {
	return &contentConfiguration{}
}

func (cC *contentConfiguration) Validate(input []byte, contentType string) (string, error) {
	if len(input) == 0 {
		return "", errors.New(ERROR_EMPTY_INPUT)
	}

	switch contentType {
	case "json":
		return validateJSON(input)
	case "yaml":
		return validateYAML(input)
	default:

		return "", errors.New(ERROR_NO_VALIDATOR)
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
		return "", errors.New(ERROR_MARSHAL_JSON)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return "", errors.New(ERROR_VALIDATING_JSON)
	}

	if !result.Valid() {
		var errorsAccumulator []string
		for _, desc := range result.Errors() {
			// Customize error messages based on schema validation errors
			switch desc.Type() {
			case "required":
				errorsAccumulator = append(errorsAccumulator, fmt.Sprintf(ERROR_REQUIRED_FIELD, desc.Field()))
			case "invalid_type":
				errorsAccumulator = append(errorsAccumulator, fmt.Sprintf(ERROR_INVALID_FIELD_TYPE, desc.Field(), desc.Details()["type"], desc.Details()["expected"]))
			default:
				errorsAccumulator = append(errorsAccumulator, desc.String())
			}
		}
		return "", errors.New(fmt.Sprintf(ERROR_DOCUMENT_INVALID, fmt.Sprint(errorsAccumulator)))
	}

	return string(jsonBytes), nil
}
