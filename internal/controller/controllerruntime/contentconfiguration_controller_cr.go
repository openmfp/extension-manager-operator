/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllerruntime

import (
	"context"
	"net/http"

	openmfpconfig "github.com/platform-mesh/golang-commons/config"
	"github.com/platform-mesh/golang-commons/controller/lifecycle/builder"
	"github.com/platform-mesh/golang-commons/controller/lifecycle/controllerruntime"
	"github.com/platform-mesh/golang-commons/controller/lifecycle/subroutine"
	"github.com/platform-mesh/golang-commons/logger"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/internal/config"
	"github.com/openmfp/extension-manager-operator/internal/controller"
	"github.com/openmfp/extension-manager-operator/pkg/subroutines"
	"github.com/openmfp/extension-manager-operator/pkg/validation"
)

var contentConfigurationReconcilerName = "ContentConfigurationReconcilerCR"

// ContentConfigurationReconcilerCR ContentConfigurationReconciler reconciles a ContentConfiguration object
type ContentConfigurationReconcilerCR struct {
	lifecycle *controllerruntime.LifecycleManager
}

func NewContentConfigurationReconcilerCR(log *logger.Logger, mgr ctrl.Manager, cfg config.OperatorConfig) *ContentConfigurationReconcilerCR {
	var subs []subroutine.Subroutine
	if cfg.Subroutines.ContentConfiguration.Enabled {
		subs = append(subs, subroutines.NewContentConfigurationSubroutine(validation.NewContentConfiguration(), http.DefaultClient))
	}
	return &ContentConfigurationReconcilerCR{lifecycle: builder.
		NewBuilder(controller.OperatorName, contentConfigurationReconcilerName, subs, log).
		WithSpreadingReconciles().
		WithConditionManagement().
		BuildControllerRuntime(mgr.GetClient())}
}

func (r *ContentConfigurationReconcilerCR) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.lifecycle.Reconcile(ctx, req, &v1alpha1.ContentConfiguration{})
}

func (r *ContentConfigurationReconcilerCR) SetupWithManager(mgr ctrl.Manager, cfg *openmfpconfig.CommonServiceConfig, log *logger.Logger, eventPredicates ...predicate.Predicate) error {
	return r.lifecycle.SetupWithManager(mgr, cfg.MaxConcurrentReconciles, contentConfigurationReconcilerName, &v1alpha1.ContentConfiguration{}, cfg.DebugLabelValue, r, log, eventPredicates...)
}
