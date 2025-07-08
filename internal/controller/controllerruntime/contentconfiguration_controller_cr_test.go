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
	"fmt"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/openmfp/extension-manager-operator/api/v1alpha1"
	commonTesting "github.com/openmfp/extension-manager-operator/pkg/util/testing"
	"github.com/openmfp/extension-manager-operator/pkg/validation/validation_test"
)

func TestContentConfigurationCRTestSuite(t *testing.T) {
	suite.Run(t, new(ContentConfigurationControllerTestSuite))
}

func (suite *ContentConfigurationControllerTestSuite) TestContentConfigurationCRCreation() {
	remoteURL := "https://this-address-should-be-mocked-by-httpmock"

	// Define the test cases
	testCases := []struct {
		name           string
		instanceName   string
		spec           v1alpha1.ContentConfigurationSpec
		expectedResult string
	}{
		{
			name:         "TestInlineContentConfiguration",
			instanceName: "inline",
			spec: v1alpha1.ContentConfigurationSpec{
				InlineConfiguration: &v1alpha1.InlineConfiguration{
					ContentType: "yaml",
					Content:     validation_test.GetValidYAML(),
				},
			},
			expectedResult: validation_test.GetValidJSON(),
		},
		{
			name:         "TestBothInlineAndRemoteConfiguration",
			instanceName: "inline-and-remote",
			spec: v1alpha1.ContentConfigurationSpec{
				InlineConfiguration: &v1alpha1.InlineConfiguration{
					ContentType: "yaml",
					Content:     validation_test.GetValidYAML(),
				},
				RemoteConfiguration: &v1alpha1.RemoteConfiguration{
					URL:         "this-url-should-not-be-used",
					ContentType: "yaml",
				},
			},
			expectedResult: validation_test.GetValidJSON(),
		},
		{
			name:         "TestRemoteContentConfiguration",
			instanceName: "remote",
			spec: v1alpha1.ContentConfigurationSpec{
				RemoteConfiguration: &v1alpha1.RemoteConfiguration{
					ContentType: "json",
					URL:         remoteURL,
				},
			},
			expectedResult: validation_test.GetValidJSON(),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(
				"GET", remoteURL, httpmock.NewStringResponder(200, validation_test.GetValidJSON()),
			)

			// Given
			testContext := context.Background()
			instance := &v1alpha1.ContentConfiguration{
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
			createdInstance := v1alpha1.ContentConfiguration{}
			suite.Assert().Eventually(
				func() bool {
					err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
						Name:      tc.instanceName,
						Namespace: defaultNamespace,
					}, &createdInstance)

					equal, cErr := commonTesting.CompareJSON(tc.expectedResult, createdInstance.Status.ConfigurationResult)
					return err == nil && cErr == nil && equal
				},
				defaultTestTimeout, defaultTickInterval,
			)
		})
	}
}

