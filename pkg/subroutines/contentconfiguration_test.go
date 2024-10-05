package subroutines

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/openmfp/golang-commons/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"

	golangCommonErrors "github.com/openmfp/golang-commons/errors"
	apimachinery "k8s.io/apimachinery/pkg/apis/meta/v1"

	cachev1alpha1 "github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/pkg/subroutines/mocks"
	"github.com/openmfp/extension-content-operator/pkg/validation"
	"github.com/openmfp/extension-content-operator/pkg/validation/validation_test"
)

type ContentConfigurationSubroutineTestSuite struct {
	suite.Suite

	testObj *ContentConfigurationSubroutine

	// mocks
	clientMock *mocks.Client
}

func TestContentConfigurationSubroutineTestSuit(t *testing.T) {
	suite.Run(t, new(ContentConfigurationSubroutineTestSuite))
}

func (suite *ContentConfigurationSubroutineTestSuite) SetupTest() {
	// create new mock
	suite.clientMock = new(mocks.Client)

	// create new test object
	suite.testObj = NewContentConfigurationSubroutine(validation.NewContentConfiguration(), http.DefaultClient)
}

func (suite *ContentConfigurationSubroutineTestSuite) TestCreateAndUpdate_OK() {
	// Given
	contentConfiguration := &cachev1alpha1.ContentConfiguration{
		Spec: cachev1alpha1.ContentConfigurationSpec{
			InlineConfiguration: cachev1alpha1.InlineConfiguration{
				Content:     validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
				ContentType: "yaml",
			},
		},
	}

	// When
	_, err := suite.testObj.Process(context.Background(), contentConfiguration)

	// Then
	suite.Require().Nil(err)

	cmp, cmpErr := compareYAML(
		validation_test.GetJSONFixture(validation_test.GetValidJSON()),
		contentConfiguration.Status.ConfigurationResult,
	)
	suite.Require().Nil(cmpErr)
	suite.Require().True(cmp)

	// Now lets take the same object and update it
	// Given
	contentConfiguration.Spec.InlineConfiguration.Content = validation_test.GetYAMLFixture(
		validation_test.GetValidYAMLFixtureButDifferentName())

	// When
	_, err2 := suite.testObj.Process(context.Background(), contentConfiguration)

	// Then
	suite.Require().Nil(err2)
	cmp, cmpErr = compareYAML(
		validation_test.GetJSONFixture(validation_test.GetValidJSONButDifferentName()),
		contentConfiguration.Status.ConfigurationResult,
	)
	suite.Require().Nil(cmpErr)
	suite.Require().True(cmp)
}

func (suite *ContentConfigurationSubroutineTestSuite) TestCreateAndUpdate_Error() {
	// Given
	contentConfiguration := &cachev1alpha1.ContentConfiguration{
		Spec: cachev1alpha1.ContentConfigurationSpec{
			InlineConfiguration: cachev1alpha1.InlineConfiguration{
				Content:     validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
				ContentType: "yaml",
			},
		},
	}

	// When
	_, errCmp := suite.testObj.Process(context.Background(), contentConfiguration)

	// Then
	suite.Require().Nil(errCmp)

	// compare configuration and result YAMLs
	cmp, cmpErr := compareYAML(
		validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
		contentConfiguration.Status.ConfigurationResult,
	)
	suite.Require().Nil(cmpErr)
	suite.Require().True(cmp)

	// Given invalid configuration
	contentConfiguration.Spec.InlineConfiguration.Content = "invalid"

	// When
	_, errProcessInvalidConfig := suite.testObj.Process(context.Background(), contentConfiguration)
	time.Sleep(1 * time.Second)

	// Then
	suite.Require().Nil(errProcessInvalidConfig)
	// result shoundn't change
	equal, cmpErr := compareYAML(
		validation_test.GetYAMLFixture(validation_test.GetValidYAML()), contentConfiguration.Status.ConfigurationResult)
	suite.Require().Nil(cmpErr)
	suite.Require().True(equal)
}

func (suite *ContentConfigurationSubroutineTestSuite) TestGetName_OK() {
	// When
	result := suite.testObj.GetName()

	// Then
	suite.Equal(ContentConfigurationSubroutineName, result)
}

