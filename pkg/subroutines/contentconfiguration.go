package subroutines

import (
	"context"

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

func (r *ContentConfigurationSubroutine) Finalize(ctx context.Context, runtimeObj lifecycle.RuntimeObject) (ctrl.Result, errors.OperatorError) {
	return ctrl.Result{}, nil
}

func (r *ContentConfigurationSubroutine) Finalizers() []string { // coverage-ignore
	return []string{"contentconfiguration.core.openmfp.io/finalizer"}
}

func (r *ContentConfigurationSubroutine) Process(ctx context.Context, runtimeObj lifecycle.RuntimeObject) (ctrl.Result, errors.OperatorError) {
	// TODO: processing logic
	// instance := runtimeObj.(*v1alpha1.ContentConfiguration)

	// logger.Log.Info("Processing ContentConfiguration", "namespace", instance.Namespace, "name", instance.Name)

	return ctrl.Result{}, nil
}
