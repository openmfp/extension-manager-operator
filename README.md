> [!WARNING]
> This Repository is under development and not ready for productive use. It is in an alpha stage. That means APIs and concepts may change on short notice including breaking changes or complete removal of apis.

# OpenMFP - Extension Manager Operator

## Description

The extension-manager-operator implements the lifecycle management of a Kubernetes CRD `ContentConfiguration` resource, which is a Kubernetes Resource/API for configuration of Micro Frontends in OpenMFP.

For reference, see the [RFC for OpenMFP Extension Management - CDM Processing](https://github.com/openmfp/architecture/blob/main/rfc/002-extension-content-configuration-processing.md).

## Features
- Support for inline and remote content configurations. 
- Validation of content configuration and generation of a JSON Schema that can be used by contributors to validate their content configuration.
- Services to allow validation of content configuration at runtime while developing a micro frontend on the developers system.
- Ability to provide validation feedback while keeping the last validated content configuration.

## Getting Started
For running OpenMFP locally checkout our [getting started guide](https://openmfp.github.io/openmfp.org/docs/getting-started). The extension-manager-operator can be deployed on a kubernetes cluster using the helm-chart [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator) and for CRDs [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator-crds).

## Releasing

The release is performed automatically through a GitHub Actions Workflow. New Versions will be updated in the helm-chart of the extension-manager-operator located [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator). There is a separate helm chart for the extension-manager-operator CRDS located [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator-crds).

## Requirements

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
task install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
task deploy IMG=<some-registry>/extension-manager-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
task uninstall
```

**UnDeploy the controller from the cluster:**

```sh
task undeploy
```

### Run tests
```sh
task test
task cover
```
### Debug locally
```sh
task generate
task manifests
task install
task run
```

## Using a `validation` library

To install the validation library, add the following import statement in your Go project:

```go
import "github.com/openmfp/extension-manager-operator/pkg/validation"
```

Example usage:

```go
package main

import (
    "fmt"
    "github.com/openmfp/extension-manager-operator/pkg/validation"
)

func main() {
    cC := validation.NewContentConfiguration()

    input := []byte(`{ "name": "example" }`)
    contentType := "json"

    result, err := cC.Validate(input, contentType)
    if err != nil {
        fmt.Println("Validation failed:", err)
    } else {
        fmt.Println("Validation succeeded:", result)
    }
}
```

## Using a `/validate` HTTP endpoint

```shell
# run with 'server' argument
go run main.go server

# validate docs/assets/test.json local file
curl http://localhost:8088/validate -X POST -d @docs/assets/test.json   -H "Content-Type: application/json"
```


### Debug Helm chart locally

```sh
# create local KIND cluster
cat <<EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
EOF
kind create cluster --config ./kind-config.yaml

IMG_TAG=0.16.0

# build docker local chart image
docker build . --no-cache --tag local-extension-manager-operator:$IMG_TAG

# load image to kind
kind load docker-image local-extension-manager-operator:$IMG_TAG

# apply CRDS
kubectl apply -f chart/crds/core.openmfp.io_contentconfigurations.yaml

# change in imagePullPolicy in chart/templates/deployment.yaml
imagePullPolicy: IfNotPresent

# apply chart with test configuration
helm template -f ./chart/test-values.yaml extension-manager-operator --include-crds ./chart/ | kubectl apply -f -

# create sample resources
kubectl apply -f config/samples/v1alpha1_contentconfiguration.yaml

# cleanup
kubectl delete -f config/samples/v1alpha1_contentconfiguration.yaml
helm template -f ./chart/test-values.yaml extension-manager-operator ./chart/ --include-crds | kubectl delete -f -
kubectl delete -f chart/crds/core.openmfp.io_contentconfigurations.yaml
docker image rm local-extension-manager-operator:test
kind delete cluster
```

### Updating the JSONSchema

The JSON Schema used to validate `ContentConfiguration` resources is automatically generated from the structs in `pkg/validation/model.go`. To update it, modify model.go and run `task schemagen` or `task generate`.

## Support, Feedback, Contributing
This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/openmfp/extension-manager-operator/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/openmfp/extension-manager-operator/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Contributing

Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file in this repository for instructions on how to contribute to OpenMFP.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2024 SAP SE or an SAP affiliate company and OpenMFP contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/openmfp/extension-manager-operator).
