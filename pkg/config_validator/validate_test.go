package config_validator

import "testing"

func TestValidate(t *testing.T) {
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
								"pathSegment": "home",
								"label": "Overview"
							}
						]
					}
				}
			]
		}`

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

	// Test valid JSON
	result, err := cC.Validate([]byte(validJSON), "json")
	if err != nil {
		t.Errorf("expected valid JSON to pass validation, got error: %v", err)
	}

	// Test invalid JSON
	result1, err := cC.Validate([]byte(invalidJSON), "json")
	if err == nil {
		t.Error("expected invalid JSON to fail validation, but it passed")
	}

	// Test valid YAML
	result2, err := cC.Validate([]byte(validYAML), "yaml")
	if err != nil {
		t.Errorf("expected valid YAML to pass validation, got error: %v", err)
	}

	// Test invalid YAML
	result3, err := cC.Validate([]byte(invalidYAML), "yaml")
	if err == nil {
		t.Error("expected invalid YAML to fail validation, but it passed")
	}

	// Test unsupported content type
	result4, err := cC.Validate([]byte(validJSON), "xml")
	if err == nil || err.Error() != "no validator found for content type" {
		t.Errorf("expected error for unsupported content type, got: %v", err)
	}

	// Test invalid content type
	result5, err := cC.Validate([]byte(validJSON), "invalid")
	if err == nil || err.Error() != "no validator found for content type" {
		t.Errorf("expected error for invalid content type, got: %v", err)
	}

	// Test empty input
	result6, err := cC.Validate([]byte{}, "json")
	if err == nil || err.Error() != "empty input provided" {
		t.Errorf("expected error for empty input, got: %v", err)
	}

	_ = result
	_ = result1
	_ = result2
	_ = result3
	_ = result4
	_ = result5
	_ = result6
}
