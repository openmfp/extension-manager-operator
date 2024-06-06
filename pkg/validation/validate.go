package validation

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

var (
	ErrorEmptyInput       = errors.New("empty input provided")
	ErrorNoValidator      = errors.New("no validator found for content type")
	ErrorMarshalJSON      = errors.New("error marshaling input to JSON")
	ErrorValidatingJSON   = errors.New("error validating JSON data")
	ErrorDocumentInvalid  = errors.New("The document is not valid:\n%s")
	ErrorRequiredField    = errors.New("field '%s' is required")
	ErrorInvalidFieldType = errors.New("field '%s' is invalid, got '%s', expected '%s'")
)

type contentConfiguration struct {
	schema []byte
}

//go:embed default_schema_core.openmfp.io_contentconfigurations_gen1.json
var schemaJSON []byte

func NewContentConfiguration() ExtensionConfiguration {

	return &contentConfiguration{
		schema: schemaJSON,
	}
}

func (cC *contentConfiguration) LoadSchema(schema []byte) error {
	if len(schema) == 0 {
		return ErrorEmptyInput
	}
	cC.schema = schema
	return nil
}

func (cC *contentConfiguration) Validate(input []byte, contentType string) (string, error) {
	if len(input) == 0 {
		return "", ErrorEmptyInput
	}

	switch contentType {
	case "json":
		return validateJSON(cC.schema, input)
	case "yaml":
		return validateYAML(cC.schema, input)
	default:

		return "", ErrorNoValidator
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
		return "", ErrorMarshalJSON
	}

	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return "", ErrorValidatingJSON
	}

	if !result.Valid() {
		var errorsAccumulator []string
		for _, desc := range result.Errors() {
			switch desc.Type() {
			case "required":
				errorsAccumulator = append(errorsAccumulator, fmt.Sprintf(ErrorRequiredField.Error(), desc.Field()))
			case "invalid_type":
				errorsAccumulator = append(errorsAccumulator, fmt.Sprintf(
					ErrorInvalidFieldType.Error(),
					desc.Field(),
					desc.Details()["type"],
					desc.Details()["expected"]))
			default:
				errorsAccumulator = append(errorsAccumulator, desc.String())
			}
		}
		return "", errors.Errorf(ErrorDocumentInvalid.Error(), fmt.Sprint(errorsAccumulator))
	}

	return string(jsonBytes), nil
}
