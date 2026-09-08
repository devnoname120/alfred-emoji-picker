package main

import (
	"testing"

	"github.com/devnoname120/alfred-emoji-picker/usage"
	"github.com/devnoname120/turtle"
)

func TestFrequentResultsResolvesNormalizedUsageKey(t *testing.T) {
	t.Setenv("frequent_emoji_limit", "10")

	results := frequentResults(usage.Stats{
		usage.NormalizeEmoji("🖐️"): 2,
	})

	if len(results) != 1 {
		t.Fatalf("expected 1 frequent result, got %d", len(results))
	}

	if results[0].Char != "🖐️" {
		t.Fatalf("expected emoji with variation selector, got %q", results[0].Char)
	}
}

func TestFrequentResultsSortByUsageDescending(t *testing.T) {
	t.Setenv("frequent_emoji_limit", "10")

	results := frequentResults(usage.Stats{
		usage.NormalizeEmoji("🕶️"): 4,
		usage.NormalizeEmoji("☀️"): 2,
		usage.NormalizeEmoji("🆒"):  1,
	})

	if len(results) != 3 {
		t.Fatalf("expected 3 frequent results, got %d", len(results))
	}

	got := []string{results[0].Char, results[1].Char, results[2].Char}
	want := []string{"🕶️", "☀️", "🆒"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected order %v, got %v", want, got)
		}
	}
}

func TestSearchKeepsExactKeywordBeforePrefixKeyword(t *testing.T) {
	results := search("cool", nil)

	var coolButtonIndex = -1
	var sunglassesIndex = -1
	for i, result := range results {
		switch result.Char {
		case turtle.EmojisByChar["🆒"].Char:
			coolButtonIndex = i
		case turtle.EmojisByChar["🕶️"].Char:
			sunglassesIndex = i
		}
	}

	if coolButtonIndex == -1 || sunglassesIndex == -1 {
		t.Fatalf("expected both cool button and sunglasses in results, got cool=%d sunglasses=%d", coolButtonIndex, sunglassesIndex)
	}

	if coolButtonIndex > sunglassesIndex {
		t.Fatalf("expected exact keyword match before prefix keyword match, got cool=%d sunglasses=%d", coolButtonIndex, sunglassesIndex)
	}
}

func TestSearchOffersResetFrequentFirst(t *testing.T) {
	results := search("reset", nil)

	if len(results) == 0 {
		t.Fatalf("expected reset item in results")
	}

	if results[0].Char != "__reset_frequent__" {
		t.Fatalf("expected reset item first, got %q", results[0].Char)
	}
}

func TestSearchOffersResetFrequentForLongerResetQuery(t *testing.T) {
	results := search("reset frequent", nil)

	if len(results) == 0 {
		t.Fatalf("expected reset item in results")
	}

	if results[0].Char != "__reset_frequent__" {
		t.Fatalf("expected reset item first, got %q", results[0].Char)
	}
}
