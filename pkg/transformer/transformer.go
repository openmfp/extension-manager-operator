package transformer

import (
	"github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/pkg/validation"
)

type ContentConfigurationTransformer interface {
	Transform(contentConfiguration *validation.ContentConfiguration, instance *v1alpha1.ContentConfiguration) error
}
