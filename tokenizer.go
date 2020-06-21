package sastrawi

import (
	"html"
	"strings"
)

// Tokenize remove symbols and URLs from sentence, then split it into words
func Tokenize(sentence string) []string {
	// Normalize sentence and remove all symbol
	sentence = strings.ToLower(sentence)
	sentence = html.UnescapeString(sentence)
	sentence = rxURL.ReplaceAllString(sentence, "")
	sentence = rxEmail.ReplaceAllString(sentence, "")
	sentence = rxTwitter.ReplaceAllString(sentence, "")
	sentence = rxEscapeStr.ReplaceAllString(sentence, "")
	sentence = rxSymbol.ReplaceAllString(sentence, " ")
	sentence = strings.TrimSpace(sentence)

	return strings.Fields(sentence)
}