func (suite *ContentConfigurationSubroutineTestSuite) TestFinalize_OK() {
	// Given
	contentConfiguration := &cachev1alpha1.ContentConfiguration{}

	// When
	result, err := suite.testObj.Finalize(context.Background(), contentConfiguration)

	// Then
	suite.False(result.Requeue)
	suite.Assert().Zero(result.RequeueAfter)
	suite.Nil(err)
}

func (suite *ContentConfigurationSubroutineTestSuite) TestProcessingConfig() {
	remoteURL := "https://this-address-should-be-mocked-by-httpmock"

	tests := []struct {
		name                 string
		spec                 cachev1alpha1.ContentConfigurationSpec
		remoteURL            string
		statusCode           int
		expectedError        golangCommonErrors.OperatorError
		expectedConfigResult string
	}{
		{
			name: "InlineConfigYAML_OK",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					Content:     validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
					ContentType: "yaml",
				},
			},
			expectedConfigResult: validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
		},
		{
			name: "InlineConfigYAML_ValidationError",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					Content:     "I am not a valid yaml",
					ContentType: "yaml",
				},
			},
		},
		{
			name: "InlineConfigJSON_OK",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					Content:     validation_test.GetJSONFixture(validation_test.GetValidJSON()),
					ContentType: "json",
				},
			},
			expectedConfigResult: validation_test.GetJSONFixture(validation_test.GetValidJSON()),
		},
		{
			name: "InlineConfigJSON_ValidationError",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					Content:     "I am not a valid json",
					ContentType: "json",
				},
			},
		},
		{
			name: "RemoteConfig_OK",
			spec: cachev1alpha1.ContentConfigurationSpec{
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					ContentType: "json",
					URL:         remoteURL,
				},
			},
			remoteURL:            remoteURL,
			statusCode:           http.StatusOK,
			expectedConfigResult: validation_test.GetValidJSON(),
		},
		{
			name: "RemoteConfig_http_error",
			spec: cachev1alpha1.ContentConfigurationSpec{
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					URL: remoteURL,
				},
			},
			remoteURL:     remoteURL,
			statusCode:    http.StatusInternalServerError,
			expectedError: golangCommonErrors.NewOperatorError(errors.New("received non-200 status code: 500"), false, true),
		},
		{
			name:          "NoConfigProvider_Error",
			spec:          cachev1alpha1.ContentConfigurationSpec{},
			expectedError: golangCommonErrors.NewOperatorError(errors.New("no configuration provided"), false, true),
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			if tt.remoteURL != "" {
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder(
					"GET", tt.remoteURL, httpmock.NewStringResponder(tt.statusCode, validation_test.GetValidJSON()),
				)
			}

			// When
			contentConfiguration := cachev1alpha1.ContentConfiguration{
				Spec: tt.spec,
			}
			_, err := suite.testObj.Process(context.Background(), &contentConfiguration)

			// Then
			if tt.expectedError != nil {
				if err == nil {
					suite.Fail("expected error, but got nil")
				}
				suite.Require().Equal(tt.expectedError.Err().Error(), err.Err().Error())
			} else {
				suite.Nil(err)
			}

			cmp, cmpErr := compareYAML(tt.expectedConfigResult, contentConfiguration.Status.ConfigurationResult)
			suite.Require().Nil(cmpErr)
			suite.Require().True(cmp)
		})
	}
}

func (suite *ContentConfigurationSubroutineTestSuite) TestFinalizers_OK() {
	// Given
	contentConfiguration := &cachev1alpha1.ContentConfiguration{}

	// When
	result, err := suite.testObj.Finalize(context.Background(), contentConfiguration)

	// Then
	suite.False(result.Requeue)
	suite.Assert().Zero(result.RequeueAfter)
	suite.Nil(err)

	// When
	finalizers := suite.testObj.Finalizers()

	// Then
	suite.Equal([]string{}, finalizers)

}

