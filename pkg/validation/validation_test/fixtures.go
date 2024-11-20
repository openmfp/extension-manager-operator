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

func GetValidJSON_learnings() string {
	return `{
        "name": "learning",
        "luigiConfigFragment": {
              "data": {
                  "nodes": [
                    {
                      "pathSegment": "help-portal-documentation",
                      "label": "Documentation",
                      "icon": "document-text",
                      "entityType": "global",
                      "dxpOrder": 6,
                      "order": 6,
                      "url": "https://uacptraining.int.hana.ondemand.com/docs/HYPERSPACE",
                      "visibleForFeatureToggles": ["helpPortalDocumentation"],
                      "visibleForPlugin": true,
                      "networkVisibility": "internal",
                      "hideSideNav": true,
                      "virtualTree": false,
                      "isolateView": true
                    },
                    {
                      "pathSegment": "learning",
                      "label": "{{learning}}",
                      "icon": "education",
                      "entityType": "global",
                      "hideSideNav": true,
                      "dxpOrder": 5,
                      "order": 5,
                      "tabNav": true,
                      "showBreadcrumbs": false,
                      "children": [
                        {
                          "pathSegment": "home",
                          "label": "{{home}}",
                          "icon": "home",
                          "url": "about:blank",
                          "compound": {
                            "renderer": {
                              "use": "grid",
                              "config": {
                                "columns": "minmax(0,1fr) minmax(0,1fr) minmax(0,1fr)",
                                "rows": "[first] repeat(20, auto ) [last]",
                                "layouts": [
                                  {
                                    "minWidth": 0,
                                    "maxWidth": 600,
                                    "columns": "minmax(0,1fr)",
                                    "rows": "[first] repeat(20, auto ) [last]",
                                    "gap": "0px"
                                  },
                                  {
                                    "minWidth": 600,
                                    "maxWidth": 1024,
                                    "columns": "minmax(0,1fr) minmax(0,1fr)",
                                    "rows": "[first] repeat(20, auto ) [last]",
                                    "gap": "0px"
                                  }
                                ]
                              }
                            },
                            "children": [
                              {
                                "urlSuffix": "/microfrontends/feature.js",
                                "context": {
                                  "title": "Hyperspace Portal Documentation",
                                  "content": "Learn how to use the various functionalities of the Hyperspace Portal.",
                                  "alt_label": "{{hero_alt}}",
                                  "gradientColor1": "#02172D",
                                  "gradientColor2": "#203046",
                                  "build_button_label": "Read the documentation",
                                  "feature_link": "/learning/documentation"
                                },
                                "layoutConfig": {
                                  "slot": "content",
                                  "row": "first",
                                  "column": "1"
                                }
                              },
                              {
                                "urlSuffix": "/microfrontends/feature.js",
                                "context": {
                                  "title": "Hyperspace Academy",
                                  "content": "Learn how to use Hyperspace Paved Roads, tools, services, and apply community best practices.",
                                  "alt_label": "{{hero_alt}}",
                                  "gradientColor1": "#DB1F77",
                                  "gradientColor2": "#29313A",
                                  "build_button_label": "Visit the Hyperspace Academy",
                                  "feature_link": "https://pages.github.tools.sap/hyperspace/academy/"
                                },
                                "layoutConfig": {
                                  "slot": "content",
                                  "row": "first",
                                  "column": "2"
                                }
                              },
                              {
                                "urlSuffix": "/microfrontends/feature.js",
                                "context": {
                                  "title": "Hyperspace SharePoint",
                                  "content": "Get an overview of what Hyperspace is all about and stay informed about recent updates.",
                                  "alt_label": "{{hero_alt}}",
                                  "gradientColor1": "#57CC99",
                                  "gradientColor2": "#29313A",
                                  "build_button_label": "Visit the Hyperspace SharePoint",
                                  "feature_link": "https://sap.sharepoint.com/sites/124706"
                                },
                                "layoutConfig": {
                                  "slot": "content",
                                  "row": "first",
                                  "column": "3"
                                }
                              },
                              {
                                "urlSuffix": "/microfrontends/quicklinks.js",
                                "context": {
                                  "title": "Get Going",
                                  "description": "Helpful links in the context of the Hyperspace Portal",
                                  "links": [
                                    {
                                      "label": "User guide",
                                      "url": "https://portal.hyperspace.tools.sap/projects/dxp/documentation/User-Guide/Getting-Started/Overview"
                                    },
                                    {
                                      "label": "Extension catalog",
                                      "url": "https://portal.hyperspace.tools.sap/projects/dxp/documentation/Extension-Catalog"
                                    },
                                    {
                                      "label": "How to contribute",
                                      "url": "https://portal.hyperspace.tools.sap/projects/dxp/documentation/Extend-&-Contribute/Contribution-Guidelines/Overview"
                                    }
                                  ]
                                },
                                "layoutConfig": {
                                  "slot": "content"
                                }
                              },
                              {
                                "urlSuffix": "/microfrontends/quicklinks.js",
                                "context": {
                                  "title": "Find Guidance",
                                  "description": "Relevant links within the Hyperspace Academy",
                                  "links": [
                                    {
                                      "label": "Paved Roads",
                                      "url": "https://pages.github.tools.sap/hyperspace/academy/pavedroad/"
                                    },
                                    {
                                      "label": "Tools documentation",
                                      "url": "https://pages.github.tools.sap/hyperspace/academy/tools/"
                                    },
                                    {
                                      "label": "Service documentations",
                                      "url": "https://pages.github.tools.sap/hyperspace/academy/services/"
                                    },
                                    {
                                      "label": "Community-driven content",
                                      "url": "https://pages.github.tools.sap/hyperspace/academy/communitycontent/"
                                    }
                                  ]
                                },
                                "layoutConfig": {
                                  "slot": "content"
                                }
                              },
                              {
                                "urlSuffix": "/microfrontends/quicklinks.js",
                                "context": {
                                  "title": "Learn More",
                                  "description": "Relevant links within the Hyperspace SharePoint",
                                  "links": [
                                    {
                                      "label": "Development Platform offerings",
                                      "url": "https://sap.sharepoint.com/sites/124706/SitePages/Hyperspace-Development-Platform.aspx"
                                    },
                                    {
                                      "label": "What's next (roadmap)",
                                      "url": "https://sap.sharepoint.com/sites/124706/SitePages/What's-next-(roadmap).aspx"
                                    },
                                    {
                                      "label": "What's new (release notes)",
                                      "url": "https://sap.sharepoint.com/sites/124706/SitePages/What's-new-(release-notes).aspx"
                                    },
                                    {
                                      "label": "Join events & communities",
                                      "url": "https://sap.sharepoint.com/sites/124706/SitePages/Join-Events.aspx"
                                    }
                                  ]
                                },
                                "layoutConfig": {
                                  "slot": "content"
                                }
                              },
                              {
                                "urlSuffix": "/microfrontends/quicklinks.js",
                                "context": {
                                  "title": "Ask a Question",
                                  "description": "Link to search stack",
                                  "links": [
                                    {
                                      "label": "Help Center",
                                      "url": "https://portal.d1.hyperspace.tools.sap/projects/hyperspace-academy-experimental-laboratory/documentation/Home?modal=%2Fhelp&modalParams=%7B%22size%22:%22fullscreen%22%7D"
                                    }
                                  ]
                                },
                                "layoutConfig": {
                                  "slot": "content"
                                }
                              },
                              {
                                "urlSuffix": "/microfrontends/quicklinks.js",
                                "context": {
                                  "title": "Get Support",
                                  "description": "Links to the tools support channels",
                                  "links": [
                                    {
                                      "label": "Tool support channels",
                                      "url": "https://pages.github.tools.sap/hyperspace/academy/tools/"
                                    }
                                  ]
                                },
                                "layoutConfig": {
                                  "slot": "content"
                                }
                              }
                            ]
                          }
                        },
                        {
                          "pathSegment": "documentation",
                          "hideFromNav": true,
                          "url": "https://md-html.portal.{context.serviceProviderConfig.clusterHost}/#/",
                          "context": {
                            "projectId": "dxp"
                          },
                          "virtualTree": true
                        },
                        {
                          "pathSegment": "goldenPath",
                          "label": "{{goldenPath}}",
                          "url": "https://pages-ght.{context.serviceProviderConfig.clusterHost}/dxp/golden-path/",
                          "requiredIFramePermissions": {
                            "allow": ["clipboard-read", "clipboard-write", "fullscreen"]
                          },
                          "visibleForFeatureToggles": ["gp"],
                          "context": {
                            "pages": {
                              "login": "dxp",
                              "repoName": "golden-path"
                            }
                          },
                          "virtualTree": true,
                          "clientPermissions": {
                            "urlParameters": {
                              "url": {
                                "read": true,
                                "write": true
                              }
                            }
                          }
                        }
                      ]
                    }
                  ],
                  "texts": [
                    {
                      "locale": "",
                      "textDictionary": {
                        "home": "Home",
                        "hero_title": "Learn how Hyperspace Portal can help you build solutions",
                        "hero_content": "Explore projects, re-use templates and deploy your first components.",
                        "hero_alt": "Or you can...",
                        "hero_button_build": "Let's go build something",
                        "hero_button_docs": "Read Docs",
                        "learning": "Learning",
                        "dxp": "Hyperspace Portal",
                        "hyperspace": "Hyperspace",
                        "introduction": "Introduction",
                        "academy": "Academy",
                        "goldenPath": "Golden Path"
                      }
                    },
                    {
                      "locale": "en",
                      "textDictionary": {
                        "home": "Home",
                        "hero_title": "Learn how Hyperspace Portal can help you build solutions",
                        "hero_content": "Explore projects, re-use templates and deploy your first components.",
                        "hero_alt": "Or you can...",
                        "hero_button_build": "Let's go build something",
                        "hero_button_docs": "Read Docs",
                        "learning": "Learning",
                        "dxp": "Hyperspace Portal",
                        "hyperspace": "Hyperspace",
                        "introduction": "Introduction",
                        "academy": "Academy",
                        "goldenPath": "Golden Path"
                      }
                    },
                    {
                      "locale": "de",
                      "textDictionary": {
                        "home": "Home",
                        "hero_title": "Lerne, wie Hyperspace Portal bei der Erstellung von Lösungen helfen kann",
                        "hero_content": "Erkunde Projekte, verwende Templates und entwickle Deine ersten Komponenten",
                        "hero_alt": "oder...",
                        "hero_button_build": "Erschaffe etwas",
                        "hero_button_docs": "Lies die Dokumentation",
                        "learning": "Learning",
                        "dxp": "Hyperspace Portal",
                        "hyperspace": "Hyperspace",
                        "introduction": "Einführung",
                        "academy": "Akademie",
                        "goldenPath": "Golden Path"
                      }
                    }
                  ]
              }
          }
      }`
}

