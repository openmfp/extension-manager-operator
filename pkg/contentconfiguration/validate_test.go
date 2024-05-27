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
