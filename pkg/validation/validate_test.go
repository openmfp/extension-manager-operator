package validation

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

	type testCase struct {
		name          string
		input         []byte
		contentType   string
		expectedError string
	}

	testCases := []testCase{
		{
			name:          "valid JSON",
			input:         []byte(validJSON),
			contentType:   "json",
			expectedError: "",
		},
		{
			name:          "invalid JSON",
			input:         []byte(invalidJSON),
			contentType:   "json",
			expectedError: "some expected error message for invalid JSON", // replace with actual error message
		},
		{
			name:          "valid YAML",
			input:         []byte(validYAML),
			contentType:   "yaml",
			expectedError: "",
		},
		{
			name:          "invalid YAML",
			input:         []byte(invalidYAML),
			contentType:   "yaml",
			expectedError: "some expected error message for invalid YAML", // replace with actual error message
		},
		{
			name:          "unsupported content type",
			input:         []byte(validJSON),
			contentType:   "xml",
			expectedError: "no validator found for content type",
		},
		{
			name:          "invalid content type",
			input:         []byte(validJSON),
			contentType:   "invalid",
			expectedError: "no validator found for content type",
		},
		{
			name:          "empty input",
			input:         []byte{},
			contentType:   "json",
			expectedError: "empty input provided",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cC.Validate(tc.input, tc.contentType)
			if tc.expectedError == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil || err.Error() != tc.expectedError {
					t.Errorf("expected error: %v, got: %v", tc.expectedError, err)
				}
			}
		})
	}
}

//func TestValidate(t *testing.T) {
//	cC := NewContentConfiguration()
//
//	validJSON := `{
//		"name": "overview",
//		"luigiConfigFragment": [
//			{
//				"data": {
//					"nodes": [
//						{
//							"entityType": "global",
//							"pathSegment": "home",
//							"label": "Overview",
//							"icon": "home"
//						}
//					]
//				}
//			}
//		]
//	}`
//
//	invalidJSON := `{
//			"name": "overview",
//			"luigiConfigFragment": [
//				{
//					"data": {
//						"nodes": [
//							{
//								"entityType": "global",
//								"pathSegment": "home",
//								"label": "Overview"
//							}
//						]
//					}
//				}
//			]
//		}`
//
//	validYAML := `
//name: overview
//luigiConfigFragment:
//  - data:
//      nodes:
//        - entityType: global
//          pathSegment: home
//          label: Overview
//          icon: home
//`
//
//	invalidYAML := `
//name: overview
//luigiConfigFragment:
//  - data:
//      nodes:
//        - entityType: global
//          pathSegment: home
//`
//
//	// Test valid JSON
//	_, err := cC.Validate([]byte(validJSON), "json")
//	if err != nil {
//		t.Errorf("expected valid JSON to pass validation, got error: %v", err)
//	}
//
//	// Test invalid JSON
//	_, err = cC.Validate([]byte(invalidJSON), "json")
//	if err == nil {
//		t.Error("expected invalid JSON to fail validation, but it passed")
//	}
//
//	// Test valid YAML
//	_, err = cC.Validate([]byte(validYAML), "yaml")
//	if err != nil {
//		t.Errorf("expected valid YAML to pass validation, got error: %v", err)
//	}
//
//	// Test invalid YAML
//	_, err = cC.Validate([]byte(invalidYAML), "yaml")
//	if err == nil {
//		t.Error("expected invalid YAML to fail validation, but it passed")
//	}
//
//	// Test unsupported content type
//	_, err = cC.Validate([]byte(validJSON), "xml")
//	if err == nil || err.Error() != "no validator found for content type" {
//		t.Errorf("expected error for unsupported content type, got: %v", err)
//	}
//
//	// Test invalid content type
//	_, err = cC.Validate([]byte(validJSON), "invalid")
//	if err == nil || err.Error() != "no validator found for content type" {
//		t.Errorf("expected error for invalid content type, got: %v", err)
//	}
//
//	// Test empty input
//	_, err = cC.Validate([]byte{}, "json")
//	if err == nil || err.Error() != "empty input provided" {
//		t.Errorf("expected error for empty input, got: %v", err)
//	}
//}
