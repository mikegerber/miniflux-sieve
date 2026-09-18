package main

import (
	"github.com/adrg/xdg"
	"github.com/goccy/go-yaml"
	"log/slog"
	"os"
)

type Config struct {
	MinifluxUrl    string `yaml:"MINIFLUX_URL"`
	MinifluxApiKey string `yaml:"MINIFLUX_API_KEY"`
}

func read_config() Config {
	config_path, err := xdg.ConfigFile("miniflux-sieve/config.yml")
	if err != nil {
		slog.Error("Could not find config.yml", "error", err)
		os.Exit(1)
	}

	yml, err := os.ReadFile(config_path)
	if err != nil {
		slog.Error("Could not read config file", "error", err)
		os.Exit(1)
	}

	var config Config

	if err := yaml.Unmarshal([]byte(yml), &config); err != nil {
		slog.Error("Could not unmarshal config file from YAML", "error", err)
		os.Exit(1)
	}

	return config
}
