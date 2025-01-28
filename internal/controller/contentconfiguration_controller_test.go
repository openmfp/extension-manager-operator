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
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	cachev1alpha1 "github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/pkg/validation/validation_test"
)

func TestContentConfigurationTestSuite(t *testing.T) {
	suite.Run(t, new(ContentConfigurationTestSuite))
}

func (suite *ContentConfigurationTestSuite) TestContentConfigurationCreation() {
	remoteURL := "https://this-address-should-be-mocked-by-httpmock"

	// Define the test cases
	testCases := []struct {
		name           string
		instanceName   string
		spec           cachev1alpha1.ContentConfigurationSpec
		expectedResult string
	}{
		{
			name:         "TestInlineContentConfiguration",
			instanceName: "inline",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					ContentType: "yaml",
					Content:     validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
				},
			},
			expectedResult: validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
		},
		{
			name:         "TestBothInlineAndRemoteConfiguration",
			instanceName: "inline-and-remote",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					ContentType: "yaml",
					Content:     validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
				},
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					URL: "this-url-should-not-be-used",
				},
			},
			expectedResult: validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
		},
		{
			name:         "TestRemoteContentConfiguration",
			instanceName: "remote",
			spec: cachev1alpha1.ContentConfigurationSpec{
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					ContentType: "json",
					URL:         remoteURL,
				},
			},
			expectedResult: validation_test.GetJSONFixture(validation_test.GetValidJSON()),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(
				"GET", remoteURL, httpmock.NewStringResponder(200, validation_test.GetJSONFixture(validation_test.GetValidJSON())),
			)

			// Given
			testContext := context.Background()
			instance := &cachev1alpha1.ContentConfiguration{
				ObjectMeta: metaV1.ObjectMeta{
					Name:      tc.instanceName,
					Namespace: defaultNamespace,
				},
				Spec: tc.spec,
			}

			// When
			err := suite.kubernetesClient.Create(testContext, instance)
			suite.Nil(err)

			// Then
			createdInstance := cachev1alpha1.ContentConfiguration{}
			suite.Assert().Eventually(
				func() bool {
					err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
						Name:      tc.instanceName,
						Namespace: defaultNamespace,
					}, &createdInstance)
					return err == nil && createdInstance.Status.ConfigurationResult == tc.expectedResult
				},
				defaultTestTimeout, defaultTickInterval,
			)
		})
	}
}

func (suite *ContentConfigurationTestSuite) TestUpdateReconcile() {
	remoteURL := "https://this-address-should-be-mocked-by-httpmock"

	suite.Run("TearDown", func() {

		// Given
		testContext := context.Background()
		contentConfiguration := &cachev1alpha1.ContentConfiguration{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      "extension-manager",
				Namespace: defaultNamespace,
			},
			Spec: cachev1alpha1.ContentConfigurationSpec{
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					ContentType: "json",
					URL:         remoteURL,
				},
			},
		}

		// setup mocks
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder(
			"GET", remoteURL, httpmock.NewStringResponder(200, validation_test.GetJSONFixture(validation_test.GetValidJSON())),
		)

		// When
		err := suite.kubernetesClient.Create(testContext, contentConfiguration)
		suite.Nil(err)

		// Then
		createdInstance := cachev1alpha1.ContentConfiguration{}
		suite.Assert().Eventually(
			func() bool {
				err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
					Name:      contentConfiguration.Name,
					Namespace: contentConfiguration.Namespace,
				}, &createdInstance)
				return err == nil && createdInstance.Status.ConfigurationResult == validation_test.GetJSONFixture(validation_test.GetValidJSON())
			},
			defaultTestTimeout, defaultTickInterval,
		)

		// Update ContentConfiguration and check for 2nd reconcile
		// Given
		createdInstance.Spec.RemoteConfiguration.URL = "https://new.url"
		httpmock.RegisterResponder(
			"GET", "https://new.url", httpmock.NewStringResponder(200, validation_test.GetJSONFixture(validation_test.GetValidJSONButDifferentName())),
		)

		// When
		err = suite.kubernetesClient.Update(testContext, &createdInstance)
		suite.Nil(err)

		// Then
		updatedInstance := cachev1alpha1.ContentConfiguration{}
		suite.Assert().Eventually(
			func() bool {
				err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
					Name:      contentConfiguration.Name,
					Namespace: contentConfiguration.Namespace,
				}, &updatedInstance)

				return err == nil &&
					updatedInstance.Status.ConfigurationResult == validation_test.GetJSONFixture(validation_test.GetValidJSONButDifferentName())
			},
			defaultTestTimeout, defaultTickInterval,
		)

		// 3rd reconcile: the same URL but it returns a different content; changed labels
		// Given
		updatedInstance.Spec.RemoteConfiguration.URL = "https://new.url"
		httpmock.Reset()
		httpmock.RegisterResponder(
			"GET", "https://new.url", httpmock.NewStringResponder(200, validation_test.GetJSONFixture(validation_test.GetValidJSON())),
		)

		updatedInstance.ObjectMeta.Labels = map[string]string{
			"somelabel": "somevalue",
		}

		// When
		err = suite.kubernetesClient.Update(testContext, &updatedInstance)
		suite.Nil(err)
		println("--------------------- resource updated a 2nd time ---------------------")

		// Then
		updatedInstanceSameUrl := cachev1alpha1.ContentConfiguration{}
		suite.Assert().Eventually(
			func() bool {
				err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
					Name:      contentConfiguration.Name,
					Namespace: contentConfiguration.Namespace,
				}, &updatedInstanceSameUrl)

				return err == nil &&
					updatedInstanceSameUrl.Status.ConfigurationResult == validation_test.GetJSONFixture(validation_test.GetValidJSON())
			},
			defaultTestTimeout, defaultTickInterval,
		)
		suite.Assert().True(updatedInstanceSameUrl.ObjectMeta.Labels["somelabel"] == "somevalue")

	})
}
