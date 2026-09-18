package main

import (
	miniflux "miniflux.app/client"
)

func main() {
	config := read_config()
	rules := read_rules()

	client := miniflux.New(config.MinifluxUrl, config.MinifluxApiKey)
	apply_rules(client, rules)
}
