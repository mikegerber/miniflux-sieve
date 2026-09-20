package main

import (
	"testing"

	miniflux "miniflux.app/client"
)

func TestRuleFilterTitleMatches(t *testing.T) {
	rule := rule{
		When: condition{
			TitleMatches: StringList{`Sesame Street`},
		},
	}
	entry := miniflux.Entry{Title: "News from Sesame Street"}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if !matched {
		t.Error("ruleFilter() = false, want true")
	}
}
