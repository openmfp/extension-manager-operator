package config

import (
	openmfpconfig "github.com/openmfp/golang-commons/config"
)

// Config struct to hold the app config
type Config struct {
	openmfpconfig.CommonServiceConfig `mapstructure:",squash"`
	IsLocal                           bool   `mapstructure:"is-local"`
	ServerPort                        string `mapstructure:"server-port"`
	Subroutines                       struct {
		ContentConfiguration struct {
			Enabled bool `mapstructure:"subroutines-contentconfiguration-enabled"`
		}
	}
}
