package validation

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/openmfp/extension-content-operator/pkg/validation/validation_test"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		contentType string
		schema      []byte
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid_JSON",
			input:       validation_test.GetJSONFixture(validation_test.GetValidJSON()),
			contentType: "json",
			expected:    validation_test.GetJSONFixture(validation_test.GetValidJSON()),
			expectError: false,
		},
		{
			name:        "invalid_JSON_empty_input_ERROR",
			input:       validation_test.GetJSONFixture(`{"name": "overview",`),
			contentType: "json",
			expected:    "",
			expectError: true,
			errorMsg:    "empty input provided",
		},
		{
			name:        "valid_YAML",
			input:       validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
			contentType: "yaml",
			expected:    validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
			expectError: false,
		},
		{
			name:        "unmarshalling_YAML_ERROR",
			input:       validation_test.GetYAMLFixture("!2"),
			contentType: "yaml",
			expected:    "",
			expectError: true,
			errorMsg:    "error unmarshalling YAML: yaml:",
		},
		{
			name:        "the_document_is_not_valid_ERROR",
			input:       validation_test.GetYAMLFixture(`2!`),
			contentType: "yaml",
			expected:    "",
			expectError: true,
			errorMsg: "error unmarshalling YAML: yaml: unmarshal errors:\n  line 1: " +
				"cannot unmarshal !!str `2!` into map[string]interface {}",
		},
		{
			name:        "unsupported_content_type_ERROR",
			input:       validation_test.GetJSONFixture(validation_test.GetValidJSON()),
			contentType: "xml",
			expected:    "",
			expectError: true,
			errorMsg:    "no validator found for content type",
		},
		{
			name:        "empty_input_ERROR",
			input:       "",
			contentType: "json",
			expected:    "",
			expectError: true,
			errorMsg:    "empty input provided",
		},
		{
			name:        "validating_with_incorrect_schema_ERROR",
			schema:      []byte("123"),
			contentType: "json",
			input:       validation_test.GetJSONFixture(validation_test.GetValidJSON()),
			expected:    "",
			expectError: true,
			errorMsg:    "error validating JSON data",
		},
		{
			name:        "unmarshal_string_into_Go_struct_ERROR",
			schema:      getJSONSchemaFixture(),
			input:       validation_test.GetYAMLFixture(validation_test.GetInvalidTypeYAML()),
			contentType: "yaml",
			expected:    "",
			expectError: true,
			errorMsg:    "The document is not valid:",
		},
		{
			name:        "valid_JSON_empty_locale",
			input:       validation_test.GetJSONFixture(validation_test.GetValidJSONWithEmptyLocale()),
			contentType: "json",
			expected:    validation_test.GetJSONFixture(validation_test.GetValidJSONWithEmptyLocale()),
			expectError: false,
		},
		{
			name:        "test_luigiConfigFragment",
			input:       validation_test.GetJSONFixture(validation_test.GetluigiConfigFragment()),
			contentType: "json",
			expected:    validation_test.GetJSONFixture(validation_test.GetluigiConfigFragment()),
			expectError: false,
		},
		{
			name:        "test_luigiConfigFragment",
			input:       validation_test.GetYAMLFixture(validation_test.GetValidYaml_targetAppConfig_viewGroup()),
			contentType: "yaml",
			expected:    validation_test.GetYAMLFixture(validation_test.GetValidYaml_targetAppConfig_viewGroup()),
			expectError: false,
			schema:      nil,
		},
		{
			name:        "test_node_category_string",
			input:       validation_test.GetYAMLFixture(validation_test.GetValidYAML_node_category_string()),
			contentType: "yaml",
			expected:    validation_test.GetYAMLFixture(validation_test.GetValidYAML_node_category_string()),
			expectError: false,
			schema:      nil,
		},
		{
			name:        "test_node_category_object",
			input:       validation_test.GetYAMLFixture(validation_test.GetValidYAML_node_category_object()),
			contentType: "yaml",
			expected:    validation_test.GetYAMLFixture(validation_test.GetValidYAML_node_category_object()),
			expectError: false,
			schema:      nil,
		},
		{
			name:        "test_node_category_invalidobject",
			input:       validation_test.GetYAMLFixture(validation_test.GetInalidYAML_node_category_object()),
			contentType: "yaml",
			expected:    "",
			expectError: true,
			schema:      nil,
		},
		{
			name:        "test_luigiConfigFragment",
			input:       validation_test.GetYAMLFixture(validation_test.GetValidYaml_targetAppConfig_viewGroup2()),
			contentType: "yaml",
			expected:    validation_test.GetYAMLFixture(validation_test.GetValidYaml_targetAppConfig_viewGroup2()),
			expectError: false,
			schema:      nil,
		},
	}

	cC := NewContentConfiguration()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.schema != nil {
				cC.WithSchema(tc.schema) // nolint: errcheck
			}
			result, err := cC.Validate([]byte(tc.input), tc.contentType)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expected, result)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func Test_validateSchema(t *testing.T) {
	type ContentConfigurationTypeMock struct {
		Name    interface{} `json:"name"`
		Surname interface{} `json:"surname"`
	}

	schema := getJSONSchemaFixture()

	tests := []struct {
		name           string
		input          interface{}
		expectedErrMsg string
	}{
		{
			name: "Invalid_Type",
			input: ContentConfigurationTypeMock{
				Name: 1, // wrong type
			},
			expectedErrMsg: "The document is not valid:\n[luigiConfigFragment is required (root):" +
				" Additional property surname is not allowed field 'name' is invalid, got '%!s(<nil>)', expected 'string']",
		},
		{
			name: "Invalid_JSON",
			input: ContentConfigurationTypeMock{
				Name:    "John",
				Surname: make(chan int), // invalid type for JSON marshaling
			},
			expectedErrMsg: "error validating JSON data",
		},
		{
			name: "luigiConfigFragment_is_required",
			input: []byte(`{
				"name": "overview"
			}`),
			expectedErrMsg: "The document is not valid:\n[field '(root)' is invalid",
		},
		{
			name: "name_is_required",
			input: ContentConfiguration{
				LuigiConfigFragment: LuigiConfigFragment{
					Data: LuigiConfigData{
						Nodes: []Node{
							{
								EntityType: "global",
							},
						},
					},
				},
			},
			expectedErrMsg: "The document is not valid:\n[(root): Must validate one and only one schema (oneOf) name is required]",
		},
		{
			name: "nodes_is_required",
			input: ContentConfiguration{
				Name:                "overview",
				LuigiConfigFragment: LuigiConfigFragment{},
			},
			expectedErrMsg: "The document is not valid:\n[luigiConfigFragment.data: " +
				"Must validate one and only one schema (oneOf) nodes is required]",
		},
		{
			name: "textDictionary_is_required",
			input: ContentConfiguration{
				Name: "overview",
				LuigiConfigFragment: LuigiConfigFragment{
					Data: LuigiConfigData{
						Nodes: []Node{
							{
								EntityType: "global",
							},
						},
						Texts: []Text{{
							Locale: "de",
						}},
					},
				},
			},
			expectedErrMsg: "The document is not valid:\n[field 'luigiConfigFragment.data.texts.0.textDictionary'" +
				" is invalid, got '%!s(<nil>)', expected 'object']",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byteArray, _ := json.Marshal(tt.input)

			err := validateSchemaBytes(schema, byteArray)
			assert.Error(t, err)
			actualStr := err.Error()
			expectedStr := tt.expectedErrMsg
			assert.Contains(t, actualStr, expectedStr)
			assert.Contains(t, err.Error(), tt.expectedErrMsg)
		})
	}
}

func getJSONSchemaFixture() []byte {
	schemaFilePath := "./schema/schema_autogen.json"
	schemaJSON, err := loadSchemaJSONFromFile(schemaFilePath)
	if err != nil {
		log.Fatalf("failed to load schema JSON from file: %v", err)
	}

	return schemaJSON
}

func TestWithSchema(t *testing.T) {
	cC := NewContentConfiguration()
	empty := ""
	err := cC.WithSchema([]byte(empty))
	assert.Error(t, err)
}