func TestService_Do(t *testing.T) {
	log, err := logger.New(logger.DefaultConfig())
	require.NoError(t, err)
	tests := []struct {
		name           string
		url            string
		mockResponse   string
		mockStatusCode int
		mockError      error
		expectedBody   string
		expectError    bool
	}{
		{
			name:           "GET_request_OK",
			url:            "https://example.com/success",
			mockResponse:   `{"message": "success"}`,
			mockStatusCode: http.StatusOK,
			expectedBody:   `{"message": "success"}`,
			expectError:    false,
		},
		{
			name:           "status_code_500_Error",
			url:            "https://example.com/error",
			mockResponse:   `{"message": "error"}`,
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:           "status_code_404_Error",
			url:            "https://example.com/error",
			mockResponse:   `{"message": "error"}`,
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:        "network_Error",
			url:         "https://example.com/network-error",
			mockError:   errors.New("network error"),
			expectError: true,
		},
		{
			name:        "invalidurl",
			url:         "://invalid-url",
			mockError:   errors.New("network error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			if tt.mockError != nil {
				httpmock.RegisterResponder(http.MethodGet, tt.url,
					httpmock.NewErrorResponder(tt.mockError))
			} else {
				httpmock.RegisterResponder(http.MethodGet, tt.url,
					httpmock.NewStringResponder(tt.mockStatusCode, tt.mockResponse))
			}

			r := NewContentConfigurationSubroutine(validation.NewContentConfiguration(), http.DefaultClient)

			body, err, _ := r.getRemoteConfig(tt.url, log)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, string(body))
			}
		})
	}
}

func compareYAML(doc1 string, doc2 string) (bool, error) {
	var obj1, obj2 map[string]interface{}
	errConfig := yaml.Unmarshal([]byte(doc1), &obj1)
	if errConfig != nil {
		return false, errConfig
	}
	errResult := yaml.Unmarshal([]byte(doc2), &obj2)
	if errResult != nil {
		return false, errResult
	}

	return reflect.DeepEqual(obj1, obj2), nil
}

func (suite *ContentConfigurationSubroutineTestSuite) Test_IncompatibleSchemaUpdate() {
	// Given
	contentConfiguration := &cachev1alpha1.ContentConfiguration{
		Spec: cachev1alpha1.ContentConfigurationSpec{
			InlineConfiguration: cachev1alpha1.InlineConfiguration{
				Content:     validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
				ContentType: "yaml",
			},
		},
		Status: cachev1alpha1.ContentConfigurationStatus{
			Conditions: []apimachinery.Condition{
				{
					Type:    "Ready",
					Status:  "True",
					Message: "The resource is ready",
					Reason:  "Complete",
				},
				{
					Message: "The subroutine is complete",
					Reason:  "Complete",
					Status:  "True",
					Type:    "ContentConfigurationSubroutine_Ready",
				},
			},
			ConfigurationResult: validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
		},
	}

	// Simulate incompatible schema update
	contentConfiguration.Spec.InlineConfiguration.Content = validation_test.GetYAMLFixture(
		validation_test.GetValidIncompatibleYAML(),
	)

	// When
	result, err := suite.testObj.Process(context.Background(), contentConfiguration)

	// Then: should keep previously valid and currently invalid result
	suite.Require().Nil(err)
	suite.Require().Empty(result)

	cmp, cmpErr := compareYAML(
		validation_test.GetJSONFixture(validation_test.GetValidJSON()),
		contentConfiguration.Status.ConfigurationResult,
	)
	suite.Require().Nil(cmpErr)
	suite.Require().True(cmp)
	suite.Require().True(
		getCondition(contentConfiguration.Status.Conditions, ValidationConditionType).Status == apimachinery.ConditionFalse,
	)
	suite.Require().Equal(
		"ValidationFailed", getCondition(contentConfiguration.Status.Conditions, ValidationConditionType).Reason,
	)

	// make it valid and check if condition is removed
	contentConfiguration.Spec.InlineConfiguration.Content = validation_test.GetYAMLFixture(validation_test.GetValidYAML())

	// When
	result, err = suite.testObj.Process(context.Background(), contentConfiguration)
	suite.Require().Nil(err)
	suite.Require().Empty(result)

	cmp, cmpErr = compareYAML(
		validation_test.GetYAMLFixture(validation_test.GetValidYAML()),
		contentConfiguration.Status.ConfigurationResult,
	)
	suite.Require().NoError(cmpErr)
	suite.Require().True(cmp)

	suite.Require().Equal(
		"ValidationSucceeded", getCondition(contentConfiguration.Status.Conditions, ValidationConditionType).Reason,
	)
}

func getCondition(conditions []apimachinery.Condition, conditionType string) apimachinery.Condition { // nolint: unparam
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return condition
		}
	}
	return apimachinery.Condition{}
}
