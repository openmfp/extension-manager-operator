/*
Package configvalidator provides functionality to validate configuration data against a JSON schema.
This package includes a schema definition and functions to validate JSON and YAML configuration inputs.

The primary component of this package is the ContentConfiguration, which can be used to validate
configuration data. The schema is defined as a JSON schema and is used to ensure that the configuration
data conforms to the expected structure.

Example usage:

	package main

	import (
	    "fmt"
	    "log"
	    "your-module-path/configvalidator"
	)

	func main() {
	    configValidator := configvalidator.NewContentConfiguration()

	    // Example JSON input
	    jsonInput := []byte(`{
	        "name": "ExampleConfig",
	        "luigiConfigFragment": [{
	            "data": {
	                "nodes": [{
	                    "entityType": "exampleType",
	                    "pathSegment": "examplePath",
	                    "label": "ExampleLabel",
	                    "icon": "exampleIcon"
	                }]
	            }
	        }]
	    }`)

	    // Validate JSON input
	    result, err := configValidator.Validate(jsonInput, "json")
	    if err != nil {
	        log.Fatalf("JSON validation failed: %v", err)
	    }
	    fmt.Println("JSON validation succeeded:", result)

	    // Example YAML input
	    yamlInput := []byte(`
	    name: ExampleConfig
	    luigiConfigFragment:
	      - data:
	          nodes:
	            - entityType: exampleType
	              pathSegment: examplePath
	              label: ExampleLabel
	              icon: exampleIcon
	    `)

	    // Validate YAML input
	    result, err = configValidator.Validate(yamlInput, "yaml")
	    if err != nil {
	        log.Fatalf("YAML validation failed: %v", err)
	    }
	    fmt.Println("YAML validation succeeded:", result)
	}

Functions:

- NewContentConfiguration: Creates a new instance of ContentConfigurationInterface.
- Validate: Validates the given input data against the content schema based on the content type (JSON or YAML).

The schema is defined as follows:

	var ContentConfigurationSchema = []byte(`{
	    "$schema": "https://json-schema.org/draft/2020-12/schema",
	    "$id": "https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration",
	    "$defs": {
	        "LuigiConfigData": {
	            "properties": {
	                "nodes": {
	                    "items": {"$ref": "#/$defs/Node"},
	                    "type": "array"
	                }
	            },
	            "additionalProperties": false,
	            "type": "object",
	            "required": ["nodes"]
	        },
	        "LuigiConfigFragment": {
	            "properties": {
	                "data": {"$ref": "#/$defs/LuigiConfigData"}
	            },
	            "additionalProperties": false,
	            "type": "object",
	            "required": ["data"]
	        },
	        "Node": {
	            "properties": {
	                "entityType": {"type": "string"},
	                "pathSegment": {"type": "string"},
	                "label": {"type": "string"},
	                "icon": {"type": "string"}
	            },
	            "additionalProperties": false,
	            "type": "object",
	            "required": ["entityType", "pathSegment", "label", "icon"]
	        }
	    },
	    "properties": {
	        "name": {"type": "string"},
	        "luigiConfigFragment": {
	            "items": {"$ref": "#/$defs/LuigiConfigFragment"},
	            "type": "array"
	        }
	    },
	    "additionalProperties": false,
	    "type": "object",
	    "required": ["name", "luigiConfigFragment"]
	}`)

For detailed usage and examples, refer to the main package documentation.
*/
package config_validator
