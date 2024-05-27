package contentconfiguration

import (
	"encoding/json"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
)

type ContentConfigurationSchema struct {
	Schema     string `json:"$schema"`
	Type       string `json:"type"`
	Properties struct {
		Name struct {
			Type string `json:"type"`
		} `json:"name"`
		LuigiConfigFragment struct {
			Type  string `json:"type"`
			Items struct {
				Type       string `json:"type"`
				Properties struct {
					Data struct {
						Type       string `json:"type"`
						Properties struct {
							Nodes struct {
								Type  string `json:"type"`
								Items struct {
									Type       string `json:"type"`
									Properties struct {
										EntityType struct {
											Type string `json:"type"`
										} `json:"entityType"`
										PathSegment struct {
											Type string `json:"type"`
										} `json:"pathSegment"`
										Label struct {
											Type string `json:"type"`
										} `json:"label"`
										Icon struct {
											Type string `json:"type"`
										} `json:"icon"`
									} `json:"properties"`
									Required []string `json:"required"`
								} `json:"items"`
							} `json:"nodes"`
						} `json:"properties"`
						Required []string `json:"required"`
					} `json:"data"`
				} `json:"properties"`
				Required []string `json:"required"`
			} `json:"items"`
		} `json:"luigiConfigFragment"`
	} `json:"properties"`
	Required []string `json:"required"`
}

//const ContentConfigurationSchema = `{
//  "$schema": "http://json-schema.org/draft-07/schema#",
//  "type": "object",
//  "properties": {
//    "name": {
//      "type": "string"
//    },
//    "luigiConfigFragment": {
//      "type": "array",
//      "items": {
//        "type": "object",
//        "properties": {
//          "data": {
//            "type": "object",
//            "properties": {
//              "nodes": {
//                "type": "array",
//                "items": {
//                  "type": "object",
//                  "properties": {
//                    "entityType": {
//                      "type": "string"
//                    },
//                    "pathSegment": {
//                      "type": "string"
//                    },
//                    "label": {
//                      "type": "string"
//                    },
//                    "icon": {
//                      "type": "string"
//                    }
//                  },
//                  "required": ["entityType", "pathSegment", "label", "icon"]
//                }
//              }
//            },
//            "required": ["nodes"]
//          }
//        },
//        "required": ["data"]
//      }
//    }
//  },
//  "required": ["name", "luigiConfigFragment"]
//}`

type contentConfiguration struct{}

func NewContentConfiguration() ContentConfigurationInterface {
	return &contentConfiguration{}
}

func (cC *contentConfiguration) ValidateJSON(input []byte) (bool, []string) {
	var config ContentConfiguration
	if err := json.Unmarshal(input, &config); err != nil {
		return false, []string{err.Error()}
	}
	return validateSchema(config)
	//return validateSchema(input)
}

func (cC *contentConfiguration) ValidateYAML(input []byte) (bool, []string) {
	var config ContentConfiguration
	if err := yaml.Unmarshal(input, &config); err != nil {
		return false, []string{err.Error()}
	}
	return validateSchema(config)
	//return validateSchema(input)
}

//func validateConfig(config ContentConfiguration) (bool, []string) {
//	var errors []string
//
//	if config.Name == "" {
//		errors = append(errors, "Name is a mandatory parameter and should be specified")
//	}
//
//	for _, fragment := range config.LuigiConfigFragment {
//		for _, node := range fragment.Data.Nodes {
//			if node.EntityType == "" {
//				errors = append(errors, "EntityType is a mandatory parameter and should be specified")
//			}
//			if node.PathSegment == "" {
//				errors = append(errors, "PathSegment is a mandatory parameter and should be specified")
//			}
//			if node.Label == "" {
//				errors = append(errors, "Label is a mandatory parameter and should be specified")
//			}
//			if node.Icon == "" {
//				errors = append(errors, "Icon is a mandatory parameter and should be specified")
//			}
//		}
//	}
//
//	return len(errors) == 0, errors
//}

//func (cC *contentConfiguration) ValidateSchema(input []byte, schema string) (bool, []string) {
//	schemaLoader := gojsonschema.NewStringLoader(schema)
//	documentLoader := gojsonschema.NewBytesLoader(input)
//
//	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
//	if err != nil {
//		return false, []string{err.Error()}
//	}
//
//	if !result.Valid() {
//		var errors []string
//		for _, desc := range result.Errors() {
//			errors = append(errors, desc.String())
//		}
//		return false, errors
//	}
//
//	return true, nil
//}

//func (cC *contentConfiguration) ValidateSchema(config ContentConfiguration, schema string) (bool, []string) {
//
//	s := jsonschema.Reflect(config)
//	data, err := json.MarshalIndent(s, "", "  ")
//	if err != nil {
//		panic(err.Error())
//	}
//	fmt.Println(string(data))
//
//	return true, nil
//	result, err := jsonschema.Validate(gojsonschema.NewStringLoader(schema), gojsonschema.NewGoLoader(config))
//
//	schemaLoader := gojsonschema.NewStringLoader(schema)
//	documentLoader := gojsonschema.NewBytesLoader(input)
//
//	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
//	if err != nil {
//		return false, []string{err.Error()}
//	}
//
//	if !result.Valid() {
//		var errors []string
//		for _, desc := range result.Errors() {
//			errors = append(errors, desc.String())
//		}
//		return false, errors
//	}
//
//	return true, nil
//}

// func validateSchema(input []byte) (bool, []string) {
func validateSchema(config ContentConfiguration) (bool, []string) {

	// Unmarshal the input to a generic map
	//var document map[string]interface{}
	//if err := json.Unmarshal(input, &document); err != nil {
	//	return false, []string{err.Error()}
	//}

	result := jsonschema.Reflect(config)

	_ = result

	//// Load the schema from the constant
	//var schema jsonschema.Schema
	//if err := json.Unmarshal([]byte(ContentConfigurationSchema), &schema); err != nil {
	//	return false, []string{err.Error()}
	//}

	// Validate the document against the schema

	//errs := jsonschema..Validate(document)

	// Validate the document against the schema
	//errs := schema.Validate(document)
	//if len(errs) > 0 {
	//	var errorMessages []string
	//	for _, e := range errs {
	//		errorMessages = append(errorMessages, e.Description)
	//	}
	//	return false, errorMessages
	//}
	//
	return true, nil
}
