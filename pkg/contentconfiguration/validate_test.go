package contentconfiguration

import "testing"

func TestValidateJSON(t *testing.T) {
	cC := NewContentConfiguration()

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

	invalidJSON := `{
        "name": "overview",
        "luigiConfigFragment": [
            {
                "data": {
                    "nodes": [
                        {
                            "entityType": "global",
                            "pathSegment": "home"
                        }
                    ]
                }
            }
        ]
    }`

	valid, errors := cC.ValidateJSON([]byte(validJSON))
	if !valid {
		t.Errorf("expected valid JSON, got errors: %v", errors)
	}

	valid, errors = cC.ValidateJSON([]byte(invalidJSON))
	if valid {
		t.Errorf("expected invalid JSON, got valid")
	}
	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateYAML(t *testing.T) {
	t.Skip("skipping test")
	cC := NewContentConfiguration()

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

	invalidYAML := `
name: overview
luigiConfigFragment:
  - data:
      nodes:
        - entityType: global
          pathSegment: home
`

	valid, errors := cC.ValidateYAML([]byte(validYAML))
	if !valid {
		t.Errorf("expected valid YAML, got errors: %v", errors)
	}

	valid, errors = cC.ValidateYAML([]byte(invalidYAML))
	if valid {
		t.Errorf("expected invalid YAML, got valid")
	}
	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errors), errors)
	}
}

func Test_validateSchema(t *testing.T) {
	schema := `
	{
		"type": "object",
		"properties": {
			"name": { "type": "string" },
			"age": { "type": "integer" }
		},
		"required": ["name", "age"]
	}`

	validJSON := []byte(`{"name": "John Doe", "age": 30}`)

	invalidJSON := []byte(`{"name": "John Doe", "age": "thirty"}`)

	valid, errors := validateSchema(validJSON, schema)
	if !valid {
		t.Errorf("Expected valid JSON document, got errors: %v", errors)
	}

	valid, errors = validateSchema(invalidJSON, schema)
	if valid {
		t.Errorf("Expected invalid JSON document, got no errors")
	} else {
		expectedError := "age: Invalid type. Expected: integer, given: string"
		found := false
		for _, err := range errors {
			if err == expectedError {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected error %q, but did not find it in errors: %v", expectedError, errors)
		}
	}
}
