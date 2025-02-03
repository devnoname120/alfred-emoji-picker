package scoring

import (
	"strings"

	"github.com/devnoname120/turtle"
)

func IsScoredHigher(query string, emojiLeft *turtle.Emoji, emojiRight *turtle.Emoji) bool {
	scoreLeft := Score(query, emojiLeft)
	scoreRight := Score(query, emojiRight)
	return scoreLeft > scoreRight
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

func normalizeEmoji(e string) string {
	return strings.ReplaceAll(e, "\uFE0F", "")
}

func positionToScore(emojiChar string, emojiChars []string) int {
	norm := normalizeEmoji(emojiChar)
	for i, curEmojiChar := range emojiChars {
		if norm == normalizeEmoji(curEmojiChar) {
			return 2 + (len(emojiChars) - i)
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
	return IsScoredHigher(s.Query, (*s.Emojis)[i], (*s.Emojis)[j])
}

func (s SortedByScoreDsc) Swap(i, j int) {
	(*s.Emojis)[i], (*s.Emojis)[j] = (*s.Emojis)[j], (*s.Emojis)[i]
}
