package validation

type ContentConfigurationInterface interface {
	Validate([]byte, string) (string, error)
}
