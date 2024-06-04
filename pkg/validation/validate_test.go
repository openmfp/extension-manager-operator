package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestValidate(t *testing.T) {
	cC := NewContentConfiguration()
	schema := getJSONSchemaFixture()
	validJSON := getValidJSONFixture()
	invalidJSON := getInvalidJSONFixture()
	validYAML := getValidYAMLFixture()
	invalidYAML := getInvalidYAMLFixture()

	// Test valid JSON
	expected := validJSON
	result, err := cC.Validate(schema, []byte(validJSON), "json")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test invalid JSON
	expected = ""
	result, err = cC.Validate(schema, []byte(invalidJSON), "json")
	if err == nil {
		t.Error("expected invalid JSON to fail validation, but it passed")
	}
	assert.Error(t, err)
	assert.Equal(t, expected, result)
	assert.Contains(t, err.Error(), "The document is not valid:")

	// Test valid YAML
	expected = validJSON
	result, err = cC.Validate(schema, []byte(validYAML), "yaml")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test invalid YAML
	expected = ""
	result, err = cC.Validate(schema, []byte(invalidYAML), "yaml")
	assert.Error(t, err)
	assert.Equal(t, expected, result)
	assert.Contains(t, err.Error(), "The document is not valid:")

	// Test unsupported content type
	result, err = cC.Validate(schema, []byte(validJSON), "xml")
	assert.Error(t, err)
	assert.Equal(t, expected, result)
	assert.Contains(t, err.Error(), "no validator found for content type")

	// Test invalid content type
	result, err = cC.Validate(schema, []byte(validJSON), "invalid")
	assert.Error(t, err)
	assert.Equal(t, expected, result)
	assert.Contains(t, err.Error(), "no validator found for content type")

	// Test empty input
	result, err = cC.Validate(schema, []byte{}, "json")
	assert.Error(t, err)
	assert.Equal(t, expected, result)
	assert.Contains(t, err.Error(), "empty input provided")

	// Test error Marshal
	result, err = cC.Validate(schema, []byte(getInvalidTypeYAMLFixture()), "yaml")
	assert.Error(t, err)
	assert.Equal(t, expected, result)
	assert.Contains(t, err.Error(), "yaml: unmarshal errors:\n  line 3: "+
		"cannot unmarshal !!str `string` into []validation.Node")
}

func Test_validateJSON(t *testing.T) {
	schema := getJSONSchemaFixture()
	validJSON := getValidJSONFixture()
	invalidJSON := getInvalidJSONFixture()

	// Test valid JSON
	result, err := validateJSON(schema, []byte(validJSON))
	assert.NoError(t, err)
	assert.Equal(t, validJSON, result)

	// Test invalid JSON
	result, err = validateJSON(schema, []byte(invalidJSON))
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "The document is not valid:")

	// Test empty JSON
	result, err = validateJSON(schema, []byte{})
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")

	// Test invalid schema
	result, err = validateJSON([]byte{}, []byte(validJSON))
	assert.Error(t, err)
	assert.Equal(t, "", result)
	//assert.Contains(t, err.Error(), "invalid character '}' looking for beginning of value")

	// Test invalid schema and JSON
	result, err = validateJSON([]byte{}, []byte{})
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")

	// Test error Marshal
	result, err = validateJSON(schema, []byte(getInvalidTypeYAMLFixture()))
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "invalid character 'l' looking for beginning of value")
}

func Test_validateYAML(t *testing.T) {
	t.Skipf("skipping test")
	schema := getJSONSchemaFixture()
	validYAML := getValidYAMLFixture()
	invalidYAML := getInvalidYAMLFixture()

	// Test valid YAML
	result, err := validateYAML(schema, []byte(validYAML))
	assert.NoError(t, err)
	assert.Equal(t, validYAML, result)

	// Test invalid YAML
	result, err = validateYAML(schema, []byte(invalidYAML))
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Contains(t, err.Error(), "The document is not valid:")
}

