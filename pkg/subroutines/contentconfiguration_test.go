package subroutines

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"

	cachev1alpha1 "github.com/openmfp/extension-content-operator/api/v1alpha1"
	"github.com/openmfp/extension-content-operator/pkg/subroutines/mocks"
	golangCommonErrors "github.com/openmfp/golang-commons/errors"
)

type ContentConfigurationSubroutineTestSuite struct {
	suite.Suite

	testObj *ContentConfigurationSubroutine

	// mocks
	clientMock *mocks.Client
}

func (suite *ContentConfigurationSubroutineTestSuite) SetupTest() {
	// create new mock
	suite.clientMock = new(mocks.Client)

	// create new test object
	suite.testObj = NewContentConfigurationSubroutine(nil)
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
	remoteUrl := "https://this-address-should-be-mocked-by-httpmock"
	inlineContent := "inline content"
	remoteContent := "remote content"

	tests := []struct {
		name                 string
		spec                 cachev1alpha1.ContentConfigurationSpec
		remoteURL            string
		statusCode           int
		expectedError        golangCommonErrors.OperatorError
		expectedConfigResult string
	}{
		{
			name: "InlineConfig_OK",
			spec: cachev1alpha1.ContentConfigurationSpec{
				InlineConfiguration: cachev1alpha1.InlineConfiguration{
					Content: inlineContent,
				},
			},
			expectedConfigResult: inlineContent,
		},
		{
			name: "RemoteConfig_OK",
			spec: cachev1alpha1.ContentConfigurationSpec{
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					URL: remoteUrl,
				},
			},
			remoteURL:            remoteUrl,
			statusCode:           http.StatusOK,
			expectedConfigResult: remoteContent,
		},
		{
			name: "RemoteConfig_http_error",
			spec: cachev1alpha1.ContentConfigurationSpec{
				RemoteConfiguration: cachev1alpha1.RemoteConfiguration{
					URL: remoteUrl,
				},
			},
			remoteURL:     remoteUrl,
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
					"GET", tt.remoteURL, httpmock.NewStringResponder(tt.statusCode, remoteContent),
				)
			}

			// When
			contentConfiguration := cachev1alpha1.ContentConfiguration{
				Spec: tt.spec,
			}
			_, err := suite.testObj.Process(context.Background(), &contentConfiguration)

			// Then
			if tt.expectedError != nil {
				suite.Require().Equal(tt.expectedError.Err().Error(), err.Err().Error())
			} else {
				suite.Nil(err)
			}

			suite.Require().Equal(tt.expectedConfigResult, contentConfiguration.Status.ConfigurationResult)
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
	suite.Equal([]string{"contentconfiguration.core.openmfp.io/finalizer"}, finalizers)

}

func TestContentConfigurationSubroutineTestSuit(t *testing.T) {
	suite.Run(t, new(ContentConfigurationSubroutineTestSuite))
}
