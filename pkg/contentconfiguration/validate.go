package contentconfiguration

import (
	"encoding/json"
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

func (cC *contentConfiguration) ValidateJSON(input []byte) (bool, []string) {
	var config ContentConfiguration
	if err := json.Unmarshal(input, &config); err != nil {
		return false, []string{err.Error()}
	}
	return validateSchema(input, ContentConfigurationSchema)
}

func (cC *contentConfiguration) ValidateYAML(input []byte) (bool, []string) {
	var config ContentConfiguration
	if err := yaml.Unmarshal(input, &config); err != nil {
		return false, []string{err.Error()}
	}

	return validateSchema(input, ContentConfigurationSchema)

}

func validateSchema(input []byte, schema string) (bool, []string) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(input)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return false, []string{err.Error()}
	}

	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		return false, errors
	}

	return true, nil
}
