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

func GetValidJSON_extension_manager_ui1() string {
	return `      {
        "name": "extension-manager",
        "luigiConfigFragment": {
          "data": {
            "viewGroup": {
              "preloadSuffix": "/#/preload"
            },
            "nodes": [
              {
                "pathSegment": "catalog",
                "label": "{{extensions}}",
                "icon": "cart",
                "entityType": "project",
                "navSlot": "settings",
                "dxpOrder": 10,
                "order": 10,
                "urlSuffix": "/#/catalog",
                "testId": "dxp-frame-navigation-project-extensions-catalog",
                "defineEntity": {
                  "id": "account",
                  "useBack": true
                },
                "keepSelectedForChildren": true
              },
              {
                "pathSegment": "catalog",
                "label": "{{extensions}}",
                "icon": "cart",
                "entityType": "team",
                "navSlot": "settings",
                "dxpOrder": 10,
                "order": 10,
                "urlSuffix": "/#/catalog",
                "testId": "dxp-frame-navigation-team-extensions-catalog",
                "defineEntity": {
                  "id": "account",
                  "useBack": true
                },
                "keepSelectedForChildren": true
              },
              {
                "entityType": "project.account",
                "pathSegment": "create-res/:scope/:extClassName/account/:accountType",
                "hideFromNav": true,
                "urlSuffix": "/#/extensions/:scope/:extClassName/account/:accountType/create-resource",
                "context": {
                  "extClassName": ":extClassName"
                }
              },
              {
                "entityType": "project.account",
                "pathSegment": "edit-res/:scope/:extClassName/account/:accountType/:name/:nspace",
                "hideFromNav": true,
                "urlSuffix": "/#/extensions/:scope/:extClassName/account/:accountType/edit-resource/:name/:nspace",
                "context": {
                  "extClassName": ":extClassName"
                }
              },
              {
                "entityType": "project",
                "pathSegment": "accounts",
                "hideFromNav": true,
                "urlSuffix": "/#/catalog",
                "context": {
                  "layout": "TwoColumnsMidExpanded",
                  "extClassName": "dxp-github-ui"
                }
              },
              {
                "entityType": "project",
                "pathSegment": "install-extensions",
                "hideFromNav": true,
                "urlSuffix": "/#/install-extensions"
              },
              {
                "entityType": "team",
                "pathSegment": "install-extensions",
                "hideFromNav": true,
                "urlSuffix": "/#/install-extensions"
              },
              {
                "pathSegment": "extension-missing-mandatory-data",
                "hideFromNav": true,
                "context": {
                  "providesMissingMandatoryDataUrl": true
                },
                "urlSuffix": "/#/extension-missing-mandatory-data/:extClassName",
                "entityType": "project",
                "testId": "dxp-frame-navigation-project-extension-missing-mandatory-data"
              },
              {
                "entityType": "project",
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
              },
              {
                "entityType": "team",
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
              },
              {
                "category": {
                  "id": "community-extensions",
                  "label": "{{communityExtensions}}",
                  "isGroup": true
                },
                "dxpOrder": 10,
                "order": 20,
                "entityType": "project"
              },
              {
                "category": {
                  "id": "community-extensions",
                  "label": "{{communityExtensions}}",
                  "isGroup": true
                },
                "dxpOrder": 20,
                "order": 20,
                "entityType": "project.component"
              }
            ],
            "texts": [
              {
                "locale": "",
                "textDictionary": {
                  "extensions": "Extensions",
                  "all": "All",
                  "communityExtensions": "Community Extensions"
                }
              },
              {
                "locale": "en",
                "textDictionary": {
                  "extensions": "Extensions",
                  "all": "All",
                  "communityExtensions": "Community Extensions"
                }
              },
              {
                "locale": "de",
                "textDictionary": {
                  "extensions": "Erweiterungen",
                  "all": "Alle",
                  "communityExtensions": "Community Erweiterungen"
                }
              }
            ]
          }
        }
      }
`
}

