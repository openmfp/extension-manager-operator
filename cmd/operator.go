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

package cmd

import (
	"context"
	"crypto/tls"

	openmfpcontext "github.com/openmfp/golang-commons/context"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/openmfp/extension-manager-operator/internal/controller"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "operator to reconcile ContentConfiguration",
	Run:   RunController,
}

func RunController(cmd *cobra.Command, args []string) { // coverage-ignore
	ctrl.SetLogger(log.ComponentLogger("controller-runtime").Logr())

	ctx, _, shutdown := openmfpcontext.StartContext(log, appConfig, appConfig.ShutdownTimeout)
	defer shutdown()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: appConfig.MetricsBindAddress,
			TLSOpts: []func(*tls.Config){
				func(c *tls.Config) {
					log.Info().Msg("disabling http/2")
					c.NextProtos = []string{"http/1.1"}
				},
			},
		},
		BaseContext:                   func() context.Context { return ctx },
		HealthProbeBindAddress:        appConfig.HealthProbeBindAddress,
		LeaderElection:                appConfig.LeaderElection.Enabled,
		LeaderElectionID:              "eengiex4.openmfp.org",
		LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start manager")
	}

	contentConfigurationReconciler := controller.NewContentConfigurationReconciler(log, mgr, appConfig)
	if err := contentConfigurationReconciler.SetupWithManager(mgr, appConfig, log); err != nil {
		log.Fatal().Err(err).Str("controller", "ContentConfiguration").Msg("unable to create controller")
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Fatal().Err(err).Msg("unable to set up health check")
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Fatal().Err(err).Msg("unable to set up ready check")
	}

	log.Info().Msg("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal().Err(err).Msg("problem running manager")
	}
}
