> [!WARNING]
> This Repository is under development and not ready for productive use. It is in an alpha stage. That means APIs and concepts may change on short notice including breaking changes or complete removal of apis.

# openMFP - Extension Manager Operator

## Description

The extension-manager-operator implements the lifecycle management of a Kubernetes CRD `ContentConfiguration` resource, which is a Kubernetes Resource/API for configuration of Micro Frontends in openMFP.

For reference, see the [RFC for openMFP Extension Management - CDM Processing](https://github.com/openmfp/architecture/blob/main/rfc/002-extension-content-configuration-processing.md).

## Features
- Support for inline and remote content configurations. 
- Validation of content configuration and generation of a JSON Schema that can be used by contributors to validate their content configuration.
- Services to allow validation of content configuration at runtime while developing a micro frontend on the developers system.
- Ability to provide validation feedback while keeping the last validated content configuration.

## Getting Started
For running OpenMFP locally checkout our [getting started guide](https://openmfp.github.io/openmfp.org/docs/getting-started). The extension-manager-operator can be deployed on a kubernetes cluster using the helm-chart [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator) and for CRDs [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator-crds).

## Releasing

The release is performed automatically through a GitHub Actions Workflow. New Versions will be updated in the helm-chart of the extension-manager-operator located [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator). There is a separate helm chart for the extension-manager-operator CRDS located [here](https://github.com/openmfp/helm-charts/tree/main/charts/extension-manager-operator-crds).

## Support, Feedback, Contributing
This project is open to feature requests/suggestions, bug reports etc. via [GitHub issues](https://github.com/openmfp/extension-manager-operator/issues). Contribution and feedback are encouraged and always welcome. For more information about how to contribute, the project structure, as well as additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Security / Disclosure
If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/openmfp/extension-manager-operator/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct
We as members, contributors, and leaders pledge to make participation in our community a harassment-free experience for everyone. By participating in this project, you agree to abide by its [Code of Conduct](CODE_OF_CONDUCT.md) at all times.

## Licensing

Copyright 2024 SAP SE or an SAP affiliate company and openMFP contributors. Please see our [LICENSE](LICENSE) for copyright and license information. Detailed information including third-party components and their licensing/copyright information is available [via the REUSE tool](https://api.reuse.software/info/github.com/openmfp/extension-manager-operator).

