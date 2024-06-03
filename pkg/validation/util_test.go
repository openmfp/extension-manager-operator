package validation

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadSchemaJSONFromFile(t *testing.T) {
	// Create a temporary file for testing (in /tmp FS)
	tmpFile, err := ioutil.TempFile("", "mock_schema.json.out")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	schemaJSONContent := `{
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://github.com/openmfp/extension-content-operator/pkg/validation/content-configuration",
        "$defs": {
            "LuigiConfigData": {
                "properties": {
                    "nodes": {
                        "items": {
                            "$ref": "#/$defs/Node"
                        },
                        "type": "array"
                    }
                },
                "additionalProperties": false,
                "type": "object",
                "required": ["nodes"]
            },
            "LuigiConfigFragment": {
                "properties": {
                    "data": {
                        "$ref": "#/$defs/LuigiConfigData"
                    }
                },
                "additionalProperties": false,
                "type": "object",
                "required": ["data"]
            },
            "Node": {
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
                "additionalProperties": false,
                "type": "object",
                "required": ["entityType", "pathSegment", "label", "icon"]
            }
        },
        "properties": {
            "name": {
                "type": "string"
            },
            "luigiConfigFragment": {
                "items": {
                    "$ref": "#/$defs/LuigiConfigFragment"
                },
                "type": "array"
            }
        },
        "additionalProperties": false,
        "type": "object",
        "required": ["name", "luigiConfigFragment"]
    }`
	if _, err := tmpFile.Write([]byte(schemaJSONContent)); err != nil {
		t.Fatal(err)
	}

	// Close the file
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test loading the schema JSON from the temporary file
	schemaJSON, err := loadSchemaJSONFromFile(tmpFile.Name())

	assert.NoError(t, err)
	assert.NotNil(t, schemaJSON)
	assert.Equal(t, []byte(schemaJSONContent), schemaJSON)
}
