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
	"net/http"
	"time"

	openmfpcontext "github.com/openmfp/golang-commons/context"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/openmfp/extension-content-operator/internal/config"
	"github.com/openmfp/extension-content-operator/internal/server"
	"github.com/openmfp/extension-content-operator/pkg/validation"
)

var (
	appCfg config.Config
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server with configuration validation endpoint",
	Run:   RunServer,
}

func init() { // coverage-ignore
	var err error
	appCfg, err = config.NewFromEnv()
	if err != nil {
		setupLog.Error(err, "unable to load config")
		panic(err)
	}
}

func RunServer(cmd *cobra.Command, args []string) { // coverage-ignore
	log := initLog()
	ctrl.SetLogger(log.ComponentLogger("server").Logr())

	ctx, _, shutdown := openmfpcontext.StartContext(log, nil, appCfg.ShutdownTimeout)
	defer shutdown()

	rt := server.CreateRouter(appCfg, log, validation.NewContentConfiguration())

	server := &http.Server{
		Addr:         ":" + appCfg.ServerPort,
		Handler:      rt,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		server.ListenAndServe()
	}()
	log.Info().Msg("Server started")

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Panic().Err(err).Msg("Graceful shutdown failed")
	}
	log.Info().Msg("Server stopped")
}
