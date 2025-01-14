package display

import (
	"strings"
	"unicode"
)

const ellipsis = "..."

func getStringDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if unicode.Is(unicode.Han, r) ||
			unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Hangul, r) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

func TruncateAndPad(s string, maxWidth int) string {
	if maxWidth < len(ellipsis) {
		return strings.Repeat(" ", maxWidth)
	}

	displayWidth := getStringDisplayWidth(s)
	if displayWidth <= maxWidth {
		return s + strings.Repeat(" ", maxWidth-displayWidth)
	}

	runes := []rune(s)
	currentWidth := 0
	truncateIndex := 0

	for i, r := range runes {
		charWidth := 1
		if unicode.Is(unicode.Han, r) ||
			unicode.Is(unicode.Hiragana, r) ||
			unicode.Is(unicode.Katakana, r) ||
			unicode.Is(unicode.Hangul, r) {
			charWidth = 2
		}

		if currentWidth+charWidth+len(ellipsis) > maxWidth {
			break
		}
		currentWidth += charWidth
		truncateIndex = i + 1
	}

	truncated := string(runes[:truncateIndex]) + ellipsis
	return truncated + strings.Repeat(" ", maxWidth-getStringDisplayWidth(truncated))
}
