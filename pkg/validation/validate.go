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

func (cC *contentConfiguration) WithSchema(schema []byte) error {
	if len(schema) == 0 {
		return ErrorEmptyInput
	}
	cC.schema = schema
	return nil
}

// Validate does the validation against strict schema
// and includes any optional field that is not defined there:
// Steps:
// 1. Gets raw JSON or YAML data from the input
// 2. In case of YAML it converts data to raw JSON
// 3. Passes the raw JSON to the schema validator
// 4. In case of success returns the original raw JSON
func (cC *contentConfiguration) Validate(input []byte, contentType string) (string, error) {
	if len(input) == 0 {
		return "", ErrorEmptyInput
	}

	var rawJSON []byte
	var err error
	switch contentType {
	case "json":
		rawJSON = input
	case "yaml":
		rawJSON, err = convertYAMLToJSON(input)
		if err != nil {
			return "", err
		}
	default:
		return "", ErrorNoValidator
	}

	return validateJSON(cC.schema, rawJSON)
}

func validateJSON(schema, input []byte) (string, error) {
	var config ContentConfiguration
	if err := json.Unmarshal(input, &config); err != nil {
		return "", err
	}

	if err := validateSchema(schema, config); err != nil {
		return "", err
	}

	return string(input), nil
}

// func validateSchema(schema []byte, input ContentConfiguration) (string, error) {
func validateSchema(schema []byte, input interface{}) error {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return ErrorMarshalJSON
	}

	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return ErrorValidatingJSON
	}

	if !result.Valid() {
		var errorsAccumulator []string
		for _, desc := range result.Errors() {
			switch desc.Type() {
			case "required":
				errorsAccumulator = append(errorsAccumulator, desc.Description())
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
		return errors.Errorf(ErrorDocumentInvalid.Error(), fmt.Sprint(errorsAccumulator))
	}

	return nil
}

// ConvertYAMLToJSON converts a YAML byte array to a JSON byte array
func convertYAMLToJSON(yamlData []byte) ([]byte, error) {
	// Unmarshal YAML into a map
	var data map[string]interface{}
	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML: %w", err)
	}

	// Marshal the map into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %w", err)
	}

	return jsonData, nil
}
