package main

import (
	"fmt"
	"math"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// VERY basic function to parse e.g. "7d" to a time.Duration
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

// Match against multiple regular expressions, true if any matches
func matchStringAny(patterns []string, s string) bool {
	for _, pattern := range patterns {
		pattern = "(?i)" + pattern // case-insensitive by default
		patternMatched, err := regexp.MatchString(pattern, s)
		if err != nil {
			slog.Error("Error in regex", "error", err)
			os.Exit(1)
		}
		if patternMatched {
			return true
		}
	}
	return false
}
