package validation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_nonstrict_createJson(t *testing.T) {
	type Node struct {
		EntityType  string `json:"entityType,omitempty" yaml:"entityType,omitempty"`
		PathSegment string `json:"pathSegment,omitempty" yaml:"pathSegment,omitempty"`
		Label       string `json:"label,omitempty" yaml:"label,omitempty"`
		Icon        string `json:"icon,omitempty" yaml:"icon,omitempty"`
	}

	type LuigiConfigData struct {
		Nodes []Node `json:"nodes,omitempty" yaml:"nodes,omitempty"`
	}

	type LuigiConfigFragment struct {
		Data LuigiConfigData `json:"data,omitempty" yaml:"data,omitempty"`
	}

	type ContentConfiguration struct {
		Name                string                `json:"name,omitempty" yaml:"name,omitempty"`
		LuigiConfigFragment []LuigiConfigFragment `json:"luigiConfigFragmenty,omitempty" yaml:"luigiConfigFragment,omitempty"`
	}

	result := createJSON(&ContentConfiguration{})
	expected := getCreateJSONNonStrictFixture()

	fmt.Println(string(result))

	assert.IsType(t, []byte{}, result)
	assert.Equal(t, expected, result)
}

func Test_strict_createJson(t *testing.T) {
	type Node struct {
		EntityType  string `json:"entityType" yaml:"entityType"`
		PathSegment string `json:"pathSegment" yaml:"pathSegment"`
		Label       string `json:"label" yaml:"label"`
		Icon        string `json:"icon" yaml:"icon"`
	}

	type LuigiConfigData struct {
		Nodes []Node `json:"nodes" yaml:"nodes"`
	}

	type LuigiConfigFragment struct {
		Data LuigiConfigData `json:"data" yaml:"data"`
	}

	type ContentConfiguration struct {
		Name                string                `json:"name" yaml:"name"`
		LuigiConfigFragment []LuigiConfigFragment `json:"luigiConfigFragmenty" yaml:"luigiConfigFragment"`
	}

	result := createJSON(&ContentConfiguration{})
	expected := getCreateJSONStrictFixture()

	assert.IsType(t, []byte{}, result)
	assert.Equal(t, expected, result)
}

func getCreateJSONNonStrictFixture() []byte {
	originalJSON := []byte(`{
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
                "type": "object"
            },
            "LuigiConfigFragment": {
                "properties": {
                    "data": {
                        "$ref": "#/$defs/LuigiConfigData"
                    }
                },
                "additionalProperties": false,
                "type": "object"
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
                "type": "object"
            }
        },
        "properties": {
            "name": {
                "type": "string"
            },
            "luigiConfigFragmenty": {
                "items": {
                    "$ref": "#/$defs/LuigiConfigFragment"
                },
                "type": "array"
            }
        },
        "additionalProperties": false,
        "type": "object"
    }`)

	var compactJSON bytes.Buffer
	if err := json.Compact(&compactJSON, originalJSON); err != nil {
		return nil
	}

	return compactJSON.Bytes()
}

func getCreateJSONStrictFixture() []byte {
	originalJSON := []byte(`{
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
                "required": [
                    "entityType",
                    "pathSegment",
                    "label",
                    "icon"
                ]
            }
        },
        "properties": {
            "name": {
                "type": "string"
            },
            "luigiConfigFragmenty": {
                "items": {
                    "$ref": "#/$defs/LuigiConfigFragment"
                },
                "type": "array"
            }
        },
        "additionalProperties": false,
        "type": "object",
        "required": [
            "name",
            "luigiConfigFragmenty"
        ]
    }`)

	var compactJSON bytes.Buffer
	if err := json.Compact(&compactJSON, originalJSON); err != nil {
		return nil
	}

	return compactJSON.Bytes()
}

func Test_loadSchemaJSONFromFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "mock_schema.json.out")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	// Write JSON content to the temporary file
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

	// Assert no error occurred
	assert.NoError(t, err)

	// Assert schemaJSON is not nil and contains the expected content
	assert.NotNil(t, schemaJSON)
	assert.Equal(t, []byte(schemaJSONContent), schemaJSON)
}
