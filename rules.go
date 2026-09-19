package main

import (
	"github.com/adrg/xdg"
	"github.com/goccy/go-yaml"
	"log/slog"
	miniflux "miniflux.app/client"
	"os"
	"slices"
	"time"
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
	UrlMatches      StringList `yaml:"url_matches"`
	ContentMatches  StringList `yaml:"content_matches"`
}

func read_rules() []Rule {
	rules_path, err := xdg.ConfigFile("miniflux-sieve/rules.yml")
	if err != nil {
		slog.Error("Could not find rules.yml", "error", err)
		os.Exit(1)
	}

	yml, err := os.ReadFile(rules_path)
	if err != nil {
		slog.Error("Could not read rules file", "error", err)
		os.Exit(1)
	}

	var rulesfile RulesFile

	if err := yaml.Unmarshal([]byte(yml), &rulesfile); err != nil {
		slog.Error("Could not unmarshale rules file from YAML", "error", err)
		os.Exit(1)
	}
	if rulesfile.Version != "1" {
		slog.Error("Unsupported rules version", "version", rulesfile.Version)
		os.Exit(1)
	}

	rules := rulesfile.Rules

	return rules
}

// Construct a function that filters according to the Rule.
func RuleFilter(rule Rule) func(entry miniflux.Entry) bool {
	filter := func(entry miniflux.Entry) bool {
		filter_older_than := func(entry miniflux.Entry) bool {
			// Always true if not filtered on tags
			if rule.When.OlderThan == "" {
				return true
			}

			olderthan, err := ParseDuration(rule.When.OlderThan)
			if err != nil {
				slog.Error("Invalid duration", "rule", rule)
				os.Exit(1)
			}

			return entry.Date.Before(time.Now().Add(-olderthan))
		}

		filter_tagged := func(entry miniflux.Entry) bool {
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

		filter_title_matches := func(entry miniflux.Entry) bool {
			// Always true if not filtered on title match
			if len(rule.When.TitleMatches) == 0 {
				return true
			}
			return MatchStringAny(rule.When.TitleMatches, entry.Title)
		}

		filter_not_title_matches := func(entry miniflux.Entry) bool {
			// Always true if not filtered on "not title match"
			if len(rule.When.NotTitleMatches) == 0 {
				return true
			}
			return !MatchStringAny(rule.When.NotTitleMatches, entry.Title)
		}

		return filter_older_than(entry) &&
			filter_tagged(entry) &&
			filter_title_matches(entry) &&
			filter_not_title_matches(entry)
	}
	return filter
}

func apply_rules(client *miniflux.Client, rules []Rule) {
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
			len(rule.When.UrlMatches) == 0 &&
			len(rule.When.ContentMatches) == 0 {
			slog.Warn("Skipping rule with no valid conditions", "rule", rule)
			continue
		}
		if len(rule.When.UrlMatches) > 0 {
			slog.Warn("Skipping rule: url_matches not implemented yet", "rule", rule)
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

				var ids_to_mark_read []int64
				for _, entry := range entries.Entries {
					if filter(*entry) {
						ids_to_mark_read = append(ids_to_mark_read, entry.ID)
					}
				}
				client.UpdateEntries(
					ids_to_mark_read,
					miniflux.EntryStatusRead,
				)

			}
		}
	}
}
