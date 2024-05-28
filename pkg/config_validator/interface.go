package config_validator

type ContentConfigurationInterface interface {
	Validate([]byte, string) (string, error)
}