func GetValidJSON_organization_ui() string {
	return `{
        "name": "organization-ui",
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
                "entityType": "global",
                "pathSegment": "products",
                "label": "{{products}}",
                "hideSideNav": true,
                "icon": "product",
                "urlSuffix": "/{i18n.currentLocale}/#/products",
                "dxpOrder": 2,
                "navigationContext": "projects",
                "visibleForFeatureToggles": ["splitProjectByType"]
              },
              {
                "entityType": "global",
                "pathSegment": "experiments",
                "label": "{{experiments}}",
                "hideSideNav": true,
                "icon": "lab",
                "urlSuffix": "/{i18n.currentLocale}/#/experiments",
                "dxpOrder": 2.1,
                "navigationContext": "projects",
                "visibleForFeatureToggles": ["splitProjectByType"]
              },
              {
                "entityType": "global",
                "pathSegment": "projects",
                "label": "{{projects}}",
                "hideSideNav": true,
                "icon": "curriculum",
                "urlSuffix": "/{i18n.currentLocale}/#/projects",
                "dxpOrder": 2,
                "navigationContext": "projects",
                "visibleForFeatureToggles": ["!splitProjectByType"],
                "context": {
                  "_tmpUseBreadcrumbsForTitle": true
                },
                "children": [
                  {
                    "pathSegment": ":projectId",
                    "hideFromNav": true,
                    "navHeader": {
                      "useTitleResolver": true
                    },
                    "titleResolver": {
                      "request": {
                        "method": "GET",
                        "url": "${frameContext.accountSearchServiceApiUrl}?q=name:${projectId}",
                        "headers": {
                          "authorization": "Bearer ${token}"
                        }
                      },
                      "titlePropertyChain": "docs[0].displayName",
                      "prerenderFallback": false,
                      "fallbackTitle": "{{project}}",
                      "fallbackIcon": "curriculum"
                    },
                    "defineEntity": {
                      "id": "project",
                      "contextKey": "projectId",
                      "dynamicFetchId": "project",
                      "useBack": true,
                      "label": "{{project}}",
                      "pluralLabel": "{{projects}}",
                      "notFoundConfig": {
                        "entityListNavigationContext": "projects",
                        "sapIllusSVG": "Scene-NoSearchResults"
                      }
                    },
                    "context": {
                      "projectId": ":projectId"
                    },
                    "navigationContext": "project",
                    "children": [
                      {
                        "defineSlot": "main"
                      },
                      {
                        "defineSlot": ""
                      },
                      {
                        "defineSlot": "devopsMetrics",
                        "category": {
                          "label": "{{devopsMetrics}}",
                          "collapsible": false
                        }
                      },
                      {
                        "defineSlot": ""
                      },
                      {
                        "defineSlot": "settings",
                        "category": {
                          "label": "{{settings}}",
                          "collapsible": false
                        }
                      }
                    ]
                  }
                ]
              },
              {
                "entityType": "global",
                "pathSegment": "projects",
                "label": "{{projects}}",
                "hideFromNav": true,
                "hideSideNav": true,
                "urlSuffix": "/{i18n.currentLocale}/#/products",
                "navigationContext": "projects",
                "visibleForFeatureToggles": ["splitProjectByType"],
                "children": [
                  {
                    "pathSegment": ":projectId",
                    "hideFromNav": true,
                    "navHeader": {
                      "useTitleResolver": true
                    },
                    "titleResolver": {
                      "request": {
                        "method": "GET",
                        "url": "${frameContext.accountSearchServiceApiUrl}?q=name:${projectId}",
                        "headers": {
                          "authorization": "Bearer ${token}"
                        }
                      },
                      "titlePropertyChain": "docs[0].displayName",
                      "prerenderFallback": false,
                      "fallbackTitle": "{{project}}",
                      "fallbackIcon": "curriculum"
                    },
                    "defineEntity": {
                      "id": "project",
                      "contextKey": "projectId",
                      "dynamicFetchId": "project",
                      "useBack": true,
                      "label": "{{project}}",
                      "pluralLabel": "{{projects}}",
                      "notFoundConfig": {
                        "entityListNavigationContext": "projects",
                        "sapIllusSVG": "Scene-NoSearchResults"
                      }
                    },
                    "context": {
                      "projectId": ":projectId"
                    },
                    "navigationContext": "project",
                    "children": [
                      {
                        "defineSlot": "main"
                      },
                      {
                        "defineSlot": ""
                      },
                      {
                        "defineSlot": "devopsMetrics",
                        "category": {
                          "label": "{{devopsMetrics}}",
                          "collapsible": false
                        }
                      },
                      {
                        "defineSlot": ""
                      },
                      {
                        "defineSlot": "settings",
                        "category": {
                          "label": "{{settings}}",
                          "isGroup": true
                        }
                      }
                    ]
                  }
                ]
              },
              {
                "entityType": "global",
                "pathSegment": "teams",
                "label": "{{teams}}",
                "urlSuffix": "/{i18n.currentLocale}/#/teams",
                "hideSideNav": true,
                "icon": "group",
                "dxpOrder": 3,
                "navigationContext": "teams",
                "children": [
                  {
                    "pathSegment": ":teamId",
                    "hideFromNav": true,
                    "navHeader": {
                      "useTitleResolver": true
                    },
                    "titleResolver": {
                      "request": {
                        "method": "GET",
                        "url": "${frameContext.accountSearchServiceApiUrl}?q=name:${teamId}&fuzzy=false&fq=accountRole%3A%22Team%22",
                        "headers": {
                          "authorization": "Bearer ${token}"
                        }
                      },
                      "titlePropertyChain": "docs[0].displayName",
                      "prerenderFallback": false,
                      "fallbackTitle": "{{team}}",
                      "fallbackIcon": "group"
                    },
                    "defineEntity": {
                      "id": "team",
                      "contextKey": "teamId",
                      "dynamicFetchId": "team",
                      "useBack": true,
                      "label": "{{team}}",
                      "pluralLabel": "{{teams}}",
                      "notFoundConfig": {
                        "entityListNavigationContext": "teams",
                        "sapIllusSVG": "Scene-NoSearchResults"
                      }
                    },
                    "context": {
                      "teamId": ":teamId"
                    },
                    "navigationContext": "team",
                    "children": [
                      {
                        "defineSlot": "main"
                      },
                      {
                        "defineSlot": ""
                      },
                      {
                        "defineSlot": "settings",
                        "category": {
                          "label": "{{settings}}",
                          "collapsible": false
                        }
                      }
                    ]
                  }
                ]
              },
              {
                "pathSegment": "projects-create-dialog",
                "entityType": "global",
                "hideFromNav": true,
                "hideSideNav": true,
                "navigationContext": "projects",
                "context": {
                  "_tmpUseBreadcrumbsForTitle": true
                },
                "urlSuffix": "/{i18n.currentLocale}/#/projects-create-dialog"
              },
              {
                "pathSegment": "create-product",
                "entityType": "global",
                "hideFromNav": true,
                "hideSideNav": true,
                "label": "{{createProduct}}",
                "visibleForFeatureToggles": ["splitProjectByType"],
                "urlSuffix": "/{i18n.currentLocale}/#/create-project?type=product"
              },
              {
                "pathSegment": "create-experiment",
                "entityType": "global",
                "hideFromNav": true,
                "hideSideNav": true,
                "label": "{{createExperiment}}",
                "visibleForFeatureToggles": ["splitProjectByType"],
                "urlSuffix": "/{i18n.currentLocale}/#/create-project?type=experiment"
              },
              {
                "pathSegment": "create-project",
                "entityType": "global",
                "hideFromNav": true,
                "hideSideNav": true,
                "context": {
                  "_tmpUseBreadcrumbsForTitle": true
                },
                "urlSuffix": "/{i18n.currentLocale}/#/choose-project-type"
              },
              {
                "pathSegment": "create-project-details",
                "entityType": "global",
                "hideFromNav": true,
                "hideSideNav": true,
                "context": {
                  "_tmpUseBreadcrumbsForTitle": true
                },
                "urlSuffix": "/{i18n.currentLocale}/#/create-project"
              },
              {
                "pathSegment": "edit-project",
                "entityType": "project",
                "hideFromNav": true,
                "hideSideNav": true,
                "label": "Edit Project",
                "context": {
                  "_tmpUseBreadcrumbsForTitle": true
                },
                "urlSuffix": "/{i18n.currentLocale}/#/edit-project"
              },
              {
                "pathSegment": "edit-team",
                "entityType": "team",
                "hideFromNav": true,
                "hideSideNav": true,
                "label": "Edit Team",
                "context": {
                  "_tmpUseBreadcrumbsForTitle": true
                },
                "urlSuffix": "/{i18n.currentLocale}/#/edit-team"
              },
              {
                "pathSegment": "create-team",
                "entityType": "global",
                "hideFromNav": true,
                "hideSideNav": true,
                "label": "{{createTeam}}",
                "context": {
                  "_tmpUseBreadcrumbsForTitle": true
                },
                "urlSuffix": "/{i18n.currentLocale}/#/create-team"
              }
            ],
            "texts": [
              {
                "locale": "",
                "textDictionary": {
                  "projects": "Projects",
                  "project": "Project",
                  "products": "Products",
                  "experiments": "Experiments",
                  "createProduct": "Create Product",
                  "createExperiment": "Create Experiment",
                  "createTeam": "Create Team",
                  "teams": "Teams",
                  "team": "Team",
                  "settings": "Settings & Accounts",
                  "devopsMetrics": "DevOps Metrics"
                }
              },
              {
                "locale": "en",
                "textDictionary": {
                  "projects": "Projects",
                  "project": "Project",
                  "products": "Products",
                  "experiments": "Experiments",
                  "createProduct": "Create Product",
                  "createExperiment": "Create Experiment",
                  "createTeam": "Create Team",
                  "teams": "Teams",
                  "team": "Team",
                  "settings": "Settings & Accounts",
                  "devopsMetrics": "DevOps Metrics"
                }
              },
              {
                "locale": "de",
                "textDictionary": {
                  "projects": "Projekte",
                  "project": "Projekt",
                  "products": "Produkte",
                  "experiments": "Experimente",
                  "createProduct": "Produkt erstellen",
                  "createExperiment": "Experiment erstellen",
                  "createTeam": "Team erstellen",
                  "teams": "Teams",
                  "team": "Team",
                  "settings": "Einstellungen & Accounts",
                  "devopsMetrics": "DevOps Metrics"
                }
              }
            ]
          }
        }
      }`
}

