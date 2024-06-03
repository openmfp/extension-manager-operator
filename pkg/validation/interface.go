package validation

type ContentConfigurationInterface interface {
	Validate([]byte, []byte, string) (string, error)
}
