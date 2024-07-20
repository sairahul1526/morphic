package config

import (
	"fmt"

	"github.com/raystack/salt/db"
)

type Config struct {
	Service  serviceConfig `mapstructure:"service"`
	Metrics  metricsConfig `mapstructure:"metrics"`
	Log      logConfig     `mapstructure:"log"`
	Database db.Config     `mapstructure:"database"`
	Auth     AuthConfig    `mapstructure:"auth"`
}

type AuthConfig struct {
	Secret        string
	SecretEncoded string `mapstructure:"secret_encoded"`
	Expiry        int    `mapstructure:"expiry" default:"24"`
}

type serviceConfig struct {
	Host           string   `mapstructure:"host" default:""`
	Port           int      `mapstructure:"port" default:"8080"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type metricsConfig struct {
	Port    int    `mapstructure:"port" default:"9090"`
	Enabled string `mapstructure:"enabled"`
}

type logConfig struct {
	Level string `mapstructure:"level" default:"info"`
}

func (serviceCfg serviceConfig) Addr() string {
	return fmt.Sprintf("%s:%d", serviceCfg.Host, serviceCfg.Port)
}
