package main

import (
	miniflux "miniflux.app/client"
)

func main() {
	config := ReadConfig()
	rules := ReadRules()

	client := miniflux.New(config.MinifluxURL, config.MinifluxAPIKey)
	ApplyRules(client, rules)
}
