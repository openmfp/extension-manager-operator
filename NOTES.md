We need to create a Golang re-use module which implements validation functions for the ContentConfiguration.

This module will be used by the controllers in the Extension Manager Operator; and also by the Extension Manager Service, where this module will be used to validate the local JSON before being fed to Jukebox (so developers can receive the same errors locally than they'd receive once the ContentConfiguration is deployed on-cluster).



Requirements:
Create an exposed/public package in the extension-content-operator https://github.com/openmfp/extension-content-operator to validate Content configurations

The package will contain the following features:

validate input configuration as JSON or YAML for correctness
Generate configurationResult that will be stored in ContentConfiguration.status.configurationResult
Return descriptive error messages in case validation fails
Support inlineConfiguration
Implement a json-schema validation for the content configuration that is reflected in go-structs


An example of the usage of json schemas can be found in

https://github.tools.sap/dxp/json-schemas

File: go.mod

```go
module github.tools.sap/dxp/json-schemas

go 1.22.2

require (
	github.com/invopop/jsonschema v0.12.0
	github.tools.sap/dxp/metadata-registry-operator v0.0.0-20240212142158-393bffc1d6a3
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/evanphx/json-patch/v5 v5.8.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic-models v0.6.9-0.20230804172637-c7be7c783f49 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/oauth2 v0.17.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.29.1 // indirect
	k8s.io/apimachinery v0.29.1 // indirect
	k8s.io/client-go v0.29.1 // indirect
	k8s.io/klog/v2 v2.110.1 // indirect
	k8s.io/kube-openapi v0.0.0-20231010175941-2dd684a91f00 // indirect
	k8s.io/utils v0.0.0-20240102154912-e7106e64919e // indirect
	sigs.k8s.io/controller-runtime v0.17.1 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
````


File: main.go

```go
package main

import (
	"os"

	"github.com/invopop/jsonschema"
	"github.tools.sap/dxp/metadata-registry-operator/pkg/service/parser"
)

func main() {
	err := os.Mkdir("dist", 0750)
	if !os.IsExist(err) && err != nil {
		panic(err)
	}

	createJson(&parser.ProjectMetadata{}, "dist/project-metadata.json")

	createJson(&parser.TeamMetadata{}, "dist/team-metadata.json")

	createJson(&parser.ComponentMetadata{}, "dist/component-metadata.json")

	createJson(&parser.BoundedContextMetadata{}, "dist/bounded-context.json")

	createJson(&parser.APIDefinitionMetadata{}, "dist/api-definition.json")
}

func createJson(v any, path string) {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	s := r.Reflect(v)
	json, err := s.MarshalJSON()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, json, 0644)
	if err != nil {
		panic(err)
	}
}
FooterSAP SE
GitHub@SAP Docs
SNOW Support Ticket
Maintenance & Upgrade Roadmap
FAQ (General)
FAQ (github.tools)
SAP SE
SAP SE

```

File: renovate.json

```json
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": [
    "local>dxp/.github:renovate-config"
  ]
}
```


using https://github.com/invopop/jsonschema