func Test_validateSchema(t *testing.T) {
	schema := getJSONSchemaFixture()
	validJSON := getValidJSONFixture()
	//invalidJSON := getInvalidJSONFixture()

	// Example of valid ContentConfiguration
	validConfig := ContentConfiguration{
		Name: "overview",
		LuigiConfigFragment: []LuigiConfigFragment{
			{
				Data: LuigiConfigData{
					Nodes: []Node{
						{
							EntityType:  "global",
							PathSegment: "home",
							Label:       "Overview",
							Icon:        "home", // Ensure this matches schema if required
						},
					},
				},
			},
		},
	}

	// Example of invalid ContentConfiguration
	invalidConfig := ContentConfiguration{
		// Populate with invalid data according to the schema
	}

	// Test valid schema
	result, err := validateSchema(schema, validConfig)
	fmt.Printf("Valid schema test: result=%s, err=%v\n", result, err)
	assert.NoError(t, err)
	assert.Equal(t, validJSON, result)

	// Test invalid schema
	result, err = validateSchema(schema, invalidConfig)
	fmt.Printf("Invalid schema test: result=%s, err=%v\n", result, err)
	assert.Error(t, err)
	//assert.Equal(t, invalidJSON, result)

	//// Test error Marshal
	//result, err = validateJSON(schema, []byte(getInvalidTypeYAMLFixture()))
	//assert.Error(t, err)
	//assert.Equal(t, "", result)
	//assert.Contains(t, err.Error(), "invalid character 'l' looking for beginning of value")

	// Test error marshal
	type invalidConfig2 struct {
		Name    string
		Channel chan struct{}
	}

	iC2 := invalidConfig2{}
	result, err = validateSchema([]byte{}, iC2)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	//assert.Contains(t, err.Error(), "json: error calling Marshal: json: unsupported type: struct")
}

//func Test_validateSchema(t *testing.T) {
//	//t.Skipf("skipping test")
//	schema := getJSONSchemaFixture()
//	validJSON := getValidJSONFixture()
//	invalidJSON := getInvalidJSONFixture()
//
//	//input, err := validateYAML(schema, []byte(validYAML))
//
//	// Test valid schema
//	result, err := validateSchema(schema, ContentConfiguration{})
//	assert.NoError(t, err)
//	assert.Equal(t, validJSON, result)
//
//	// Test invalid schema
//	result, err = validateSchema(schema, ContentConfiguration{})
//	assert.NoError(t, err)
//	assert.Equal(t, invalidJSON, result)
//}

func getJSONSchemaFixture() []byte {
	schemaFilePath := "./example_schema.json"
	schemaJSON, err := loadSchemaJSONFromFile(schemaFilePath)
	if err != nil {
		log.Fatalf("failed to load schema JSON from file: %v", err)
	}

	return schemaJSON
}

func getValidJSONFixture() string {
	validJSON := `{
		"name": "overview",
		"luigiConfigFragment": [
			{
				"data": {
					"nodes": [
						{
							"entityType": "global",
							"pathSegment": "home",
							"label": "Overview",
							"icon": "home"
						}
					]
				}
			}
		]
	}`

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(validJSON)); err != nil {
		return ""
	}

	return buf.String()
}

func getInvalidJSONFixture() string {
	invalidJSON := `{
			"name": "overview",
			"luigiConfigFragment": [
				{
					"data": {
						"nodes": [
							{
								"entityType": "global",
								"pathSegment": "home",
								"label": "Overview"
							}
						]
					}
				}
			]
		}`

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(invalidJSON)); err != nil {
		return ""
	}

	return buf.String()
}

func getValidYAMLFixture() string {
	validYAML := `
name: overview
luigiConfigFragment:
 - data:
     nodes:
       - entityType: global
         pathSegment: home
         label: Overview
         icon: home
`

	var data interface{}
	err := yaml.Unmarshal([]byte(validYAML), &data)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	compactYAML, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalf("failed to marshal YAML: %v", err)
	}

	return string(compactYAML)
}

func getInvalidYAMLFixture() string {
	invalidYAML := `
name: overview
luigiConfigFragment:
 - data:
     nodes:
       - entityType: global
         pathSegment: home
`

	var data interface{}
	err := yaml.Unmarshal([]byte(invalidYAML), &data)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	compactYAML, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalf("failed to marshal YAML: %v", err)
	}

	return string(compactYAML)
}

func getInvalidTypeYAMLFixture() string {
	invalidYAML := `
name: overview
luigiConfigFragment:
 - data:
     nodes: "string"
`

	var data interface{}
	err := yaml.Unmarshal([]byte(invalidYAML), &data)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	compactYAML, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalf("failed to marshal YAML: %v", err)
	}

	return string(compactYAML)
}
