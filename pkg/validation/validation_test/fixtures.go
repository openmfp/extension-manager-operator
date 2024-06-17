package validation_test

import (
	"bytes"
	"encoding/json"
	"log"

	"gopkg.in/yaml.v3"
)

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

func GetInvalidTypeYAML() string {
	return `
name: overview
luigiConfigFragment:
 - data:
     nodes: "string"
`
}

func GetValidJSONButDifferentName() string {
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
		"name": "overview2"
	}`
}

func GetValidYAMLFixtureButDifferentName() string {
	return `
iAmOptionalCustomFieldThatShouldBeStored: iAmOptionalCustomValue
name: overview2
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
