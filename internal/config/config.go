package config

// Config struct to hold the app config
type ServerConfig struct {
	IsLocal    bool   `mapstructure:"is-local"`
	ServerPort string `mapstructure:"server-port"`
}

type OperatorConfig struct {
	Subroutines struct {
		ContentConfiguration struct {
			Enabled bool `mapstructure:"subroutines-content-configuration-enabled"`
		} `mapstructure:",squash"`
	} `mapstructure:",squash"`
}
