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

type contentConfiguration struct{}

func NewContentConfiguration() ContentConfigurationInterface {
	return &contentConfiguration{}
}

func (cC *contentConfiguration) Validate(schema, input []byte, contentType string) (string, error) {
	if len(input) == 0 {
		return "", errors.New(ErrorEmptyInput)
	}

	switch contentType {
	case "json":
		return validateJSON(schema, input)
	case "yaml":
		return validateYAML(schema, input)
	default:

		return "", errors.New(ErrorNoValidator)
	}
}

func validateJSON(schema, input []byte) (string, error) {
	var config ContentConfiguration
	if err := json.Unmarshal(input, &config); err != nil {
		return "", err
	}
	return validateSchema(schema, config)
}

func validateYAML(schema, input []byte) (string, error) {
	var config ContentConfiguration
	if err := yaml.Unmarshal(input, &config); err != nil {
		return "", err
	}
	return validateSchema(schema, config)
}

// func validateSchema(schema []byte, input ContentConfiguration) (string, error) {
func validateSchema(schema []byte, input interface{}) (string, error) {
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
				errorsAccumulator = append(errorsAccumulator, fmt.Sprintf(
					ErrorInvalidFieldType,
					desc.Field(),
					desc.Details()["type"],
					desc.Details()["expected"]))
			default:
				errorsAccumulator = append(errorsAccumulator, desc.String())
			}
		}
		return "", errors.Errorf(ErrorDocumentInvalid, fmt.Sprint(errorsAccumulator))
	}

	return string(jsonBytes), nil
}
