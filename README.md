# extension-content-operator

## About extension-content-operator

The *extension-content-operator* implements implements the lifecycle management of a Kubernetes CRD `ContentConfiguration` resource, which is the API for configuration of Extensions in openMFP.

For reference, see the [RFC for openMFP Extension Management - CDM Processing](https://github.com/openmfp/architecture/pull/2/files?short_path=8a071a3#diff-8a071a31a02919a613572237f1e968fe02b9cf7d350c2cf796ba6b35495ec09b).

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
task docker-build docker-push IMG=<some-registry>/extension-content-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
task install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
task deploy IMG=<some-registry>/extension-content-operator:tag
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
import "github.com/openmfp/extension-content-operator/pkg/validation"
```

Example usage:

```go
package main

import (
    "fmt"
    "github.com/openmfp/extension-content-operator/pkg/validation"
)

func main() {
    cC := validation.NewContentConfiguration()

    schema := []byte(`{
        "$schema": "https://json-schema.org/draft/2020-12/schema",
        "$id": "https://example.com/schema.json",
        "type": "object",
        "properties": {
            "name": { "type": "string" }
        },
        "required": ["name"]
    }`)

    input := []byte(`{ "name": "example" }`)
    contentType := "json"

    result, err := cC.Validate(schema, input, contentType)
    if err != nil {
        fmt.Println("Validation failed:", err)
    } else {
        fmt.Println("Validation succeeded:", result)
    }
}
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
docker build . --no-cache --tag local-extension-content-operator:$IMG_TAG

# load image to kind
kind load docker-image local-extension-content-operator:$IMG_TAG

# apply CRDS
kubectl apply -f chart/crds/core.openmfp.io_contentconfigurations.yaml

# change in imagePullPolicy in chart/templates/deployment.yaml
imagePullPolicy: IfNotPresent

# apply chart with test configuration
helm template -f ./chart/test-values.yaml extension-content-operator --include-crds ./chart/ | kubectl apply -f -

# create sample resources
kubectl apply -f config/samples/v1alpha1_contentconfiguration.yaml

# cleanup
kubectl delete -f config/samples/v1alpha1_contentconfiguration.yaml
helm template -f ./chart/test-values.yaml extension-content-operator ./chart/ --include-crds | kubectl delete -f -
kubectl delete -f chart/crds/core.openmfp.io_contentconfigurations.yaml
docker image rm local-extension-content-operator:test
kind delete cluster
```

## Support, Feedback, Contributing

This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/openmfp/extension-content-operator/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/openmfp/extension-content-operator/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](https://github.com/openmfp/extension-content-operator/.github/blob/main/CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright (20xx-)20xx SAP SE or an SAP affiliate company and *extension-content-operator* contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/<your-project>).
