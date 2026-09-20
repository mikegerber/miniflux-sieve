package main

import (
	"log/slog"
	"os"

	"github.com/adrg/xdg"
	"github.com/goccy/go-yaml"
)

type config struct {
	MinifluxURL    string `yaml:"MINIFLUX_URL"`
	MinifluxAPIKey string `yaml:"MINIFLUX_API_KEY"`
}

func readConfig() config {
	configPath, err := xdg.ConfigFile("miniflux-sieve/config.yml")
	if err != nil {
		slog.Error("Could not find config.yml", "error", err)
		os.Exit(1)
	}

	yml, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("Could not read config file", "error", err)
		os.Exit(1)
	}

	var cfg config

	if err := yaml.Unmarshal([]byte(yml), &cfg); err != nil {
		slog.Error("Could not unmarshal config file from YAML", "error", err)
		os.Exit(1)
	}

	return cfg
}
