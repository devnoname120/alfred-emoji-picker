package main

import (
  "fmt"
  "os"
  "sort"
  "strings"

  "github.com/deanishe/awgo"
  "github.com/devnoname120/alfred-emoji-picker/scoring"
  "github.com/devnoname120/turtle"
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
      Subtitle(fmt.Sprintf("Input \"%s\" (%s) into foremost application", result.Char, result.Slug)).
      Arg(result.Char).
      Icon(&aw.Icon{Value: fmt.Sprintf("emojis/%s.png", result.Slug)}).
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
    nameSlugMatches,
    keywordMatches,
    categoryExactMatches,
    categoryMatches,
  }

  consolidated := lo.Uniq(lo.Flatten(results))
  return consolidated
}