func GetValidJSON_github_ui1() string {
	return `      {
        "name": "github-ui",
        "luigiConfigFragment": {
            "data": {
                "viewGroup": {
                  "preloadSuffix": "/#/preload",
                  "requiredIFramePermissions": {
                    "allow": ["clipboard-read", "clipboard-write"]
                  }
                },
                "nodes": [
                  {
                    "pathSegment": "github-loading-screen",
                    "hideFromNav": true,
                    "urlSuffix": "/#/projects/:projectId/github-loading-screen",
                    "entityType": "project"
                  },
                  {
                    "pathSegment": "github",
                    "hideFromNav": true,
                    "urlSuffix": "/#/projects/:projectId/connect-account-dialog",
                    "entityType": "project.account"
                  },
                  {
                    "pathSegment": "github-loading-screen",
                    "hideFromNav": true,
                    "urlSuffix": "/#/teams/:teamId/github-loading-screen",
                    "entityType": "team"
                  },
                  {
                    "pathSegment": "github",
                    "hideFromNav": true,
                    "urlSuffix": "/#/teams/:teamId/connect-account-dialog",
                    "entityType": "team.account"
                  },
                  {
                    "pathSegment" : "github-code",
                    "label" : "Code",
                    "url": "{context.entityContext.component.annotations[\"github.dxp.sap.com/repo-url\"]}",
                    "virtualTree": false,
                    "isolateView": true,
                    "loadingIndicator": {
                      "enabled": false
                    },
                    "entityType": "project.component",
                    "icon": "source-code",
                    "visibleForPlugin": true,
                    "visibleForContext": "serviceProviderConfig.disableGithubCode == null  || serviceProviderConfig.disableGithubCode == 'false'",
                    "category": {
                      "label": "{{development}}",
                      "collapsable": false,
                      "dxpOrder": 100,
                      "order": 100
                    }
                  },
                  {
                    "pathSegment" : "github-pulls",
                    "label" : "Pulls",
                    "url": "{context.entityContext.component.annotations[\"github.dxp.sap.com/repo-url\"]}/pulls",
                    "virtualTree": false,
                    "isolateView": true,
                    "visibleForPlugin": true,
                    "visibleForContext": "serviceProviderConfig.disableGithubPullRequests == null  || serviceProviderConfig.disableGithubPullRequests == 'false'",
                    "loadingIndicator": {
                      "enabled": false
                    },
                    "entityType": "project.component",
                    "icon": "wrench",
                    "category": { "label": "{{development}}" }
                  },
                  {
                    "pathSegment" : "github-issues",
                    "label" : "Issues",
                    "url": "{context.entityContext.component.annotations[\"github.dxp.sap.com/repo-url\"]}/issues",
                    "virtualTree": false,
                    "isolateView": true,
                    "visibleForPlugin": true,
                    "visibleForContext": "serviceProviderConfig.disableGithubIssues == null  || serviceProviderConfig.disableGithubIssues == 'false'",
                    "loadingIndicator": {
                      "enabled": false
                    },
                    "entityType": "project.component",
                    "icon": "task",
                    "category": {
                      "label": "{{issueManagement}}",
                      "collapsable": false,
                      "dxpOrder": 200,
                      "order": 200
                    }
                  }
                ],
                "texts": [
                  {
                    "locale": "",
                    "textDictionary": {
                      "quality": "Security & Quality",
                      "issueManagement": "Issue Management",
                      "development": "Development"
                    }
                  },
                  {
                    "locale": "en",
                    "textDictionary": {
                      "quality": "Security & Quality",
                      "issueManagement": "Issue Management",
                      "development": "Development"
                    }
                  },
                  {
                    "locale": "de",
                    "textDictionary": {
                      "quality": "Sicherheit & Qualität",
                      "issueManagement": "Issue Management",
                      "development": "Development"
                    }
                  }
                ]
            }
        }
      }
`
}

func GetValidJSON_github_wc() string {
	return `      {
        "name": "github-wc",
        "luigiConfigFragment": {
            "data": {
                "viewGroup": {
                    "preloadSuffix": "/#/preload",
                    "requiredIFramePermissions": {
                      "allow": ["clipboard-read", "clipboard-write"]
                  }
                },
                "nodes": [
                  {
                    "entityType": "project.overview::compound",
                    "pathSegment": "add-github-account-card",
                    "urlSuffix": "/main.js#add-github-account-card",
                    "visibleForContext": "(serviceProviderConfig.skipOnboardingCard == null  || serviceProviderConfig.skipOnboardingCard == \"false\") && ( serviceProviderConfig.githubAccountAdded == null  || serviceProviderConfig.githubAccountAdded == \"false\")",
                    "visibleForEntityContext": {
                      "project": {
                        "policies": ["iamMember"]
                      }
                    },
                    "layoutConfig": {
                      "slot": "recommended-actions",
                      "order": 10
                    },
                    "webcomponent": {
                      "selfRegistered": true
                    }
                  }
                ],
                "texts": [
                  {
                    "locale": "",
                    "textDictionary": {
                      "quality": "Security & Quality",
                      "issueManagement": "Issue Management",
                      "development": "Development"
                    }
                  },
                  {
                    "locale": "en",
                    "textDictionary": {
                      "quality": "Security & Quality",
                      "issueManagement": "Issue Management",
                      "development": "Development"
                    }
                  },
                  {
                    "locale": "de",
                    "textDictionary": {
                      "quality": "Sicherheit & Qualität",
                      "issueManagement": "Issue Management",
                      "development": "Development"
                    }
                  }
                ]
            }
        }
      }
`
}

func GetValidJSON_iam_ui() string {
	return `{
        "name": "iam-ui",
        "luigiConfigFragment": {
          "data": {
            "viewGroup": {
              "preloadSuffix": "/#/preload",
              "requiredIFramePermissions": {
                "allow": ["clipboard-read", "clipboard-write"]
              }
            },
            "nodes": [
              {
                "entityType": "project",
                "pathSegment": "members",
                "label": "{{members}}",
                "icon": "company-view",
                "hideFromNav": false,
                "urlSuffix": "/#/projects/:projectId/members",
                "navSlot": "settings",
                "dxpOrder": 30,
                "order": 30
              },
              {
                "entityType": "project",
                "pathSegment": "add-members",
                "label": "{{members}}",
                "hideFromNav": true,
                "urlSuffix": "/#/projects/:projectId/add-members"
              },
              {
                "entityType": "team",
                "pathSegment": "members",
                "label": "{{members}}",
                "icon": "company-view",
                "hideFromNav": false,
                "urlSuffix": "/#/teams/:teamId/members",
                "navSlot": "settings",
                "dxpOrder": 30,
                "order": 30
              },
              {
                "entityType": "team",
                "pathSegment": "add-members",
                "label": "{{members}}",
                "hideFromNav": true,
                "urlSuffix": "/#/teams/:teamId/add-members"
              }
            ],
            "texts": [
              {
                "locale": "",
                "textDictionary": {
                  "members": "Members"
                }
              },
              {
                "locale": "en",
                "textDictionary": {
                  "members": "Members"
                }
              },
              {
                "locale": "de",
                "textDictionary": {
                  "members": "Mitglieder"
                }
              }
            ]
          }
        }
      }
`
}
