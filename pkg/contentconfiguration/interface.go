package contentconfiguration

type ContentConfigurationInterface interface {
	ValidateJSON(input []byte) (bool, []string)
	ValidateYAML(input []byte) (bool, []string)
}
