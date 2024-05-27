package subroutines

import (
	"context"
	"github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/pkg/retryHttpClient"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/errors"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	ContentConfigurationSubroutineName                 = "ContentConfigurationSubroutine"
	ContentConfigurationSubroutineFinalizer            = "contentconfiguration.core.openmfp.io/finalizer"
	ContentConfigurationOwnerLabel                     = "contentconfiguration.core.openmfp.io/owner"
	ContentConfigurationOwnerContentConfigurationLabel = "contentconfiguration.core.openmfp.io/owner-namespace"
	ContentConfigurationNamePrefix                     = "contentconfiguration-"
)

type ContentConfigurationSubroutine struct {
	client          client.Client
	retryHttpClient retryHttpClient.Service
}

func NewContentConfigurationSubroutine(client client.Client) *ContentConfigurationSubroutine {
	return &ContentConfigurationSubroutine{
		client:          client,
		retryHttpClient: retryHttpClient.New(5, 1*time.Second, 5*time.Second),
	}
}

func (r *ContentConfigurationSubroutine) GetName() string {
	return ContentConfigurationSubroutineName
}

func (r *ContentConfigurationSubroutine) Finalize(
	ctx context.Context,
	runtimeObj lifecycle.RuntimeObject) (ctrl.Result, errors.OperatorError) {
	return ctrl.Result{}, nil
}

func (r *ContentConfigurationSubroutine) Finalizers() []string { // coverage-ignore
	return []string{"contentconfiguration.core.openmfp.io/finalizer"}
}

func (r *ContentConfigurationSubroutine) Process(
	ctx context.Context, runtimeObj lifecycle.RuntimeObject,
) (ctrl.Result, errors.OperatorError) {
	instance := runtimeObj.(*v1alpha1.ContentConfiguration)

	var rawConfig []byte
	// InlineConfiguration has higher priority than RemoteConfiguration
	if instance.Spec.InlineConfiguration.Content != "" {
		rawConfig = []byte(instance.Spec.InlineConfiguration.Content)
	} else {
		bytes, err := r.retryHttpClient.Do(instance.Spec.RemoteConfiguration.URL, http.MethodGet, nil)
		if err != nil {
			ctrl.Log.Error(err, "failed to fetch remote configuration")

			return ctrl.Result{}, errors.NewOperatorError(err, true, true)
		}
		rawConfig = bytes
	}

	// TODO replace it with validation function
	validatedConfig := string(rawConfig)

	instance.Status.ConfigurationResult = validatedConfig

	err := r.client.Status().Update(ctx, instance)
	if err != nil {
		ctrl.Log.Error(err, "failed to update instance status")
		return ctrl.Result{}, errors.NewOperatorError(err, true, true)
	}

	return ctrl.Result{}, nil
}
