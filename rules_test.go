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

func TestRuleFilterTagged(t *testing.T) {
	rule := rule{
		When: condition{
			Tagged: StringList{"ernie", "bert"},
		},
	}
	entry := miniflux.Entry{Title: "News from Sesame Street", Tags: []string{"oscar", "bert"}}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if !matched {
		t.Error("ruleFilter() = false, want true")
	}

	entry = miniflux.Entry{Title: "News from Sesame Street", Tags: []string{"tiffy", "finchen"}}
}

func TestRuleFilterTaggedReject(t *testing.T) {
	rule := rule{
		When: condition{
			Tagged: StringList{"ernie", "bert"},
		},
	}
	entry := miniflux.Entry{Title: "News from Sesame Street", Tags: []string{"tiffy", "finchen"}}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if matched {
		t.Error("ruleFilter() = true, want false")
	}

}

func TestRuleWithNoConditionsIsTrue(t *testing.T) {
	rule := rule{}
	entry := miniflux.Entry{Title: "Tiffy meets Oscar"}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if !matched {
		t.Error("ruleFilter() = false, want true")
	}
}
