package config_validator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// Hint: schema generated without `omitempty` tag in struct fields

var ContentConfigurationSchema = []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration","$defs":{"LuigiConfigData":{"properties":{"nodes":{"items":{"$ref":"#/$defs/Node"},"type":"array"}},"additionalProperties":false,"type":"object","required":["nodes"]},"LuigiConfigFragment":{"properties":{"data":{"$ref":"#/$defs/LuigiConfigData"}},"additionalProperties":false,"type":"object","required":["data"]},"Node":{"properties":{"entityType":{"type":"string"},"pathSegment":{"type":"string"},"label":{"type":"string"},"icon":{"type":"string"}},"additionalProperties":false,"type":"object","required":["entityType","pathSegment","label","icon"]}},"properties":{"name":{"type":"string"},"luigiConfigFragment":{"items":{"$ref":"#/$defs/LuigiConfigFragment"},"type":"array"}},"additionalProperties":false,"type":"object","required":["name","luigiConfigFragment"]}`)

type contentConfiguration struct{}

func NewContentConfiguration() ContentConfigurationInterface {
	return &contentConfiguration{}
}

func (cC *contentConfiguration) Validate(input []byte, contentType string) (string, error) {
	switch contentType {
	case "json":
		return validateJSON(input)
	case "yaml":
		return validateYAML(input)
	default:

		return "", errors.New("no validator found for content type")
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
	// Marshal the input to JSON
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "", errors.New("error marshaling input to JSON")
	}

	// Load the schema and JSON data for validation
	schemaLoader := gojsonschema.NewBytesLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	// Validate the JSON data against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return "", errors.New("error validating JSON data")
	}

	if !result.Valid() {
		var errors1 []string
		for _, desc := range result.Errors() {
			errors1 = append(errors1, desc.String())
		}
		return "", errors.New(fmt.Sprintf("The document is not valid:\n%s", fmt.Sprint(errors1)))
	}

	return string(jsonBytes), nil
}
