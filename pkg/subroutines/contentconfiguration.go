package subroutines

import (
	"context"
	"fmt"
	"github.com/openmfp/extension-content-operator/api/v1alpha1"
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

func (r *ContentConfigurationSubroutine) Finalizers() []string { // coverage-ignore
	return []string{"contentconfiguration.core.openmfp.io/finalizer"}
}

func (r *ContentConfigurationSubroutine) Process(
	ctx context.Context, runtimeObj lifecycle.RuntimeObject,
) (ctrl.Result, errors.OperatorError) {
	instance := runtimeObj.(*v1alpha1.ContentConfiguration)

	//if instance == nil {
	//	instance = &v1alpha1.ContentConfiguration{}
	//}

	fmt.Println("###", instance)

	//if instance.Spec.InlineConfiguration.Url != "" {
	//	instance.Status.ConfigurationResult.Url = "1213"
	//} else {
	//	instance.Status.ConfigurationResult.Url = "123"
	//}

	instance.Spec.InlineConfiguration.Content = "this is inline config"
	instance.Status.ConfigurationResult = "this is config result"

	err := r.client.Update(ctx, instance)
	if err != nil {
		fmt.Println("### err", err)
		return ctrl.Result{}, errors.NewOperatorError(err, true, true)
	}

	err = r.client.Status().Update(ctx, instance)
	if err != nil {
		return ctrl.Result{}, errors.NewOperatorError(err, true, true)
	}

	return ctrl.Result{}, nil
}
