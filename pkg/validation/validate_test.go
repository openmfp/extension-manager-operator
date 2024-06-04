package validation

import (
	"bytes"
	"encoding/json"
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
