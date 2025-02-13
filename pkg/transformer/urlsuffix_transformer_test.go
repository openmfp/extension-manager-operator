package transformer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/pkg/validation"
)

func TestUrlSuffixTransformer_Transform(t *testing.T) {
	transformer := &UrlSuffixTransformer{}
	tests := []struct {
		name     string
		before   *validation.ContentConfiguration
		instance *v1alpha1.ContentConfiguration
		expected *validation.ContentConfiguration
	}{
		{
			name: "Test UrlSuffixTransformer Transform",

			instance: &v1alpha1.ContentConfiguration{Spec: v1alpha1.ContentConfigurationSpec{RemoteConfiguration: &v1alpha1.RemoteConfiguration{URL: "https://test.com:9999/ui/cdm/config.json"}}},
			before: &validation.ContentConfiguration{
				LuigiConfigFragment: validation.LuigiConfigFragment{
					Data: validation.LuigiConfigData{
						Nodes: []validation.Node{
							{
								UrlSuffix: "test/#/my-ui?query=param&query2=param2",
								Children: []validation.Node{
									{
										UrlSuffix: "test/#/my-child-1?query=param&query3=param4",
									},
									{
										UrlSuffix: "test/#/my-child-2?query=param&query1=param5",
									},
								},
							},
						},
					},
				},
			},
			expected: &validation.ContentConfiguration{
				LuigiConfigFragment: validation.LuigiConfigFragment{
					Data: validation.LuigiConfigData{
						Nodes: []validation.Node{
							{
								Url:       "https://test.com:9999/test/#/my-ui?query=param&query2=param2",
								UrlSuffix: "",
								Children: []validation.Node{
									{
										UrlSuffix: "",
										Url:       "https://test.com:9999/test/#/my-child-1?query=param&query3=param4",
									},
									{
										UrlSuffix: "",
										Url:       "https://test.com:9999/test/#/my-child-2?query=param&query1=param5",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		err := transformer.Transform(test.before, test.instance)
		assert.NoError(t, err)
		assert.Equal(t, test.before, test.expected)
	}
}
