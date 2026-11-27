package main

import (
	"testing"
	"time"

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
		t.Fatal("ruleFilter() = false, want true")
	}
}

// TODO: Combine this with TestRuleFilterTaggedCaseInsensitive.
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
		t.Fatal("ruleFilter() = false, want true")
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
		t.Fatal("ruleFilter() = true, want false")
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
		t.Fatalf("ruleFilter() = %v, want true for case-insensitive tag match", matched)
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
		t.Fatal("ruleFilter() = false, want true")
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
		t.Error("ruleFilter() returned no error for an invalid regular expression")
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
	rule := rule{
		When: condition{
			OlderThan: "7d",
		},
	}

	entryOld := miniflux.Entry{
		Title: "not important",
		Date:  time.Now().Add(-13 * 24 * time.Hour),
	}
	matched, err := ruleFilter(rule)(entryOld)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if !matched {
		t.Fatalf("ruleFilter() = false, want true for an older entry")
	}

	entryNew := miniflux.Entry{
		Title: "not important",
		Date:  time.Now().Add(-1 * 24 * time.Hour),
	}
	matched, err = ruleFilter(rule)(entryNew)
	if err != nil {
		t.Fatalf("ruleFilter() returned an unexpected error: %v", err)
	}
	if matched {
		t.Fatalf("ruleFilter() = true, want false for a newer entry")
	}

}

func TestRuleFilterOlderThanInvalidDuration(t *testing.T) {
	invalidRule := rule{
		When: condition{
			OlderThan: "half a fortnight",
		},
	}
	_, err := ruleFilter(invalidRule)(miniflux.Entry{})
	if err == nil {
		t.Fatal("ruleFilter() did not return an error for an invalid duration")
	}
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
