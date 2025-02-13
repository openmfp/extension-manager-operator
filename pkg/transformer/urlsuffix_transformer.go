package transformer

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/openmfp/extension-manager-operator/api/v1alpha1"
	"github.com/openmfp/extension-manager-operator/pkg/validation"
)

type UrlSuffixTransformer struct{}

func (*UrlSuffixTransformer) Transform(contentConfiguration *validation.ContentConfiguration, instance *v1alpha1.ContentConfiguration) error {
	if instance.Spec.RemoteConfiguration != nil {
		parsedUrl, err := url.Parse(instance.Spec.RemoteConfiguration.URL)
		if err != nil {
			return errors.Wrap(err, "failed to parse URL")
		}
		domain := fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)

		for i := range contentConfiguration.LuigiConfigFragment.Data.Nodes {
			err = transformNode(&contentConfiguration.LuigiConfigFragment.Data.Nodes[i], domain)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return nil
}

func transformNode(node *validation.Node, domain string) error {
	if node.UrlSuffix != "" {
		domain = strings.TrimRight(domain, "/")
		urlSuffix := strings.TrimLeft(node.UrlSuffix, "/")
		url := fmt.Sprintf("%s/%s", domain, urlSuffix)
		node.Url = url
		node.UrlSuffix = ""
	}
	for i := range node.Children {
		err := transformNode(&node.Children[i], domain)
		if err != nil {
			return err
		}
	}
	return nil
}
