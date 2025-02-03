package scoring

import (
	"fmt"

	"github.com/devnoname120/turtle"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("score() function", func() {
	queries := map[string][]string{
		"raised":           {"🙌", "✋", "🤚", "🖐️", "🤨"},
		"rais":             {"🙌", "✋", "🤚", "🖐️", "🤨"},
		"hand":             {"✋", "🤚", "🖐️", "🤝"},
		"smile":            {"🙂", "😊"},
		"slight":           {"🙂", "🙁"},
		"slightly_smiling": {"🙂", "🙁"},
		"sli":              {"🙂", "🙁"},
		"ok":               {"👌", "🆗"},
		"check":            {"✅", "✔️", "☑️", "🏁", "🏨"},
		"che":              {"✅", "✔️", "☑️", "🏁", "🏨"},
	}

	for searchQuery, emojis := range queries {
		// capture the loop variable
		q := searchQuery
		When(fmt.Sprintf("the search Query is '%s'", q), func() {
			for i := 0; i < len(emojis)-1; i++ {
				emojiKey1 := emojis[i]
				emojiKey2 := emojis[i+1]

				It(fmt.Sprintf("%s > %s", emojiKey1, emojiKey2), func() {
					leftEmoji := turtle.EmojisByChar[emojiKey1]
					rightEmoji := turtle.EmojisByChar[emojiKey2]
					Expect(leftEmoji).NotTo(BeNil())
					Expect(rightEmoji).NotTo(BeNil())
					got := IsScoredHigher(q, leftEmoji, rightEmoji)
					Expect(got).To(BeTrue())
				})

				It(fmt.Sprintf("%s < %s", emojiKey2, emojiKey1), func() {
					leftEmoji := turtle.EmojisByChar[emojiKey1]
					rightEmoji := turtle.EmojisByChar[emojiKey2]
					Expect(leftEmoji).NotTo(BeNil())
					Expect(rightEmoji).NotTo(BeNil())
					got := IsScoredHigher(q, rightEmoji, leftEmoji)
					Expect(got).To(BeFalse())
				})
			}
		})
	}
})
