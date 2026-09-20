package main

import (
	"errors"
	"log/slog"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// VERY basic function to parse e.g. "7d" to a time.Duration
func parseDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)

	unit := s[len(s)-1:]
	if unit != "d" {
		return math.MaxInt64, errors.New("duration not in format <x>d")
	}

	amount := s[:len(s)-1]
	days, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return math.MaxInt64, err
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
