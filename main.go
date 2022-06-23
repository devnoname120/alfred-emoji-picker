package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/devnoname120/alfred-emoji-picker/scoring"

	"github.com/deanishe/awgo"
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

	nameMatches := lo.Flatten([][]*turtle.Emoji{nameExactMatches, namePrefixMatches})

	sort.Stable(SortedByScoreDsc{query: query, emojis: &nameMatches})

	keywordExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword == query {
				return true
			}
		}
		return false
	})

	keywordMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		for _, keyword := range e.Keywords {
			if keyword != query && strings.HasPrefix(keyword, query) {
				return true
			}
		}
		return false
	})

	categoryExactMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category == query
	})

	categoryMatches := turtle.Filter(func(e *turtle.Emoji) bool {
		return e.Category != query && strings.HasPrefix(e.Category, query)
	})

	results := [][]*turtle.Emoji{
		nameMatches,
		keywordExactMatches,
		keywordMatches,
		categoryExactMatches,
		categoryMatches,
	}

	consolidated := lo.Flatten(results)
	return consolidated
}

type SortedByScoreDsc struct {
	query  string
	emojis *[]*turtle.Emoji
}

func (s SortedByScoreDsc) Len() int {
	return len(*s.emojis)
}

func (s SortedByScoreDsc) Less(i, j int) bool {
	// Less = position in the list.
	// For us high score = left on the left, so we return true when score is higher
	return scoring.IsScoredHigher(s.query, (*s.emojis)[i], (*s.emojis)[j])
}

func (s SortedByScoreDsc) Swap(i, j int) {
	(*s.emojis)[i], (*s.emojis)[j] = (*s.emojis)[j], (*s.emojis)[i]
}
