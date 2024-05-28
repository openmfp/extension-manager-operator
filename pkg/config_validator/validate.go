package config_validator

import (
	"encoding/json"
	"errors"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

const ContentConfigurationSchema = `{
 "$schema": "http://json-schema.org/draft-07/schema#",
 "type": "object",
 "properties": {
   "name": {
     "type": "string"
   },
   "luigiConfigFragment": {
     "type": "array",
     "items": {
       "type": "object",
       "properties": {
         "data": {
           "type": "object",
           "properties": {
             "nodes": {
               "type": "array",
               "items": {
                 "type": "object",
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
                 "required": ["entityType", "pathSegment", "label", "icon"]
               }
             }
           },
           "required": ["nodes"]
         }
       },
       "required": ["data"]
     }
   }
 },
 "required": ["name", "luigiConfigFragment"]
}`

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

func validateSchema(input ContentConfiguration, schema string) (string, error) {
	// Marshal the input to JSON
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return "", errors.New("error marshaling input to JSON")
	}

	// Load the schema and JSON data for validation
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)

	// Validate the JSON data against the schema
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return "", errors.New("error validating JSON data")
	}

	// Generate the result string
	//var validationResult string
	if result.Valid() {
		//validationResult = "The document is valid"
		//validationResult = string(jsonBytes)
		//return validationResult, nil
		return string(jsonBytes), nil
		//} else {
		//validationResult = "The document is not valid. See errors:\n"
		//for _, desc := range result.Errors() {
		//	validationResult += fmt.Sprintf("- %s\n", desc)
		//}
		//return "", errors.New("The document is not valid")
	}

	//return validationResult, nil
	return "", errors.New("The document is not valid")
}
