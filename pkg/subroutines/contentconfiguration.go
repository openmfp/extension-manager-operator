package subroutines

import (
	"context"
	"net/http"

	"github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/pkg/httpclient"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/errors"
	"github.com/openmfp/golang-commons/logger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ContentConfigurationSubroutineName                 = "ContentConfigurationSubroutine"
	ContentConfigurationSubroutineFinalizer            = "contentconfiguration.core.openmfp.io/finalizer"
	ContentConfigurationOwnerLabel                     = "contentconfiguration.core.openmfp.io/owner"
	ContentConfigurationOwnerContentConfigurationLabel = "contentconfiguration.core.openmfp.io/owner-namespace"
	ContentConfigurationNamePrefix                     = "contentconfiguration-"
)

type ContentConfigurationSubroutine struct {
	client     client.Client
	httpClient httpclient.Service
	log        *logger.Logger
}

func NewContentConfigurationSubroutine(client client.Client, log *logger.Logger) *ContentConfigurationSubroutine {
	return &ContentConfigurationSubroutine{
		client:     client,
		httpClient: httpclient.NewService(),
		log:        log,
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
		bytes, err := r.httpClient.Do(http.MethodGet, instance.Spec.RemoteConfiguration.URL, nil)
		if err != nil {
			r.log.Err(err).Msg("failed to fetch remote configuration")

			return ctrl.Result{}, errors.NewOperatorError(err, true, true)
		}
		rawConfig = bytes
	}

	// TODO replace it with validation function
	validatedConfig := string(rawConfig)

	instance.Status.ConfigurationResult = validatedConfig

	err := r.client.Status().Update(ctx, instance)
	if err != nil {
		r.log.Err(err).Msg("failed to update instance status")

		return ctrl.Result{}, errors.NewOperatorError(err, true, true)
	}

	return ctrl.Result{}, nil
}
