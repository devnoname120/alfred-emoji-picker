package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/deanishe/awgo"
	"github.com/devnoname120/alfred-emoji-picker/scoring"
	"github.com/hackebrot/turtle"
	"github.com/samber/lo"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func run() {
	query := os.Args[1]
	results := search(query)

	for _, result := range results {
		wf.NewItem(result.Name).
			Subtitle(fmt.Sprintf("Input \"%s\" (%s) into foremost application", result.Char, result.Name)).
			Arg(result.Char).
			Icon(&aw.Icon{Value: fmt.Sprintf("emojis/%s.png", result.Name)}).
			Valid(true)
	}

	wf.WarnEmpty("No matching emojis", "Try a different query?")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}

// FIXME: `che` doesn't return the checkbox in Alfred
func search(query string) []*turtle.Emoji {
	if query == "" {
		return make([]*turtle.Emoji, 0)
	}

	nameExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name == query
	})

	namePrefixMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name != query && strings.HasPrefix(e.Name, query)
	})

	nameContainMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Name != query && !strings.HasPrefix(e.Name, query) && strings.Contains(e.Name, query)
	})

	nameMatches := lo.Flatten([][]*turtle.Emoji{nameExactMatches, namePrefixMatches, nameContainMatches})

	sort.Stable(scoring.SortedByScoreDsc{Query: query, Emojis: &nameMatches})

	keywordExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword == query {
				return true
			}
		}
		return false
	})

	keywordContainMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword != query && strings.HasPrefix(keyword, query) {
				return true
			}
		}
		return false
	})

	keywordMatches := lo.Flatten([][]*turtle.Emoji{keywordExactMatches, keywordContainMatches})
	sort.Stable(scoring.SortedByScoreDsc{Query: query, Emojis: &keywordMatches})

	categoryExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category == query
	})

	categoryMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category != query && strings.HasPrefix(e.Category, query)
	})

	results := [][]*turtle.Emoji{
		nameMatches,
		keywordMatches,
		categoryExactMatches,
		categoryMatches,
	}

	consolidated := lo.Uniq(lo.Flatten(results))
	return consolidated
}
