package config_validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createJson(t *testing.T) {
	result := createJson(&ContentConfiguration{})
	expected := get_createJsonFixture()

	assert.IsType(t, []byte{}, result)
	assert.Equal(t, expected, result)
}

// Hint: schema generated with `omitempty` tag in struct fields
func get_createJsonFixture() []byte {
	return []byte(`{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration","$defs":{"LuigiConfigData":{"properties":{"nodes":{"items":{"$ref":"#/$defs/Node"},"type":"array"}},"additionalProperties":false,"type":"object","required":["nodes"]},"LuigiConfigFragment":{"properties":{"data":{"$ref":"#/$defs/LuigiConfigData"}},"additionalProperties":false,"type":"object"},"Node":{"properties":{"entityType":{"type":"string"},"pathSegment":{"type":"string"},"label":{"type":"string"},"icon":{"type":"string"}},"additionalProperties":false,"type":"object"}},"properties":{"name":{"type":"string"},"luigiConfigFragment":{"items":{"$ref":"#/$defs/LuigiConfigFragment"},"type":"array"}},"additionalProperties":false,"type":"object"}`)

}