func GetValidJSON_search_ui() string {
	return `{
        "name": "search-ui",
        "luigiConfigFragment": {
          "data": {
            "viewGroup": {
              "preloadSuffix": "/#/preload",
              "requiredIFramePermissions": {
                "allow": ["clipboard-read", "clipboard-write"],
                "sandbox": ["allow-forms"]
              }
            },
            "nodes": [
              {
                "entityType": "global",
                "pathSegment": "search",
                "hideFromNav": true,
                "showBreadcrumbs": false,
                "navHeader": {
                  "label": "Category"
                },
                "navigationContext": "search",
                "urlSuffix": "/#/search",
                "children": [
                  {
                    "label": "Projects {viewGroupData.projects}",
                    "icon": "curriculum",
                    "pathSegment": "projects",
                    "urlSuffix": "/#/search/projects",
                    "visibleForFeatureToggles": ["!splitProjectByType"],
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "label": "Products {viewGroupData.products}",
                    "icon": "product",
                    "pathSegment": "products",
                    "urlSuffix": "/#/search/products",
                    "visibleForFeatureToggles": ["splitProjectByType"],
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "hideFromNav": true,
                    "pathSegment": "projects",
                    "urlSuffix": "/#/search/products",
                    "visibleForFeatureToggles": ["splitProjectByType"],
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "label": "Experiments {viewGroupData.experiments}",
                    "icon": "lab",
                    "pathSegment": "experiments",
                    "urlSuffix": "/#/search/experiments",
                    "visibleForFeatureToggles": ["splitProjectByType"],
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "label": "Components {viewGroupData.components}",
                    "icon": "course-book",
                    "pathSegment": "components",
                    "urlSuffix": "/#/search/components",
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "label": "Teams {viewGroupData.teams}",
                    "icon": "group",
                    "pathSegment": "teams",
                    "urlSuffix": "/#/search/teams",
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "label": "Documentation {viewGroupData.techdocs}",
                    "icon": "curriculum",
                    "pathSegment": "techdocs",
                    "urlSuffix": "/#/search/techdocs",
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        },
                        "url": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "label": "API Documentation (draft)",
                    "icon": "curriculum",
                    "pathSegment": "apidocs",
                    "urlSuffix": "/#/search/apidocs",
                    "visibleForFeatureToggles": ["enable-api-docs-search"],
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  },
                  {
                    "label": "Users {viewGroupData.users}",
                    "icon": "group",
                    "pathSegment": "users",
                    "urlSuffix": "/#/search/users",
                    "clientPermissions": {
                      "urlParameters": {
                        "q": {
                          "read": true,
                          "write": true
                        }
                      }
                    }
                  }
                ]
              }
            ]
          }
        }
      }`
}
