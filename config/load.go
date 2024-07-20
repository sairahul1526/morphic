package config

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/raystack/salt/config"
)

func Load(configFile string) (Config, error) {
	var cfg Config
	loader := config.NewLoader(config.WithFile(configFile))

	if err := loader.Load(&cfg); err != nil {
		if errors.As(err, &config.ConfigFileNotFoundError{}) {
			fmt.Println(err)
		} else {
			return Config{}, err
		}
	}

	// Decode the Base64 Secret Key
	decodedBytes, err := base64.StdEncoding.DecodeString(cfg.Auth.SecretEncoded)
	if err != nil {
		fmt.Println("Error:", err)
		return Config{}, err
	}

	// Convert the decoded bytes to a string
	cfg.Auth.Secret = string(decodedBytes)

	return cfg, nil
}
