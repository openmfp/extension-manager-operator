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

package controller

import (
	"context"
	"net/http"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	cachev1alpha1 "github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/internal/config"
	"github.com/openmfp/extension-manager-operator/pkg/subroutines"
	"github.com/openmfp/extension-manager-operator/pkg/validation"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/logger"
)

var (
	operatorName                       = "extension-manager-operator"
	contentConfigurationReconcilerName = "ContentConfigurationReconciler"
)

// ContentConfigurationReconciler reconciles a ContentConfiguration object
type ContentConfigurationReconciler struct {
	lifecycle *lifecycle.LifecycleManager
}

func NewContentConfigurationReconciler(log *logger.Logger, mgr ctrl.Manager, cfg config.Config) *ContentConfigurationReconciler {
	subs := []lifecycle.Subroutine{}
	if cfg.Subroutines.ContentConfiguration.Enabled {
		subs = append(subs, subroutines.NewContentConfigurationSubroutine(validation.NewContentConfiguration(), http.DefaultClient))
	}
	return &ContentConfigurationReconciler{
		lifecycle: lifecycle.NewLifecycleManager(log, operatorName, contentConfigurationReconcilerName, mgr.GetClient(), subs).WithSpreadingReconciles().WithConditionManagement(),
	}
}

func (r *ContentConfigurationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	return r.lifecycle.Reconcile(ctx, req, &cachev1alpha1.ContentConfiguration{})
}

func (r *ContentConfigurationReconciler) SetupWithManager(mgr ctrl.Manager, cfg config.Config, log *logger.Logger, eventPredicates ...predicate.Predicate) error {
	return r.lifecycle.SetupWithManager(mgr, cfg.MaxConcurrentReconciles, contentConfigurationReconcilerName, &cachev1alpha1.ContentConfiguration{}, cfg.DebugLabelValue, r, log, eventPredicates...)
}
