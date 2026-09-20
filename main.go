package main

import (
	"log/slog"
	"os"

	miniflux "miniflux.app/client"
)

func main() {
	config, err := readConfig()
	if err != nil {
		slog.Error("Error reading config", "error", err)
		os.Exit(1)
	}

	rules, err := readRules()
	if err != nil {
		slog.Error("Error reading rules", "error", err)
	}

	client := miniflux.New(config.MinifluxURL, config.MinifluxAPIKey)
	err = applyRules(client, rules)
	if err != nil {
		slog.Error("Error applying rules", "error", err)
		os.Exit(1)
	}
}
