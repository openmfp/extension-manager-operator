package subroutines

import (
	"context"
	"net/http"

	"github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/pkg/httpclient"
	"github.com/openmfp/golang-commons/logger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/errors"
)

const (
	ContentConfigurationSubroutineName                 = "ContentConfigurationSubroutine"
	ContentConfigurationSubroutineFinalizer            = "contentconfiguration.core.openmfp.io/finalizer"
	ContentConfigurationOwnerLabel                     = "contentconfiguration.core.openmfp.io/owner"
	ContentConfigurationOwnerContentConfigurationLabel = "contentconfiguration.core.openmfp.io/owner-namespace"
	ContentConfigurationNamePrefix                     = "contentconfiguration-"
)

type ContentConfigurationSubroutine struct {
	client client.Client
}

func NewContentConfigurationSubroutine(client client.Client) *ContentConfigurationSubroutine {
	return &ContentConfigurationSubroutine{client: client}
}

func (r *ContentConfigurationSubroutine) GetName() string {
	return ContentConfigurationSubroutineName
}

func (r *ContentConfigurationSubroutine) Finalize(
	ctx context.Context,
	runtimeObj lifecycle.RuntimeObject) (ctrl.Result, errors.OperatorError) {
	return ctrl.Result{}, nil
}

func (r *ContentConfigurationSubroutine) Finalizers() []string {
	return []string{"contentconfiguration.core.openmfp.io/finalizer"}
}

func (r *ContentConfigurationSubroutine) Process(
	ctx context.Context, runtimeObj lifecycle.RuntimeObject,
) (ctrl.Result, errors.OperatorError) {
	log := logger.LoadLoggerFromContext(ctx)

	instance := runtimeObj.(*v1alpha1.ContentConfiguration)

	var rawConfig []byte
	// InlineConfiguration has higher priority than RemoteConfiguration
	switch {
	case instance.Spec.InlineConfiguration.Content != "":
		rawConfig = []byte(instance.Spec.InlineConfiguration.Content)
	case instance.Spec.RemoteConfiguration.URL != "":
		bytes, err, retry := httpclient.NewService().Do(http.MethodGet, instance.Spec.RemoteConfiguration.URL, nil)
		if err != nil {
			log.Err(err).Msg("failed to fetch remote configuration")

			return ctrl.Result{}, errors.NewOperatorError(err, retry, true)
		}
		rawConfig = bytes
	default:
		return ctrl.Result{}, errors.NewOperatorError(errors.New("no configuration provided"), false, true)
	}

	// TODO replace it with validation function
	validatedConfig := string(rawConfig)

	instance.Status.ConfigurationResult = validatedConfig

	return ctrl.Result{}, nil
}
