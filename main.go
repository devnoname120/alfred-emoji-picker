package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/deanishe/awgo"
	"github.com/devnoname120/alfred-emoji-picker/scoring"
	"github.com/devnoname120/alfred-emoji-picker/usage"
	"github.com/devnoname120/turtle"
	"github.com/samber/lo"
)

var wf *aw.Workflow

func run() {
	if len(os.Args) < 2 {
		wf.Fatal("missing query argument")
	}

	query := os.Args[1]
	showUsageCount := query == ""
	stats := usage.Stats{}
	if showUsageCount {
		stats = usageStats()
	}
	results := search(query, stats)

	for _, result := range results {
		isResetItem := result.Char == "__reset_frequent__"

		subtitle := fmt.Sprintf("Input \"%s\" (%s) into foremost application", result.Char, result.Slug)
		if showUsageCount && !isResetItem {
			subtitle = fmt.Sprintf("%s · Used %d times", subtitle, stats.Count(result.Char))
		}

		item := wf.NewItem(result.Name).
			Arg(result.Char).
			Valid(true)

		if !isResetItem {
			item.Subtitle(subtitle)
		}

		item.Icon(&aw.Icon{Value: fmt.Sprintf("emojis/%s.png", result.Slug)})

		if !showUsageCount {
			item.UID(result.Char)
		}
	}

	wf.WarnEmpty("No matching emojis", "Try a different query?")
	wf.SendFeedback()
}

func usageStats() usage.Stats {
	stats, err := usage.Load()
	if err != nil {
		log.Printf("load emoji usage stats for subtitle: %v", err)
		return usage.Stats{}
	}

	return stats
}

func main() {
	if len(os.Args) > 1 {
		var err error
		handled := true
		switch os.Args[1] {
		case "--record":
			if len(os.Args) != 3 {
				err = fmt.Errorf("usage: alfred-emoji-picker --record EMOJI")
			} else {
				err = usage.Increment(os.Args[2])
			}
		case "--reset-frequent":
			if len(os.Args) != 2 {
				err = fmt.Errorf("usage: alfred-emoji-picker --reset-frequent")
			} else {
				err = usage.Reset()
			}
		default:
			handled = false
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if handled {
			return
		}
	}

	wf = aw.New()
	wf.Run(run)
}

func search(query string, stats usage.Stats) []*turtle.Emoji {
	if query == "" {
		return frequentResults(stats)
	}

	nameSlugExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name == query || e.Slug == query
	})

	nameSlugPrefixMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name != query && e.Slug != query && (strings.HasPrefix(e.Name, query) || strings.HasPrefix(e.Slug, query))
	})

	nameSlugContainMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name != query && e.Slug != query && !strings.HasPrefix(e.Name, query) && !strings.HasPrefix(e.Slug, query) && (strings.Contains(e.Name, query) || strings.Contains(e.Slug, query))
	})

	nameSlugMatches := lo.Flatten([][]*turtle.Emoji{nameSlugExactMatches, nameSlugPrefixMatches, nameSlugContainMatches})
	sort.Stable(scoring.SortedByScoreDsc{Query: query, Emojis: &nameSlugMatches})

	keywordExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword == query {
				return true
			}
		}
		return false
	})

	keywordPrefixMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword != query && strings.HasPrefix(keyword, query) {
				return true
			}
		}
		return false
	})

	keywordMatches := lo.Flatten([][]*turtle.Emoji{keywordExactMatches, keywordPrefixMatches})
	sort.Stable(scoring.SortedByScoreDsc{Query: query, Emojis: &keywordMatches})

	categoryExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category == query
	})

	categoryPrefixMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category != query && strings.HasPrefix(e.Category, query)
	})

	sort.Stable(scoring.SortedByScoreDsc{Query: query, Emojis: &nameSlugMatches})
	sort.Stable(scoring.SortedByScoreDsc{Query: query, Emojis: &keywordMatches})

	results := [][]*turtle.Emoji{
		nameSlugMatches,
		keywordMatches,
		categoryExactMatches,
		categoryPrefixMatches,
	}

	consolidated := lo.Uniq(lo.Flatten(results))
	if shouldOfferResetFrequent(query) {
		resetItem := &turtle.Emoji{
			Name: "Reset frequent emoji",
			Slug: "wastebasket",
			Char: "__reset_frequent__",
		}
		consolidated = append([]*turtle.Emoji{resetItem}, consolidated...)
	}

	return consolidated
}

func frequentResults(stats usage.Stats) []*turtle.Emoji {
	limit := frequentEmojiLimit()
	if limit <= 0 || len(stats) == 0 {
		return []*turtle.Emoji{}
	}

	results := make([]*turtle.Emoji, 0, len(stats))
	for emojiChar, count := range stats {
		if count <= 0 {
			continue
		}

		emoji := emojiByUsageKey(emojiChar)
		if emoji == nil {
			continue
		}

		results = append(results, emoji)
	}

	sortFrequentByUsage(results, stats)
	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func frequentEmojiLimit() int {
	const defaultLimit = 20

	raw := strings.TrimSpace(os.Getenv("frequent_emoji_limit"))
	if raw == "" {
		return defaultLimit
	}

	limit, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("invalid frequent_emoji_limit %q, using default %d", raw, defaultLimit)
		return defaultLimit
	}

	if limit < 0 {
		return 0
	}

	return limit
}

func shouldOfferResetFrequent(query string) bool {
	query = strings.TrimSpace(strings.ToLower(query))
	return query != "" && strings.HasPrefix(query, "reset")
}

func emojiByUsageKey(emojiChar string) *turtle.Emoji {
	for _, emoji := range turtle.EmojisByChar {
		if usage.NormalizeEmoji(emoji.Char) == emojiChar {
			return emoji
		}
	}

	return nil
}

func sortFrequentByUsage(emojis []*turtle.Emoji, stats usage.Stats) {
	slices.SortStableFunc(emojis, func(left, right *turtle.Emoji) int {
		usageLeft := stats.Count(left.Char)
		usageRight := stats.Count(right.Char)
		if usageLeft != usageRight {
			return cmp.Compare(usageRight, usageLeft)
		}

		if left.Name != right.Name {
			return cmp.Compare(left.Name, right.Name)
		}

		return cmp.Compare(left.Char, right.Char)
	})
}
