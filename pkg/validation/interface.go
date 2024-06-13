package validation

type ExtensionConfiguration interface {
	Validate([]byte, string) (string, error)
	WithSchema([]byte) error
}
