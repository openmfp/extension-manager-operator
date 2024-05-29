package config_validator

import (
	"fmt"
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

	result := createJson(&ContentConfiguration{})
	expected := getCreateJsonNonStrictFixture()

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

	result := createJson(&ContentConfiguration{})
	expected := getCreateJsonStrictFixture()

	assert.IsType(t, []byte{}, result)
	assert.Equal(t, expected, result)
}

func getCreateJsonNonStrictFixture() []byte {
	return []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration","$defs":{"LuigiConfigData":{"properties":{"nodes":{"items":{"$ref":"#/$defs/Node"},"type":"array"}},"additionalProperties":false,"type":"object"},"LuigiConfigFragment":{"properties":{"data":{"$ref":"#/$defs/LuigiConfigData"}},"additionalProperties":false,"type":"object"},"Node":{"properties":{"entityType":{"type":"string"},"pathSegment":{"type":"string"},"label":{"type":"string"},"icon":{"type":"string"}},"additionalProperties":false,"type":"object"}},"properties":{"name":{"type":"string"},"luigiConfigFragmenty":{"items":{"$ref":"#/$defs/LuigiConfigFragment"},"type":"array"}},"additionalProperties":false,"type":"object"}`)
}

func getCreateJsonStrictFixture() []byte {
	return []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration","$defs":{"LuigiConfigData":{"properties":{"nodes":{"items":{"$ref":"#/$defs/Node"},"type":"array"}},"additionalProperties":false,"type":"object","required":["nodes"]},"LuigiConfigFragment":{"properties":{"data":{"$ref":"#/$defs/LuigiConfigData"}},"additionalProperties":false,"type":"object","required":["data"]},"Node":{"properties":{"entityType":{"type":"string"},"pathSegment":{"type":"string"},"label":{"type":"string"},"icon":{"type":"string"}},"additionalProperties":false,"type":"object","required":["entityType","pathSegment","label","icon"]}},"properties":{"name":{"type":"string"},"luigiConfigFragmenty":{"items":{"$ref":"#/$defs/LuigiConfigFragment"},"type":"array"}},"additionalProperties":false,"type":"object","required":["name","luigiConfigFragmenty"]}`)
}
