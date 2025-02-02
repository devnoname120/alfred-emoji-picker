package scoring

import (
	"fmt"

	"github.com/devnoname120/turtle"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type stringTuple struct {
	a, b string
}

var _ = Describe("score() function", func() {
	queries := map[string][]string{
		"raised":           {"🙌", "✋", "🤚", "🖐", "🤨"},
		"rais":             {"🙌", "✋", "🤚", "🖐", "🤨"},
		"hand":             {"✋", "🤚", "🖐", "🤝"},
		"smile":            {"🙂", "😊"},
		"slight":           {"🙂", "🙁"},
		"slightly_smiling": {"🙂", "🙁"},
		"sli":              {"🙂", "🙁"},
		"ok":               {"👌", "🆗"},
		"check":            {"✅", "✔️", "☑️", "🏁", "🏨"},
		"che":              {"✅", "✔️", "☑️", "🏁", "🏨"},
	}

	for searchQuery, emojis := range queries {
		When(fmt.Sprintf("the search Query is '%s'", searchQuery), func() {
			for i := 0; i < len(emojis)-1; i++ {
				leftEmoji := turtle.EmojisByChar[emojis[i]]
				rightEmoji := turtle.EmojisByChar[emojis[i+1]]

				It(fmt.Sprintf("%s (%s) > %s (%s)", leftEmoji.Char, leftEmoji.Name, rightEmoji.Char, rightEmoji.Name), func(searchQuery string, leftEmoji *turtle.Emoji, rightEmoji *turtle.Emoji) func() {
					return func() {
						got := IsScoredHigher(searchQuery, leftEmoji, rightEmoji)
						Expect(got).To(BeTrue())
					}
				}(searchQuery, leftEmoji, rightEmoji))

				It(fmt.Sprintf("%s (%s) < %s (%s)", rightEmoji.Char, rightEmoji.Name, leftEmoji.Char, leftEmoji.Name), func(searchQuery string, leftEmoji *turtle.Emoji, rightEmoji *turtle.Emoji) func() {
					return func() {
						got := IsScoredHigher(searchQuery, rightEmoji, leftEmoji)
						Expect(got).To(BeFalse())
					}
				}(searchQuery, leftEmoji, rightEmoji))
			}
		})
	}
})
