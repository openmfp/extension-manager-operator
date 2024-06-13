package validation

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func loadSchemaJSONFromFile(filePath string) ([]byte, error) {
	schemaJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return schemaJSON, nil
}

func GetJSONFixture(input string) string {
	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(input)); err != nil {
		return ""
	}

	return buf.String()
}

func GetYAMLFixture(input string) string {
	var data interface{}
	err := yaml.Unmarshal([]byte(input), &data)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	compactYAML, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalf("failed to marshal YAML: %v", err)
	}

	return string(compactYAML)
}

func GetValidJSON() string {
	return `{
		"iAmOptionalCustomFieldThatShouldBeStored": "iAmOptionalCustomValue",
		"luigiConfigFragment": [
			{
				"data": {
					"nodeDefaults": {
						"entityType": "global",
						"isolateView": true
					},
					"nodes": [
						{
							"entityType": "global",
							"icon": "home",
							"label": "Overview",
							"pathSegment": "home"
						}
					],
					"texts": [
						{
							"locale": "de",
							"textDictionary": {
								"hello": "Hallo"
							}
						}
					]
				}
			}
		],
		"name": "overview"
	}`
}

func GetInvalidJSON() string {
	return `{
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
}

func GetValidYAML() string {
	return `
iAmOptionalCustomFieldThatShouldBeStored: iAmOptionalCustomValue
name: overview
luigiConfigFragment:
- data:
   nodeDefaults:
     entityType: global
     isolateView: true
   nodes:
   - entityType: global
     pathSegment: home
     label: Overview
     icon: home
   texts:
   - locale: de
     textDictionary:
       hello: Hallo
`
}

func getInvalidYAML() string {
	return `
name: overview
luigiConfigFragment:
 - data:
     nodes:
       - entityType: global
         pathSegment: home
`
}

func GetInvalidTypeYAML() string {
	return `
name: overview
luigiConfigFragment:
 - data:
     nodes: "string"
`
}
