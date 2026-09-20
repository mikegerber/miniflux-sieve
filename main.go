package main

import (
	"log/slog"
	"os"

	miniflux "miniflux.app/client"
)

func main() {
	config := readConfig()
	rules := readRules()

	client := miniflux.New(config.MinifluxURL, config.MinifluxAPIKey)
	err := applyRules(client, rules)
	if err != nil {
		slog.Error("Error applying rules", "error", err)
		os.Exit(1)
	}
}
