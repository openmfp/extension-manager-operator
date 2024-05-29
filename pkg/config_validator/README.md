# Overview

WIP

# Table of Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Playbooks](#playbooks)
  - [How to generate JSON schema from go struct playbook](#how-to-generate-json-schema-from-go-struct-playbook)

# Getting Started

## Prerequisites

- Go SDK 1.22 [download](https://golang.org/dl/)
- IDE (e.g. IntelliJ IDEA) [download](https://www.jetbrains.com/idea/download/)

# Usage

In order to use the shared package within the project please check the following code snippets within [doc.go](./doc.go) file. 

## eXtreme Programming (XP) practices

The unit tests can be triggered by running the following command:

```bash
go test -v ./... -cover
```
# Playbooks

## How to generate JSON schema from go struct playbook

In order to create a `strict` JSON schema from a go struct, the following steps need to be followed: 

```text
go test -v -run ^Test_createJson$ ./...

{"$schema":"https://json-schema.org/draft/2020-12/schema","$id":"https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration","$defs":{"LuigiConfigData":{"properties":{"nodes":{"items":{"$ref":"#/$defs/Node"},"type":"array"}},"additionalProperties":false,"type":"object","required":["nodes"]},"LuigiConfigFragment":{"properties":{"data":{"$ref":"#/$defs/LuigiConfigData"}},"additionalProperties":false,"type":"object"},"Node":{"properties":{"entityType":{"type":"string"},"pathSegment":{"type":"string"},"label":{"type":"string"},"icon":{"type":"string"}},"additionalProperties":false,"type":"object"}},"properties":{"name":{"type":"string"},"luigiConfigFragment":{"items":{"$ref":"#/$defs/LuigiConfigFragment"},"type":"array"}},"additionalProperties":false,"type":"object"}
```

By removing the `omitempty` annotation to the struct, the JSON schema will be strictly generated.

### With `omitempty`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration",
  "$defs": {
    "LuigiConfigData": {
      "properties": {
        "nodes": {
          "items": {
            "$ref": "#/$defs/Node"
          },
          "type": "array"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": ["nodes"]
    },
    "LuigiConfigFragment": {
      "properties": {
        "data": {
          "$ref": "#/$defs/LuigiConfigData"
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "Node": {
      "properties": {
        "entityType": {
          "type": "string"
        },
        "pathSegment": {
          "type": "string"
        },
        "label": {
          "type": "string"
        },
        "icon": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  },
  "properties": {
    "name": {
      "type": "string"
    },
    "luigiConfigFragmenty": {
      "items": {
        "$ref": "#/$defs/LuigiConfigFragment"
      },
      "type": "array"
    }
  },
  "additionalProperties": false,
  "type": "object"
}
```

### Without `omitempty`

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://github.com/openmfp/extension-content-operator/pkg/config_validator/content-configuration",
  "$defs": {
    "LuigiConfigData": {
      "properties": {
        "nodes": {
          "items": {
            "$ref": "#/$defs/Node"
          },
          "type": "array"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": ["nodes"]
    },
    "LuigiConfigFragment": {
      "properties": {
        "data": {
          "$ref": "#/$defs/LuigiConfigData"
        }
      },
      "additionalProperties": false,
      "type": "object"
    },
    "Node": {
      "properties": {
        "entityType": {
          "type": "string"
        },
        "pathSegment": {
          "type": "string"
        },
        "label": {
          "type": "string"
        },
        "icon": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object"
    }
  },
  "properties": {
    "name": {
      "type": "string"
    },
    "luigiConfigFragmenty": {
      "items": {
        "$ref": "#/$defs/LuigiConfigFragment"
      },
      "type": "array"
    }
  },
  "additionalProperties": false,
  "type": "object",
  "required": ["name", "luigiConfigFragmenty"]
}
```

### The delta between the two JSON schemas

Key Differences:

1. Presence of the required property:
   - With omitempty: The required property is missing at the top level.
   - Without omitempty: The required property is present at the top level with the values ["name", "luigiConfigFragmenty"].

That's why the JSON schema generated with the `omitempty` annotation is not strict.

In code the JSON schema for the [model.go](./model.go) is generated without the `omitempty` annotation.
