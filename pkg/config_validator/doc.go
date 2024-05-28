// Package config_validator provides functionality to validate configuration data.
//
// It includes functions to validate JSON and YAML configuration against a JSON Schema.
//
// Usage:
//
//	config := NewContentConfiguration()
//	jsonInput := []byte(`{"name": "example", "luigiConfigFragment": [{"data": {"nodes": [{"entityType": "example", "pathSegment": "example", "label": "example", "icon": "example"}]}}]}`)
//	result, err := config.Validate(jsonInput, "json")
//	if err != nil {
//	    fmt.Println("Validation error:", err)
//	} else {
//	    fmt.Println("Validation successful. Result:", result)
//	}
package config_validator
