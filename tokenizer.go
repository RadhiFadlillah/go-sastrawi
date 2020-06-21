package sastrawi

import (
	"html"
	"regexp"
	"strings"
)

// Tokenizer is object for tokenizing sentence
type Tokenizer struct {
	rxURL        *regexp.Regexp
	rxEmail      *regexp.Regexp
	rxTwitter    *regexp.Regexp
	rxEscapeStr  *regexp.Regexp
	rxSymbol     *regexp.Regexp
	rxWhitespace *regexp.Regexp
}

// NewTokenizer returns new Tokenizer
func NewTokenizer() Tokenizer {
	return Tokenizer{
		rxURL:        regexp.MustCompile(`(www\.|https?|s?ftp)\S+`),
		rxEmail:      regexp.MustCompile(`\S+@\S+`),
		rxTwitter:    regexp.MustCompile(`(@|#)\S+`),
		rxEscapeStr:  regexp.MustCompile(`&.*;`),
		rxSymbol:     regexp.MustCompile(`[^a-z\s]`),
		rxWhitespace: regexp.MustCompile(`\s+`),
	}
}

// Tokenize remove symbols and URLs from sentence, then split it into words
func (tokenizer Tokenizer) Tokenize(sentence string) []string {
	// Normalize sentence and remove all symbol
	sentence = strings.ToLower(sentence)
	sentence = html.UnescapeString(sentence)
	sentence = tokenizer.rxURL.ReplaceAllString(sentence, "")
	sentence = tokenizer.rxEmail.ReplaceAllString(sentence, "")
	sentence = tokenizer.rxTwitter.ReplaceAllString(sentence, "")
	sentence = tokenizer.rxEscapeStr.ReplaceAllString(sentence, "")
	sentence = tokenizer.rxSymbol.ReplaceAllString(sentence, " ")
	sentence = tokenizer.rxWhitespace.ReplaceAllString(sentence, " ")
	sentence = strings.TrimSpace(sentence)

	return strings.Fields(sentence)
}
