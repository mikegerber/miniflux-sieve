package main

import (
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/adrg/xdg"
	"github.com/goccy/go-yaml"
	miniflux "miniflux.app/client"
)

type RulesFile struct {
	Version string `yaml:"version"`
	Rules   []Rule `yaml:"rules"`
}

type Rule struct {
	Name string    `yaml:"name"`
	Feed string    `yaml:"feed"`
	When Condition `yaml:"when"`
}

type Condition struct {
	OlderThan       string     `yaml:"older_than"`
	Tagged          StringList `yaml:"tagged"`
	TitleMatches    StringList `yaml:"title_matches"`
	NotTitleMatches StringList `yaml:"not_title_matches"`
	URLMatches      StringList `yaml:"url_matches"`
	ContentMatches  StringList `yaml:"content_matches"`
}

func ReadRules() []Rule {
	rulesPath, err := xdg.ConfigFile("miniflux-sieve/rules.yml")
	if err != nil {
		slog.Error("Could not find rules.yml", "error", err)
		os.Exit(1)
	}

	yml, err := os.ReadFile(rulesPath)
	if err != nil {
		slog.Error("Could not read rules file", "error", err)
		os.Exit(1)
	}

	var rulesFile RulesFile

	if err := yaml.Unmarshal([]byte(yml), &rulesFile); err != nil {
		slog.Error("Could not unmarshale rules file from YAML", "error", err)
		os.Exit(1)
	}
	if rulesFile.Version != "1" {
		slog.Error("Unsupported rules version", "version", rulesFile.Version)
		os.Exit(1)
	}

	rules := rulesFile.Rules

	return rules
}

// Construct a function that filters according to the Rule.
func RuleFilter(rule Rule) func(entry miniflux.Entry) bool {
	filter := func(entry miniflux.Entry) bool {
		filterOlderThan := func(entry miniflux.Entry) bool {
			// Always true if not filtered on tags
			if rule.When.OlderThan == "" {
				return true
			}

			olderThanDuration, err := ParseDuration(rule.When.OlderThan)
			if err != nil {
				slog.Error("Invalid duration", "rule", rule)
				os.Exit(1)
			}

			return entry.Date.Before(time.Now().Add(-olderThanDuration))
		}

		filterTagged := func(entry miniflux.Entry) bool {
			// Always true if not filtered on tags
			if len(rule.When.Tagged) == 0 {
				return true
			}
			matched := false
			for _, tag := range entry.Tags {
				if slices.Contains(rule.When.Tagged, tag) {
					matched = true
					break
				}
			}
			return matched
		}

		filterTitleMatches := func(entry miniflux.Entry) bool {
			// Always true if not filtered on title match
			if len(rule.When.TitleMatches) == 0 {
				return true
			}
			return MatchStringAny(rule.When.TitleMatches, entry.Title)
		}

		filterNotTitleMatches := func(entry miniflux.Entry) bool {
			// Always true if not filtered on "not title match"
			if len(rule.When.NotTitleMatches) == 0 {
				return true
			}
			return !MatchStringAny(rule.When.NotTitleMatches, entry.Title)
		}

		filterURLMatches := func(entry miniflux.Entry) bool {
			// Always true if not filtered on URL match
			if len(rule.When.URLMatches) == 0 {
				return true
			}
			return MatchStringAny(rule.When.URLMatches, entry.URL)
		}

		return filterOlderThan(entry) &&
			filterTagged(entry) &&
			filterTitleMatches(entry) &&
			filterNotTitleMatches(entry) &&
			filterURLMatches(entry)
	}
	return filter
}

func ApplyRules(client *miniflux.Client, rules []Rule) {
	feeds, err := client.Feeds()
	if err != nil {
		slog.Error("Could not get feeds from Miniflux", "error", err)
		os.Exit(1)
	}

	for _, rule := range rules {
		if rule.Feed == "" {
			slog.Warn("Skipping rule with no feed", "rule", rule)
			continue
		}
		if rule.When.OlderThan == "" &&
			len(rule.When.Tagged) == 0 &&
			len(rule.When.TitleMatches) == 0 &&
			len(rule.When.NotTitleMatches) == 0 &&
			len(rule.When.URLMatches) == 0 &&
			len(rule.When.ContentMatches) == 0 {
			slog.Warn("Skipping rule with no valid conditions", "rule", rule)
			continue
		}
		if len(rule.When.ContentMatches) > 0 {
			slog.Warn("Skipping rule: content_matches not implemented yet", "rule", rule)
			continue
		}

		filter := RuleFilter(rule)

		for _, feed := range feeds {
			if rule.Feed == feed.FeedURL {
				slog.Info("Applying rule to", "feed", rule.Feed)

				// Go through all unread and non-starred items and, if they meet the conditions,
				// mark them as read.

				entries, err := client.Entries(&miniflux.Filter{
					FeedID:  feed.ID,
					Status:  miniflux.EntryStatusUnread,
					Starred: miniflux.FilterNotStarred,
				})
				if err != nil {
					slog.Error("Error getting feed entries", "feed", feed)
					continue
				}

				var idsToMarkRead []int64
				for _, entry := range entries.Entries {
					if filter(*entry) {
						idsToMarkRead = append(idsToMarkRead, entry.ID)
					}
				}
				client.UpdateEntries(
					idsToMarkRead,
					miniflux.EntryStatusRead,
				)

			}
		}
	}
}
