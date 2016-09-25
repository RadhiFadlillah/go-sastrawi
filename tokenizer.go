package sastrawi

import (
	"html"
	"regexp"
	"strings"
)

type Tokenizer struct {
	rxURL        *regexp.Regexp
	rxEmail      *regexp.Regexp
	rxTwitter    *regexp.Regexp
	rxEscapeStr  *regexp.Regexp
	rxSymbol     *regexp.Regexp
	rxWhitespace *regexp.Regexp
}

func NewTokenizer() Tokenizer {
	return Tokenizer{
		rxURL:        regexp.MustCompile(`(www\.|https?|s?ftp)\S+`),
		rxEmail:      regexp.MustCompile(`\S+@\S+`),
		rxTwitter:    regexp.MustCompile(`(@|#)\S+`),
		rxSymbol:     regexp.MustCompile(`[^a-z\s]`),
		rxWhitespace: regexp.MustCompile(`\s+`),
	}
}

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

	return strings.Split(sentence, " ")
}
