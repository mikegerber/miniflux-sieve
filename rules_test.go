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

func TestRuleFilterTaggedCaseInsensitive(t *testing.T) {
	rule := rule{
		When: condition{
			Tagged: StringList{"grover", "cookie monster"},
		},
	}
	entry := miniflux.Entry{Title: "News from Sesame Street", Tags: []string{"Grover"}}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned unexpected error: %v", err)
	}
	if !matched {
		t.Errorf("ruleFilter() = %v, want true for case-insensitive tag match", matched)
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

func TestRuleFilterRequiresAllConditions(t *testing.T) {
	rule := rule{
		When: condition{
			Tagged:       StringList{"bert"},
			TitleMatches: StringList{`Sesame Street`},
		},
	}
	entry := miniflux.Entry{Title: "News from Sesame Street", Tags: []string{"oscar"}}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if matched {
		t.Fatal("ruleFilter() = true, want false when only the title condition matches")
	}
}

func TestRuleFilterInvalidTitlePattern(t *testing.T) {
	rule := rule{
		When: condition{
			TitleMatches: StringList{`[`},
		},
	}
	entry := miniflux.Entry{Title: "News from Sesame Street"}

	matched, err := ruleFilter(rule)(entry)
	if err == nil {
		t.Fatal("ruleFilter() returned no error for an invalid regular expression")
	}
	if matched {
		t.Fatal("ruleFilter() = true, want false for an invalid regular expression")
	}
}

func TestRuleFilterNotTitleMatchesReject(t *testing.T) {
	rule := rule{
		When: condition{
			NotTitleMatches: StringList{`Sesame Street`},
		},
	}
	entry := miniflux.Entry{Title: "News from Sesame Street"}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if matched {
		t.Fatal("ruleFilter() = true, want false for an excluded title")
	}
}

func TestRuleFilterURLMatches(t *testing.T) {
	rule := rule{
		When: condition{
			URLMatches: StringList{`^https://example\.com/news/`},
		},
	}
	entry := miniflux.Entry{
		Title: "News from Sesame Street",
		URL:   "https://example.com/news/sesame-street",
	}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if !matched {
		t.Fatal("ruleFilter() = false, want true for a matching URL")
	}
}

func TestRuleFilterOlderThan(t *testing.T) {
	t.Skip("TODO: test entries older and newer than the threshold, and an invalid duration")
}

// test that a title which does not match the excluded pattern passes
func TestRuleFilterNotTitleMatchesAccept(t *testing.T) {
	rule := rule{
		When: condition{
			NotTitleMatches: StringList{`sesame street|muppets`},
		},
	}
	entry := miniflux.Entry{
		Title: "The Cookie Monster is at large again!",
	}
	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if !matched {
		t.Fatal("ruleFilter() = false, want true for a non-matching title")
	}
}

// Test that an entry passes when all configured conditions match
func TestRuleFilterAllConditionsMatch(t *testing.T) {
	rule := rule{
		When: condition{
			Tagged:       StringList{`bert`, `ernie`},
			TitleMatches: StringList{`breakfast`, `lunch`, `dinner`},
		},
	}
	entry := miniflux.Entry{
		Tags:  []string{"ernie", "cookie monster"},
		Title: "Ernie meets Cookie Monster for exclusive cookie lunch",
	}

	matched, err := ruleFilter(rule)(entry)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if !matched {
		t.Fatal("ruleFilter() = false, want true for entry that matches all conditions")
	}
}

func TestApplyRulesSkipsRuleWithNoConditions(t *testing.T) {
	t.Skip("TODO: test that applying a rule with no valid conditions skips it")
}
