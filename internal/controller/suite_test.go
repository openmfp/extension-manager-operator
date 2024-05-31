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
	"os"
	"path/filepath"
	"time"

	//+kubebuilder:scaffold:imports
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	cachev1alpha1 "github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/internal/config"
	openmfpcontext "github.com/openmfp/golang-commons/context"
	"github.com/openmfp/golang-commons/logger"
)

const (
	defaultTestTimeout  = 10 * time.Second
	defaultTickInterval = 250 * time.Millisecond
	defaultNamespace    = "default"
)

type ContentConfigurationTestSuite struct {
	suite.Suite

	kubernetesClient  client.Client
	kubernetesManager ctrl.Manager
	testEnv           *envtest.Environment

	cancel context.CancelFunc
}

func (suite *ContentConfigurationTestSuite) SetupSuite() {
	logConfig := logger.DefaultConfig()
	logConfig.NoJSON = true
	logConfig.Name = "ContentConfigurationTestSuite"
	log, err := logger.New(logConfig)
	suite.Nil(err)
	// Disable color logging as vs-code does not support color logging in the test output
	log = logger.NewFromZerolog(log.Output(&zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}))

	cfg, err := config.NewFromEnv()
	suite.Nil(err)

	testContext, _, _ := openmfpcontext.StartContext(log, cfg, cfg.ShutdownTimeout)

	testContext = logger.SetLoggerInContext(testContext, log.ComponentLogger("TestSuite"))

	suite.testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "chart", "crds")},
		ErrorIfCRDPathMissing: true,
	}

	k8scfg, err := suite.testEnv.Start()
	suite.Nil(err)

	utilruntime.Must(cachev1alpha1.AddToScheme(scheme.Scheme))
	utilruntime.Must(v1.AddToScheme(scheme.Scheme))

	// +kubebuilder:scaffold:scheme

	suite.kubernetesClient, err = client.New(k8scfg, client.Options{
		Scheme: scheme.Scheme,
	})
	suite.Nil(err)
	ctrl.SetLogger(log.Logr())
	suite.kubernetesManager, err = ctrl.NewManager(k8scfg, ctrl.Options{
		Scheme:      scheme.Scheme,
		BaseContext: func() context.Context { return testContext },
	})
	suite.Nil(err)

	contentConfigurationReconciler := NewContentConfigurationReconciler(log, suite.kubernetesManager, cfg)
	err = contentConfigurationReconciler.SetupWithManager(suite.kubernetesManager, cfg, log)
	suite.Nil(err)

	go suite.startController()
}

func (suite *ContentConfigurationTestSuite) startController() {
	var controllerContext context.Context
	controllerContext, suite.cancel = context.WithCancel(context.Background())
	err := suite.kubernetesManager.Start(controllerContext)
	suite.Nil(err)
}

func (suite *ContentConfigurationTestSuite) TearDownSuite() {
	suite.cancel()
	err := suite.testEnv.Stop()
	suite.Nil(err)
}
