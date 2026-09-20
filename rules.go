package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/adrg/xdg"
	"github.com/goccy/go-yaml"
	miniflux "miniflux.app/client"
)

type rulesFile struct {
	Version string `yaml:"version"`
	Rules   []rule `yaml:"rules"`
}

type rule struct {
	Name string    `yaml:"name"`
	Feed string    `yaml:"feed"`
	When condition `yaml:"when"`
}

type condition struct {
	OlderThan       string     `yaml:"older_than"`
	Tagged          StringList `yaml:"tagged"`
	TitleMatches    StringList `yaml:"title_matches"`
	NotTitleMatches StringList `yaml:"not_title_matches"`
	URLMatches      StringList `yaml:"url_matches"`
	ContentMatches  StringList `yaml:"content_matches"`
}

func readRules() ([]rule, error) {
	rulesPath, err := xdg.ConfigFile("miniflux-sieve/rules.yml")
	if err != nil {
		return nil, fmt.Errorf("Could not find rules.yml: %w", err)
	}

	yml, err := os.ReadFile(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("Could not read rules file: %w", err)
	}

	var rulesFile rulesFile

	if err := yaml.Unmarshal(yml, &rulesFile); err != nil {
		return nil, fmt.Errorf("Could not unmarshal rules file from YAML: %w", err)
	}
	if rulesFile.Version != "1" {
		return nil, fmt.Errorf("Unsupported rules version: %q", rulesFile.Version)
	}

	rules := rulesFile.Rules

	return rules, nil
}

type entryPredicate func(miniflux.Entry) (bool, error)

// Construct a function that filters according to the Rule.
func ruleFilter(rule rule) entryPredicate {
	return func(entry miniflux.Entry) (bool, error) {
		filterOlderThan := func(entry miniflux.Entry) (bool, error) {
			// Always true if not filtered on tags
			if rule.When.OlderThan == "" {
				return true, nil
			}

			olderThanDuration, err := parseDuration(rule.When.OlderThan)
			if err != nil {
				return false, err
			}

			return entry.Date.Before(time.Now().Add(-olderThanDuration)), nil
		}

		filterTagged := func(entry miniflux.Entry) (bool, error) {
			// Always true if not filtered on tags
			if len(rule.When.Tagged) == 0 {
				return true, nil
			}
			for _, tag := range entry.Tags {
				if slices.Contains(rule.When.Tagged, tag) {
					return true, nil
				}
			}
			return false, nil
		}

		filterTitleMatches := func(entry miniflux.Entry) (bool, error) {
			// Always true if not filtered on title match
			if len(rule.When.TitleMatches) == 0 {
				return true, nil
			}
			matched, err := matchStringAny(rule.When.TitleMatches, entry.Title)
			return matched, err
		}

		filterNotTitleMatches := func(entry miniflux.Entry) (bool, error) {
			// Always true if not filtered on "not title match"
			if len(rule.When.NotTitleMatches) == 0 {
				return true, nil
			}
			matched, err := matchStringAny(rule.When.NotTitleMatches, entry.Title)
			return !matched, err
		}

		filterURLMatches := func(entry miniflux.Entry) (bool, error) {
			// Always true if not filtered on URL match
			if len(rule.When.URLMatches) == 0 {
				return true, nil
			}
			matched, err := matchStringAny(rule.When.URLMatches, entry.URL)
			return matched, err
		}

		return all(
			[]entryPredicate{filterOlderThan, filterTagged, filterTitleMatches, filterNotTitleMatches, filterURLMatches},
			entry)
	}
}

func applyRules(client *miniflux.Client, rules []rule) error {
	feeds, err := client.Feeds()
	if err != nil {
		return fmt.Errorf("Could not get feeds from Miniflux: %w", err)
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

		filter := ruleFilter(rule)

		for _, feed := range feeds {
			if rule.Feed == feed.FeedURL {
				slog.Info("Applying rule", "rule", rule)

				// Go through all unread and non-starred items and, if they meet the conditions,
				// mark them as read.

				entries, err := client.Entries(&miniflux.Filter{
					FeedID:  feed.ID,
					Status:  miniflux.EntryStatusUnread,
					Starred: miniflux.FilterNotStarred,
				})
				if err != nil {
					return fmt.Errorf("Error getting feed entries for feed %q: %w", rule.Feed, err)
				}

				var idsToMarkRead []int64
				for _, entry := range entries.Entries {
					markRead, err := filter(*entry)
					if err != nil {
						return fmt.Errorf("Error filtering: %w", err)
					}
					if markRead {
						idsToMarkRead = append(idsToMarkRead, entry.ID)
					}
				}

				if len(idsToMarkRead) > 0 {
					err := client.UpdateEntries(idsToMarkRead, miniflux.EntryStatusRead)
					if err != nil {
						return fmt.Errorf("Mark entries read for feed %q: %w", rule.Feed, err)
					}
				}

				// Currently, a rule can only match (up to) one feed.
				break
			}
		}
	}

	return nil
}
