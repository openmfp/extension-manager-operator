# Info

Useful information about the project.

# Hints

## Reference codebase

Use the following [https://github.com/openmfp/account-operator](https://github.com/openmfp/account-operator)

## TDD: Test Driven Development

Use `testify` instead of `ginko`.

Current impl at [https://github.com/openmfp/account-operator](https://github.com/openmfp/account-operator) uses a hybrid
approach
of `kubebuilder` and `operator-sdk` CLIs.

## CD: Continuous Deployment

Don't use `kustomize`

# Drawings Collaboration

Use the following Miro Board for collaboration on
drawings: [https://miro.com/app/board/uXjVKFWUajY=/](https://miro.com/app/board/uXjVKFWUajY=/)

## Operator Stub

[https://github.com/openmfp/extension-content-operator/pull/4/files](https://github.com/openmfp/extension-content-operator/pull/4/files)

## Onboarding

On-prem docs link [https://portal.hyperspace.tools.sap/projects/dxp/team-docs/Team/Onboarding](https://portal.hyperspace.tools.sap/projects/dxp/team-docs/Team/Onboarding)

## Deployment

The `openMFT` deployment can be found at the following link: [https://portal.d1.openmfp.dxp.k8s.ondemand.com/home/overview](https://portal.d1.openmfp.dxp.k8s.ondemand.com/home/overview)

## Public Docker Registry

Registry pointer: [https://hub.docker.com/repository/docker/norbixaccenture/extension-content-operator/general](https://hub.docker.com/repository/docker/norbixaccenture/extension-content-operator/general)

Pseudo code for building and pushing the image:

```text
make docker-build docker-push IMG=norbixaccenture/extension-content-operator:latest

# Use `local` preix for local testing and `norbixaccenture` for collaboration

docker images
REPOSITORY                                   TAG       IMAGE ID       CREATED          SIZE
local/extension-content-operator             latest    b5cef5f8d92d   17 minutes ago   58.9MB
norbixaccenture/extension-content-operator   latest    b5cef5f8d92d   17 minutes ago   58.9MB
```

## Developer Workflow

Pseudo code for developer workflow:

```text
k apply -f config/crd/crd.yaml
k get crds
go run ./main.go operator

{"level":"info","service":"/home/norbix/go/pkg/mod/github.com/openmfp/golang-commons@v0.32.0/logger/logger.go","time":"2024-05-22T21:19:59+02:00","caller":"/home/norbix/go/src/corpo/extension-content-operator/cmd/operator.go:129","message":"starting manager"}
{"level":"info","service":"/home/norbix/go/pkg/mod/github.com/openmfp/golang-commons@v0.32.0/logger/logger.go","component":"controller-runtime","v":0,"logger":"controller-runtime/metrics","time":"2024-05-22T21:19:59+02:00","caller":"/home/norbix/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.2/pkg/metrics/server/server.go:208","message":"Starting metrics server"}
{"level":"info","service":"/home/norbix/go/pkg/mod/github.com/openmfp/golang-commons@v0.32.0/logger/logger.go","component":"controller-runtime","name":"health probe","addr":"[::]:8081","v":0,"time":"2024-05-22T21:19:59+02:00","caller":"/home/norbix/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.2/pkg/manager/server.go:83","message":"starting server"}
{"level":"info","service":"/home/norbix/go/pkg/mod/github.com/openmfp/golang-commons@v0.32.0/logger/logger.go","component":"controller-runtime","v":0,"logger":"controller-runtime/metrics","bindAddress":":8080","secure":false,"time":"2024-05-22T21:19:59+02:00","caller":"/home/norbix/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.2/pkg/metrics/server/server.go:247","message":"Serving metrics server"}
{"level":"info","service":"/home/norbix/go/pkg/mod/github.com/openmfp/golang-commons@v0.32.0/logger/logger.go","component":"controller-runtime","controller":"ContentConfigurationReconciler","controllerGroup":"cache.core.openmfp.io","controllerKind":"ContentConfiguration","v":0,"source":"kind source: *v1alpha1.ContentConfiguration","time":"2024-05-22T21:19:59+02:00","caller":"/home/norbix/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.2/pkg/internal/controller/controller.go:173","message":"Starting EventSource"}
{"level":"info","service":"/home/norbix/go/pkg/mod/github.com/openmfp/golang-commons@v0.32.0/logger/logger.go","component":"controller-runtime","controller":"ContentConfigurationReconciler","controllerGroup":"cache.core.openmfp.io","controllerKind":"ContentConfiguration","v":0,"time":"2024-05-22T21:19:59+02:00","caller":"/home/norbix/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.2/pkg/internal/controller/controller.go:181","message":"Starting Controller"}
{"level":"info","service":"/home/norbix/go/pkg/mod/github.com/openmfp/golang-commons@v0.32.0/logger/logger.go","component":"controller-runtime","controller":"ContentConfigurationReconciler","controllerGroup":"cache.core.openmfp.io","controllerKind":"ContentConfiguration","v":0,"worker count":10,"time":"2024-05-22T21:19:59+02:00","caller":"/home/norbix/go/pkg/mod/sigs.k8s.io/controller-runtime@v0.18.2/pkg/internal/controller/controller.go:215","message":"Starting workers"}
```