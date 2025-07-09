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
	"net/http"

	"github.com/kcp-dev/multicluster-provider/apiexport"
	openmfpcontext "github.com/platform-mesh/golang-commons/context"
	"github.com/platform-mesh/golang-commons/logger"
	"github.com/platform-mesh/golang-commons/traces"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	mcmanager "sigs.k8s.io/multicluster-runtime/pkg/manager"

	"github.com/openmfp/extension-manager-operator/internal/config"
	"github.com/openmfp/extension-manager-operator/internal/controller/controllerruntime"
	"github.com/openmfp/extension-manager-operator/internal/controller/multiclusterruntime"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "operator to reconcile ContentConfiguration",
	Run:   RunController,
}

func RunController(_ *cobra.Command, _ []string) { // coverage-ignore
	log.Info().Msg("Starting operator")
	ctrl.SetLogger(log.ComponentLogger("controller-runtime").Logr())

	ctx, _, shutdown := openmfpcontext.StartContext(log, operatorCfg, defaultCfg.ShutdownTimeout)
	defer shutdown()

	var err error
	var providerShutdown func(ctx context.Context) error
	if defaultCfg.Tracing.Enabled {
		providerShutdown, err = traces.InitProvider(ctx, defaultCfg.Tracing.Collector)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to start gRPC-Sidecar TracerProvider")
		}
	} else {
		providerShutdown, err = traces.InitLocalProvider(ctx, defaultCfg.Tracing.Collector, false)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to start local TracerProvider")
		}
	}

	defer func() {
		if err := providerShutdown(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to shutdown TracerProvider")
		}
	}()

	restCfg := ctrl.GetConfigOrDie()
	restCfg.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return otelhttp.NewTransport(rt)
	})

	if operatorCfg.KCP.Enabled {
		initializeMultiClusterManager(ctx, restCfg, log, operatorCfg)
	}

	initializeControllerRuntimeManager(ctx, restCfg)
}

func initializeMultiClusterManager(ctx context.Context, cfg *rest.Config, log *logger.Logger, operatorCfg config.OperatorConfig) {
	log.Info().Msg("Initializing multicluster manager")
	kubeconfigPath := operatorCfg.KCP.Kubeconfig
	kcpCfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		log.Fatal().Err(err).Str("controller", "ContentConfiguration").Msg("unable to construct cluster provider")
	}
	kcpCfg.Wrap(func(rt http.RoundTripper) http.RoundTripper {
		return otelhttp.NewTransport(rt)
	})

	provider, err := apiexport.New(kcpCfg, apiexport.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to construct cluster provider")
	}

	mgr, err := mcmanager.New(kcpCfg, provider, manager.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: defaultCfg.Metrics.BindAddress,
			TLSOpts: []func(*tls.Config){
				func(c *tls.Config) {
					log.Info().Msg("disabling http/2")
					c.NextProtos = []string{"http/1.1"}
				},
			},
		},
		BaseContext:                   func() context.Context { return ctx },
		HealthProbeBindAddress:        defaultCfg.HealthProbeBindAddress,
		LeaderElection:                defaultCfg.LeaderElection.Enabled,
		LeaderElectionID:              "eengiex3.openmfp.org",
		LeaderElectionReleaseOnCancel: true,
		LeaderElectionConfig:          cfg,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to set up overall controller manager")
	}

	contentConfigurationReconciler := multiclusterruntime.NewContentConfigurationReconciler(log, mgr, operatorCfg)
	if err := contentConfigurationReconciler.SetupWithManager(mgr, defaultCfg, log); err != nil {
		log.Fatal().Err(err).Str("controller", "ContentConfiguration").Msg("unable to create controller")
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Fatal().Err(err).Msg("unable to set up health check")
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		log.Fatal().Err(err).Msg("unable to set up ready check")
	}

	log.Info().Msg("Starting provider")
	go func() {
		if err := provider.Run(ctx, mgr); err != nil {
			log.Fatal().Err(err).Msg("unable to run provider")
		}
	}()
	log.Info().Msg("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Fatal().Err(err).Msg("problem running manager")
	}
}

func initializeControllerRuntimeManager(ctx context.Context, restCfg *rest.Config) {
	mgr, err := ctrl.NewManager(restCfg, ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress: defaultCfg.Metrics.BindAddress,
			TLSOpts: []func(*tls.Config){
				func(c *tls.Config) {
					log.Info().Msg("disabling http/2")
					c.NextProtos = []string{"http/1.1"}
				},
			},
		},
		BaseContext:                   func() context.Context { return ctx },
		HealthProbeBindAddress:        defaultCfg.HealthProbeBindAddress,
		LeaderElection:                defaultCfg.LeaderElection.Enabled,
		LeaderElectionID:              "eengiex4.openmfp.org",
		LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start manager")
	}

	contentConfigurationReconciler := controllerruntime.NewContentConfigurationReconcilerCR(log, mgr, operatorCfg)
	if err := contentConfigurationReconciler.SetupWithManager(mgr, defaultCfg, log); err != nil {
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
