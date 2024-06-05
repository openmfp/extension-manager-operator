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
	"bytes"
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	cachev1alpha1 "github.com/openmfp/extension-content-operator/api/v1alpha1"
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
					Content:     getValidYAMLFixture(),
				},
			},
			expectedResult: getValidJSONFixture(),
		},
		{
			name:         "TestBothInlineAndRemoteConfiguration",
			instanceName: "inline-and-remote",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					ContentType: "yaml",
					Content:     getValidYAMLFixture(),
				},
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					URL: "this-url-should-not-be-used",
				},
			},
			expectedResult: getValidJSONFixture(),
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
			expectedResult: getValidJSONFixture(),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(
				"GET", remoteURL, httpmock.NewStringResponder(200, getValidJSONFixture()),
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

func getValidYAMLFixture() string {
	validYAML := `
name: overview
luigiConfigFragment:
- data:
    nodes:
    - entityType: global
      pathSegment: home
      label: Overview
      icon: home
      hideFromNav: true
      defineEntity:
        id: example
      children:
      - pathSegment: overview
        label: Overview
        icon: home
        url: https://fiddle.luigi-project.io/examples/microfrontends/multipurpose.html
        context:
          title: Welcome to OpenMFP Portal
          content: " "
`

	var data interface{}
	err := yaml.Unmarshal([]byte(validYAML), &data)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	compactYAML, err := yaml.Marshal(&data)
	if err != nil {
		log.Fatalf("failed to marshal YAML: %v", err)
	}

	return string(compactYAML)
}

func getValidJSONFixture() string {
	validJSON := `{
		"name": "overview",
		"luigiConfigFragment": [
			{
				"data": {
					"nodes": [
						{
							"entityType": "global",
							"pathSegment": "home",
							"label": "Overview",
							"icon": "home"
						}
					]
				}
			}
		]
	}`

	var buf bytes.Buffer
	if err := json.Compact(&buf, []byte(validJSON)); err != nil {
		return ""
	}

	return buf.String()
}
