package validation

type ExtensionConfiguration interface {
	Validate([]byte, string) (string, error)
	LoadSchema([]byte) error
}
