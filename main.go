package main

import (
	"log/slog"
	"os"

	miniflux "miniflux.app/client"
)

func main() {
	config := readConfig()

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
