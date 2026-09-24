// Package main implements the miniflux-sieve command.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	miniflux "miniflux.app/client"
)

func run(config config, rules []rule) error {
	client := miniflux.New(config.MinifluxURL, config.MinifluxAPIKey)
	err := applyRules(client, rules)
	if err != nil {
		return fmt.Errorf("applying rules: %w", err)
	}
	return nil
}

func main() {
	config, err := readConfig()
	if err != nil {
		slog.Error("Error reading config", "error", err)
		os.Exit(1)
	}

	rules, err := readRules()
	if err != nil {
		slog.Error("Error reading rules", "error", err)
		os.Exit(1)
	}

	if !config.MinifluxSieveDaemon {
		err = run(config, rules)
		if err != nil {
			slog.Error("miniflux-sieve failed", "error", err)
			os.Exit(1)
		}
	} else {
		for {
			err = run(config, rules)
			if err != nil {
				slog.Error("miniflux-sieve failed", "error", err)
			}
			time.Sleep(config.MinifluxSieveInterval)
		}
	}

}
