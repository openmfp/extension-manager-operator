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
		"luigiConfigFragment": {
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
		},
		"name": "overview"
	}`
}

func GetValidJSONWithEmptyLocale() string {
	return `{
		"luigiConfigFragment": {
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
						"locale": "",
						"textDictionary": {
							"hello": "Hello"
						}
					},
					{
						"locale": "de",
						"textDictionary": {
							"hello": "Hallo"
						}
					}
				]
			}
		},
		"name": "overview"
	}`
}

func GetValidYAML() string {
	return `
name: overview
luigiConfigFragment:
 data:
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

func GetValidIncompatibleYAML() string {
	return `
iAmOptionalCustomFieldThatShouldBeStored: iAmOptionalCustomValue
name: overview
luigiConfigFragment:
 data:
  nodeDefaults:
    entityType: global
    isolateView: true
  nodes:
  - entityType: global
    pathSegment: home
    label: Overview
    icon: home
  texts:
  - textDictionary:
      hello: Hallo
`
}

func GetInvalidTypeYAML() string {
	return `
name: overview
luigiConfigFragment:
  data:
    nodes: "string"
`
}

func GetValidJSONButDifferentName() string {
	return `{
		"luigiConfigFragment": {
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
		},
		"name": "overview2"
	}`
}

func GetValidYAMLFixtureButDifferentName() string {
	return `
name: overview2
luigiConfigFragment:
 data:
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

func GetluigiConfigFragment() string {
	return ` {
        "name": "accounts",
        "luigiConfigFragment": {
            "data": {
              "nodes": [
                {
                  "pathSegment": "create",
                  "hideFromNav": true,
                  "entityType": "main",
                  "loadingIndicator": {
                    "enabled": false
                  },
                  "keepSelectedForChildren": true,
                  "url": "https://some.url/modal/create",
                  "children": []
                },
                {
                  "pathSegment": "accounts",
                  "label": "Accounts",
                  "entityType": "main",
                  "loadingIndicator": {
                    "enabled": false
                  },
                  "keepSelectedForChildren": true,
                  "url": "https://some.url/accounts",
                  "children": [
                    {
                      "pathSegment": ":accountId",
                      "hideFromNav": true,
                      "keepSelectedForChildren": false,
                      "defineEntity": {
                        "id": "account"
                      },
                      "context": {
                        "accountId": ":accountId"
                      }
                    }
                  ]
                },
                {
                  "pathSegment": "overview",
                  "label": "Overview",
                  "entityType": "main.account",
                  "loadingIndicator": {
                    "enabled": false
                  },
                  "visibleForFeatureToggles": ["oldAccount"],
                  "url": "https://some.url/accounts/:accountId"
                }
              ]
            }
          }
      }`
}
