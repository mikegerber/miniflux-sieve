package main

import (
	"fmt"
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/goccy/go-yaml"
)

type config struct {
	MinifluxURL           string        `yaml:"MINIFLUX_URL"`
	MinifluxAPIKey        string        `yaml:"MINIFLUX_API_KEY"`
	MinifluxSieveDaemon   bool          `yaml:"MINIFLUX_SIEVE_DAEMON"`
	MinifluxSieveInterval time.Duration `yaml:"MINIFLUX_SIEVE_INTERVAL"`
}

func readConfig() (config, error) {
	configPath, err := xdg.ConfigFile("miniflux-sieve/config.yml")
	if err != nil {
		return config{}, fmt.Errorf("Could not find config.yml: %w", err)
	}

	yml, err := os.ReadFile(configPath)
	if err != nil {
		return config{}, fmt.Errorf("Could not read config file: %w", err)
	}

	// Default values
	cfg := config{
		MinifluxSieveDaemon:   false,
		MinifluxSieveInterval: time.Duration(20 * time.Minute),
	}

	err = yaml.Unmarshal([]byte(yml), &cfg)
	if err != nil {
		return config{}, fmt.Errorf("Could not unmarshal config file from YAML: %w", err)
	}

	return cfg, nil
}
