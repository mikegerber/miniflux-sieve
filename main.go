package main

import (
	miniflux "miniflux.app/client"
)

func main() {
	config := readConfig()
	rules := readRules()

	client := miniflux.New(config.MinifluxURL, config.MinifluxAPIKey)
	applyRules(client, rules)
}