func (suite *ContentConfigurationControllerTestSuite) TestUpdateReconcileCR() {
	remoteURL := "https://this-address-should-be-mocked-by-httpmock"

	// Given
	testContext := context.Background()
	contentConfiguration := &v1alpha1.ContentConfiguration{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "extension-manager",
			Namespace: defaultNamespace,
		},
		Spec: v1alpha1.ContentConfigurationSpec{
			RemoteConfiguration: &v1alpha1.RemoteConfiguration{
				ContentType: "json",
				URL:         remoteURL,
			},
		},
	}

	// setup mocks
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET", remoteURL, httpmock.NewStringResponder(200, validation_test.GetValidJSON()),
	)

	// When
	err := suite.kubernetesClient.Create(testContext, contentConfiguration)
	suite.Nil(err)

	// Then
	createdInstance := v1alpha1.ContentConfiguration{}
	suite.Assert().Eventually(
		func() bool {
			err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
				Name:      contentConfiguration.Name,
				Namespace: contentConfiguration.Namespace,
			}, &createdInstance)

			equal, cErr := commonTesting.CompareJSON(validation_test.GetValidJSON(), createdInstance.Status.ConfigurationResult)
			return err == nil && cErr == nil && equal
		},
		defaultTestTimeout, defaultTickInterval,
	)

	// Update ContentConfiguration and check for 2nd reconcile
	// Given
	remoteURL = "https://new.url"
	createdInstance.Spec.RemoteConfiguration.URL = remoteURL
	httpmock.RegisterResponder(
		"GET", remoteURL, httpmock.NewStringResponder(200, validation_test.GetValidJSONButDifferentName()),
	)

	// When
	err = suite.kubernetesClient.Update(testContext, &createdInstance)
	suite.Nil(err)

	// Then
	updatedInstance := v1alpha1.ContentConfiguration{}
	suite.Assert().Eventually(
		func() bool {
			err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
				Name:      contentConfiguration.Name,
				Namespace: contentConfiguration.Namespace,
			}, &updatedInstance)

			equal, cerr := commonTesting.CompareJSON(validation_test.GetValidJSONButDifferentName(), updatedInstance.Status.ConfigurationResult)
			return err == nil && cerr == nil && equal
		},
		defaultTestTimeout, defaultTickInterval,
	)

	// 3rd reconcile: the same URL but it returns a different content; changed labels
	// Given
	remoteURL = "https://new.url2"
	updatedInstance.Spec.RemoteConfiguration.URL = remoteURL
	httpmock.Reset()
	httpmock.RegisterResponder(
		"GET", remoteURL, httpmock.NewStringResponder(200, validation_test.GetValidJSON()),
	)

	// When
	log.Info().Msg(fmt.Sprintf("Before Update ObservedGeneration: %d", updatedInstance.Status.ObservedGeneration))
	err = suite.kubernetesClient.Update(testContext, &updatedInstance)
	suite.NoError(err)
	suite.logger.Info().Msg("--------------------- resource updated a 2nd time ---------------------")
	log.Info().Msg(fmt.Sprintf("After Update ObservedGeneration: %d", updatedInstance.Status.ObservedGeneration))
	// Then
	updatedInstanceSameUrl := v1alpha1.ContentConfiguration{}
	suite.Assert().Eventually(
		func() bool {
			err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
				Name:      contentConfiguration.Name,
				Namespace: contentConfiguration.Namespace,
			}, &updatedInstanceSameUrl)

			log.Info().Msg(fmt.Sprintf("ObservedGeneration: %d", updatedInstance.Status.ObservedGeneration))
			result := err == nil && updatedInstanceSameUrl.Status.ObservedGeneration == updatedInstanceSameUrl.Generation
			return result
		},
		defaultTestTimeout, defaultTickInterval,
	)
	equal, err := commonTesting.CompareJSON(validation_test.GetValidJSON(), updatedInstanceSameUrl.Status.ConfigurationResult)
	suite.NoError(err)
	suite.True(equal)
}

func (suite *ContentConfigurationControllerTestSuite) TestContentConfigurationCreationCRInternalURL() {
	remoteURL := "https://this-address-should-be-mocked-by-httpmock"
	internalURL := "http://internal-url"

	// Given
	testContext := context.Background()
	contentConfiguration := &v1alpha1.ContentConfiguration{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "extension-manager-internal",
			Namespace: defaultNamespace,
		},
		Spec: v1alpha1.ContentConfigurationSpec{
			RemoteConfiguration: &v1alpha1.RemoteConfiguration{
				ContentType: "json",
				InternalUrl: internalURL,
				URL:         remoteURL,
			},
		},
	}

	// setup mocks
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET", remoteURL, httpmock.NewStringResponder(200, validation_test.GetValidJSON()),
	)

	// When
	err := suite.kubernetesClient.Create(testContext, contentConfiguration)
	suite.Nil(err)

	// Then
	createdInstance := v1alpha1.ContentConfiguration{}
	suite.Assert().Eventually(
		func() bool {
			err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
				Name:      contentConfiguration.Name,
				Namespace: contentConfiguration.Namespace,
			}, &createdInstance)
			return err == nil && createdInstance.Status.ConfigurationResult == ""
		},
		defaultTestTimeout, defaultTickInterval,
	)

	// Update InternalURL and make it valid by mocking
	// Given
	httpmock.RegisterResponder(
		"GET", internalURL, httpmock.NewStringResponder(200, validation_test.GetValidJSON()),
	)

	// When
	err = suite.kubernetesClient.Update(testContext, &createdInstance)
	suite.Nil(err)

	time.Sleep(1 * time.Second)

	// Then
	updatedInstance := v1alpha1.ContentConfiguration{}
	suite.Assert().Eventually(
		func() bool {
			err := suite.kubernetesClient.Get(testContext, types.NamespacedName{
				Name:      contentConfiguration.Name,
				Namespace: contentConfiguration.Namespace,
			}, &updatedInstance)

			equal, cErr := commonTesting.CompareJSON(validation_test.GetValidJSON(), updatedInstance.Status.ConfigurationResult)
			return err == nil && cErr == nil && equal
		},
		defaultTestTimeout, defaultTickInterval,
	)

}
