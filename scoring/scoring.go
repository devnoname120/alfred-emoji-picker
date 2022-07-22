package scoring

import (
	"strings"

	"github.com/hackebrot/turtle"
)

func IsScoredHigher(query string, emojiLeft *turtle.Emoji, emojiRight *turtle.Emoji) bool {
	scoreLeft := Score(query, emojiLeft)
	scoreRight := Score(query, emojiRight)

	return scoreLeft >= scoreRight
}

func Score(query string, emoji *turtle.Emoji) int {
	emojiNicknamePriorities := map[string][]string{
		"rais": {"🙌", "✋", "🤚", "🖐", "🤨"},
		"han":  {"✋", "🤚", "🖐", "🤝"},
		"smil": {"🙂", "😊"},
		"sli":  {"🙂", "🙁"},
		"ok":   {"👌", "🆗"},
		"che":  {"✅", "✔️", "☑️", "🏁", "🏨"},
		// This doesn't work because this is a keyword, and right now this function only uses names
		"lo":  {"😆", "🤣", "🍭"},
		"spa": {"✨"},
		"pra": {"🙏"},
		"cry": {"😢", "😭"},
		"thu": {"👍", "👎"},
	}

	//if emojiPriorities := emojiNicknamePriorities[query]; emojiPriorities != nil {
	//	return positionToScore(emoji.Char, emojiPriorities)
	//}

	for emojiNickname, emojiPriorities := range emojiNicknamePriorities {
		if strings.HasPrefix(query, emojiNickname) {
			return positionToScore(emoji.Char, emojiPriorities)
		}
	}

	if emoji.Name == query {
		return 2
	}

	if strings.HasPrefix(emoji.Name, query) {
		return 1
	}

	return 0
}

func positionToScore(emojiChar string, emojiChars []string) int {
	for _, curEmojiChar := range emojiChars {
		if emojiChar == curEmojiChar {
			return 2 + len(emojiChars)
		}
	}

	return 0
}

type SortedByScoreDsc struct {
	Query  string
	Emojis *[]*turtle.Emoji
}

func (s SortedByScoreDsc) Len() int {
	return len(*s.Emojis)
}

func (s SortedByScoreDsc) Less(i, j int) bool {
	// Less = position in the list.
	// For us high score = left on the left, so we return true when score is higher
	return IsScoredHigher(s.Query, (*s.Emojis)[i], (*s.Emojis)[j])
}

func (s SortedByScoreDsc) Swap(i, j int) {
	(*s.Emojis)[i], (*s.Emojis)[j] = (*s.Emojis)[j], (*s.Emojis)[i]
}
