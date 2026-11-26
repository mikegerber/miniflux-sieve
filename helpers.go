package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	miniflux "miniflux.app/client"
)

// parseDuration parses durations of whole days, for example "7d", to a time.Duration
func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[len(s)-1:] != "d" {
		return 0, fmt.Errorf("duration not in format <days>d: %q", s)
	}

	amount := s[:len(s)-1]
	days, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse number of days: %w", err)
	}

	const maxDays = int64(math.MaxInt64) / int64(24*time.Hour)
	if days > maxDays {
		return 0, fmt.Errorf("duration out of range: %q", s)
	}

	return time.Duration(days) * 24 * time.Hour, nil
}

// matchStringAny matches strings against multiple regular expressions, true if any matches
//
// Patterns are case-insensitive by default, matching stops at the first match, and an invalid pattern returns
// an error.
func matchStringAny(patterns []string, s string) (bool, error) {
	for _, pattern := range patterns {
		pattern = "(?i)" + pattern // case-insensitive by default
		patternMatched, err := regexp.MatchString(pattern, s)
		if err != nil {
			return false, fmt.Errorf("Error in regex %q: %w", pattern, err)
		}
		if patternMatched {
			return true, nil
		}
	}
	return false, nil
}

// all returns the first false or error result, or true and nil if all predicates pass.
func all(predicates []entryPredicate, entry miniflux.Entry) (bool, error) {
	for _, p := range predicates {
		ok, err := p(entry)
		if err != nil || !ok {
			return ok, err
		}
	}
	return true, nil
}
