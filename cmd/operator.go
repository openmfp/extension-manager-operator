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
	"os"

	openmfpcontext "github.com/openmfp/golang-commons/context"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/openmfp/extension-manager-operator/internal/config"
	"github.com/openmfp/extension-manager-operator/internal/controller"
)

var operatorCmd = &cobra.Command{
	Use:   "operator",
	Short: "operator to reconcile ContentConfiguration",
	Run:   RunController,
}

var (
	metricsAddr          string
	enableLeaderElection bool
	probeAddr            string
	loglevel             string
	logNoJSON            bool
	secureMetrics        bool
	enableHTTP2          bool
	cfg                  config.Config
)

func init() { // coverage-ignore
	var err error
	cfg, err = config.NewFromEnv()
	if err != nil {
		setupLog.Error(err, "unable to load config")
		os.Exit(1)
	}
	operatorCmd.Flags().StringVar(&metricsAddr, "metrics-bind-address", cfg.Metrics.BindAddress,
		"The address the metric endpoint binds to.")
	operatorCmd.Flags().StringVar(&probeAddr, "health-probe-bind-address", cfg.Probes.BindAddress,
		"The address the probe endpoint binds to.")
	operatorCmd.Flags().StringVar(&loglevel, "log-level", cfg.Log.Level,
		"The log level for the application. Default is info.")
	operatorCmd.Flags().BoolVar(&logNoJSON, "log-no-json", cfg.Log.NoJSON,
		"Flag to disable JSON logging. Default is false.")
	operatorCmd.Flags().BoolVar(&enableLeaderElection, "leader-elect", cfg.LeaderElection.Enabled,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	operatorCmd.Flags().BoolVar(&secureMetrics, "metrics-secure", cfg.Metrics.Secure,
		"If set the metrics endpoint is served securely")
	operatorCmd.Flags().BoolVar(&enableHTTP2, "enable-http2", cfg.EnableHTTP2,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")
}

func RunController(cmd *cobra.Command, args []string) { // coverage-ignore
	_, log := initApp()
	ctrl.SetLogger(log.ComponentLogger("controller-runtime").Logr())

	ctx, _, shutdown := openmfpcontext.StartContext(log, cfg, cfg.ShutdownTimeout)
	defer shutdown()

	disableHTTP2 := func(c *tls.Config) {
		log.Info().Msg("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	tlsOpts := []func(*tls.Config){}
	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: tlsOpts,
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		Metrics: metricsserver.Options{
			BindAddress:   metricsAddr,
			SecureServing: secureMetrics,
			TLSOpts:       tlsOpts,
		},
		BaseContext:                   func() context.Context { return ctx },
		WebhookServer:                 webhookServer,
		HealthProbeBindAddress:        probeAddr,
		LeaderElection:                enableLeaderElection,
		LeaderElectionID:              "eengiex4.openmfp.io",
		LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("unable to start manager")
	}

	contentConfigurationReconciler := controller.NewContentConfigurationReconciler(log, mgr, cfg)
	if err := contentConfigurationReconciler.SetupWithManager(mgr, cfg, log); err != nil {
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
