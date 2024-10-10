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

func GetValidYaml_targetAppConfig_viewGroup() string {
	return `{
  "name": "extension-manager",
  "contentType": "json",
  "luigiConfigFragment": {
      "data": {
        "targetAppConfig": {
        "_version": "1.13.0",
        "sap.integration": {
          "navMode": "inplace",
          "urlTemplateId": "urltemplate.url",
          "urlTemplateParams": {
            "query": {}
          }
        }
      },
      "viewGroup": {
        "preloadSuffix": "/#/preload"
      },
      "nodes": [
        {
          "entityType": "global",
          "pathSegment": "catalog",
          "label": "{{catalog}}",
          "icon": "business-one",
          "dxpOrder": 6,
          "order": 6,
          "hideSideNav": true,
          "tabNav": true,
          "showBreadcrumbs": false,
          "urlSuffix": "/#/global-catalog",
          "visibleForFeatureToggles": ["!global-catalog"]
        },
        {
          "entityType": "global",
          "pathSegment": "catalog",
          "label": "{{catalog}}",
          "icon": "business-one",
          "dxpOrder": 6,
          "order": 6,
          "hideSideNav": true,
          "tabNav": true,
          "showBreadcrumbs": false,
          "urlSuffix": "/#/new-global-catalog",
          "visibleForFeatureToggles": ["global-catalog"]
        },
        {
          "entityType": "global",
          "pathSegment": "extensions",
          "label": "{{extensions}}",
          "hideFromNav": true,
          "children": [
            {
              "pathSegment": ":extClassName",
              "hideFromNav": true,
              "urlSuffix": "/#/extensions/:extClassName",
              "context": {
                "extClassName": ":extClassName"
              }
            }
          ]
        }
      ],
      "texts": [
        {
          "locale": "",
          "textDictionary": {
            "catalog": "Catalog",
            "extensions": "Extensions"
          }
        },
        {
          "locale": "en",
          "textDictionary": {
            "catalog": "Catalog",
            "extensions": "Extensions"
          }
        },
        {
          "locale": "de",
          "textDictionary": {
            "catalog": "Katalog",
            "extensions": "Erweiterungen"
          }
        }
      ]
    }
  }
}`
}

func GetValidYAML_node_category_string() string {
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
    category: cat1
  texts:
  - locale: de
    textDictionary:
      hello: Hallo
`
}

func GetValidYAML_node_category_object() string {
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
    category:
      label: cat1
      icon: icon1
      collapsible: false
  texts:
  - locale: de
    textDictionary:
      hello: Hallo
`
}

func GetInalidYAML_node_category_object() string {
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
    category:
      label: cat1
      icon: icon1
      collapsible: false
      invalidfield: invalid
  texts:
  - locale: de
    textDictionary:
      hello: Hallo
`
}

func GetValidYaml_targetAppConfig_viewGroup2() string {
	return `{
    "name": "extension-manager",
    "contentType": "json",
    "luigiConfigFragment": {
        "data": {
            "userSettings": {
                "groups": {
                    "user1": {
                        "label": "label",
                        "sublabel": "sublabel",
                        "title": "title",
                        "icon": "icon",
                        "viewUrl": "viewUrl",
                        "settings": {
                            "option1": {
                                "type": "type",
                                "label": "label",
                                "style": "style",
                                "options": [],
                                "isEditable": false
                            }
                        }
                    }
                }
            },
            "nodeDefaults": {
                "entityType": "type",
                "isolateView": false
            },
            "targetAppConfig": {
                "_version": "1.13.0",
                "sap.integration": {
                    "navMode": "inplace",
                    "urlTemplateId": "urltemplate.url",
                    "urlTemplateParams": {
                        "query": {}
                    }
                }
            },
            "viewGroup": {
                "preloadSuffix": "/#/preload"
            },
            "nodes": [
                {
                    "entityType": "global",
                    "pathSegment": "catalog",
                    "label": "{{catalog}}",
                    "icon": "business-one",
                    "dxpOrder": 6,
                    "order": 6,
                    "hideSideNav": true,
                    "tabNav": true,
                    "showBreadcrumbs": false,
                    "urlSuffix": "/#/global-catalog",
                    "visibleForFeatureToggles": [
                        "!global-catalog"
                    ]
                },
                {
                    "entityType": "global",
                    "pathSegment": "catalog",
                    "label": "{{catalog}}",
                    "icon": "business-one",
                    "dxpOrder": 6,
                    "order": 6,
                    "hideSideNav": true,
                    "tabNav": true,
                    "showBreadcrumbs": false,
                    "urlSuffix": "/#/new-global-catalog",
                    "visibleForFeatureToggles": [
                        "global-catalog"
                    ]
                },
                {
                    "entityType": "global",
                    "pathSegment": "extensions",
                    "label": "{{extensions}}",
                    "hideFromNav": true,
                    "children": [
                        {
                            "pathSegment": ":extClassName",
                            "hideFromNav": true,
                            "urlSuffix": "/#/extensions/:extClassName",
                            "context": {
                                "extClassName": ":extClassName"
                            }
                        }
                    ]
                }
            ],
            "texts": [
                {
                    "locale": "",
                    "textDictionary": {
                        "catalog": "Catalog",
                        "extensions": "Extensions"
                    }
                },
                {
                    "locale": "en",
                    "textDictionary": {
                        "catalog": "Catalog",
                        "extensions": "Extensions"
                    }
                },
                {
                    "locale": "de",
                    "textDictionary": {
                        "catalog": "Katalog",
                        "extensions": "Erweiterungen"
                    }
                }
            ]
        }
    }
}
`
}
