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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	cachev1alpha1 "github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/internal/config"
	openmfpcontext "github.com/openmfp/golang-commons/context"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/logger"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

const (
	defaultTestTimeout  = 10 * time.Second
	defaultTickInterval = 250 * time.Millisecond
	defaultNamespace    = "default"
)

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "chart", "crds")},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.28.3-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = cachev1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("ContentConfiguration Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: defaultNamespace,
		}
		contentconfiguration := &cachev1alpha1.ContentConfiguration{}

		logConfig := logger.DefaultConfig()
		logConfig.NoJSON = true
		logConfig.Name = "ContentConfigurationTestSuite"
		log, _ := logger.New(logConfig)
		// Disable color logging as vs-code does not support color logging in the test output
		log = logger.NewFromZerolog(log.Output(&zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}))

		BeforeEach(func() {
			By("creating the custom resource for the Kind ContentConfiguration")
			err := k8sClient.Get(ctx, typeNamespacedName, contentconfiguration)
			if err != nil && errors.IsNotFound(err) {
				resource := &cachev1alpha1.ContentConfiguration{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: defaultNamespace,
					},
					// TODO(user): Specify other spec details if needed.
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &cachev1alpha1.ContentConfiguration{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance ContentConfiguration")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ContentConfigurationReconciler{
				lifecycle: lifecycle.NewLifecycleManager(log, operatorName, contentConfigurationReconcilerName, k8sClient, []lifecycle.Subroutine{}).WithSpreadingReconciles().WithConditionManagement(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})

type ContentConfigurationTestSuite struct {
	suite.Suite

	kubernetesClient  client.Client
	kubernetesManager ctrl.Manager
	testEnv           *envtest.Environment

	cancel context.CancelFunc
}

// TestControllers takes part in it, do not delete it
func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
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
