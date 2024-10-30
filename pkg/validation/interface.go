package validation

import "github.com/hashicorp/go-multierror"

type ExtensionConfiguration interface {
	Validate([]byte, string) (string, error, *multierror.Error)
	WithSchema([]byte) error
}
