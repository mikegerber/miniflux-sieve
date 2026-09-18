package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

// VERY basic function to parse e.g. "7d" to a time.Duration
func ParseDuration(s string) (time.Duration, error) {
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
