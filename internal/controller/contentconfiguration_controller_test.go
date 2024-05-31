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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	cachev1alpha1 "github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/internal/config"
	"github.com/openmfp/golang-commons/controller/lifecycle"
	"github.com/openmfp/golang-commons/logger"
)

var _ = Describe("ContentConfiguration Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
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
						Namespace: "default",
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
		It("should successfully call NewContentConfigurationReconciler", func() {
			By("Create reconciler with NewContentConfigurationReconciler")

			k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{})
			cfg := config.Config{}
			reconciler := NewContentConfigurationReconciler(log, k8sManager, cfg)

			Expect(err).NotTo(HaveOccurred())
			Expect(reconciler).NotTo(BeNil())

			err = reconciler.SetupWithManager(k8sManager, cfg, log)
			Expect(err).NotTo(HaveOccurred())

			result, _ := reconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			// Expect(errReconcile).To(BeNil())
			Expect(result).NotTo(BeNil())
		})
	})
})
