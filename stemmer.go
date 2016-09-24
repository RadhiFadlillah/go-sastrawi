package sastrawi

import (
	"regexp"
	"strings"
)

const (
	vowel     = "aiueo"
	consonant = "bcdfghjklmnpqrsuvwxyz"
)

type char string

func NewChar(word string, index int) char {
	if index >= len(word) {
		return char("")
	}

	return char(word[index])
}

func (c char) Is(chars string) bool {
	for _, char := range strings.Split(chars, "") {
		if string(c) == char {
			return true
		}
	}

	return false
}

func (c char) IsNot(chars string) bool {
	return !c.Is(chars)
}

type Stemmer struct {
	dictionary    Dictionary
	rxDisallowed  *regexp.Regexp
	rxPrefixFirst *regexp.Regexp
	rxParticle    *regexp.Regexp
	rxPossesive   *regexp.Regexp
	rxSuffix      *regexp.Regexp
}

func NewStemmer(dict Dictionary) *Stemmer {
	return &Stemmer{
		dictionary:    dict,
		rxDisallowed:  regexp.MustCompile(`^(ber.+i|di.+an|ke.+i|ke.+kan|me.+an|ter.+an|per.+an)$`),
		rxPrefixFirst: regexp.MustCompile(`^(be.+lah|be.+an|me.+i|di.+i|pe.+i|ter.+i)$`),
		rxParticle:    regexp.MustCompile(`-?(lah|kah|tah|pun)$`),
		rxPossesive:   regexp.MustCompile(`-?(ku|mu|nya)$`),
		rxSuffix:      regexp.MustCompile(`-?(is|isme|isasi|i|kan|an)$`),
	}
}

func (stemmer *Stemmer) Stem(word string) string {
	rootFound := false
	originalWord := word
	recodingChar := []string{}

	if len(word) < 3 {
		return word
	}

	if stemmer.dictionary.Find(word) {
		return word
	}

	// Check if prefix must be removed first
	if stemmer.rxPrefixFirst.MatchString(word) {
		// Remove prefix
		rootFound, word, recodingChar = stemmer.removePrefixes(word)
		if rootFound {
			return word
		}

		// Try recoding
		rootFound, word = stemmer.recodingPrefix(word, recodingChar)
		if rootFound {
			return word
		}

		// Remove particle
		word = stemmer.rxParticle.ReplaceAllString(word, "")
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove possesive
		word = stemmer.rxPossesive.ReplaceAllString(word, "")
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove suffix
		word = stemmer.rxSuffix.ReplaceAllString(word, "")
		if stemmer.dictionary.Find(word) {
			return word
		}
	} else {
		// Remove particle
		word = stemmer.rxParticle.ReplaceAllString(word, "")
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove possesive
		word = stemmer.rxPossesive.ReplaceAllString(word, "")
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove suffix
		word = stemmer.rxSuffix.ReplaceAllString(word, "")
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove prefix
		rootFound, word, recodingChar = stemmer.removePrefixes(word)
		if rootFound {
			return word
		}

		// Try recoding
		rootFound, word = stemmer.recodingPrefix(word, recodingChar)
		if rootFound {
			return word
		}
	}

	// When EVERYTHING failed, return original word
	return originalWord
}

func (stemmer *Stemmer) removePrefixes(word string) (bool, string, []string) {
	originalWord := word
	currentPrefix := ""
	removedPrefix := ""
	recodingChar := []string{}

	for i := 0; i < 3; i++ {
		if len(word) < 3 {
			return false, originalWord, nil
		}

		if stemmer.rxDisallowed.MatchString(word) {
			break
		}

		currentPrefix = word[:2]
		if currentPrefix == removedPrefix {
			break
		}

		removedPrefix, word, recodingChar = stemmer.removePrefix(word)
		if stemmer.dictionary.Find(word) {
			return true, word, recodingChar
		}
	}

	return false, word, recodingChar
}

func (stemmer *Stemmer) recodingPrefix(word string, recodingChar []string) (bool, string) {
	if recodingChar == nil {
		return false, word
	}

	for _, char := range recodingChar {
		if stemmer.dictionary.Find(char + word) {
			return true, char + word
		}
	}

	return false, word
}

func (stemmer *Stemmer) removePrefix(word string) (string, string, []string) {
	var (
		prefix   string   = word[:2]
		recoding []string = nil
		result   string
	)

	switch prefix {
	case "di":
		fallthrough
	case "ke":
		fallthrough
	case "se":
		fallthrough
	case "ku":
		result = word[2:]
	case "me":
		result, recoding = stemmer.removeMePrefix(word)
	case "pe":
		result, recoding = stemmer.removePePrefix(word)
	case "be":
		result, recoding = stemmer.removeBePrefix(word)
	case "te":
		result, recoding = stemmer.removeTePrefix(word)
	default:
		result = word
	}

	return prefix, result, recoding
}

func (stemmer *Stemmer) removeMePrefix(word string) (string, []string) {
	s3 := NewChar(word, 2)
	s4 := NewChar(word, 3)
	s5 := NewChar(word, 4)

	// Pattern 01
	// me[lrwy][aiueo] => [lrwy][aiueo]
	if s3.Is("lrwy") && s4.Is(vowel) {
		return word[2:], nil
	}

	// Pattern 02
	// mem[bfv] => [bfv]
	if s3.Is("m") && s4.Is("bfv") {
		return word[3:], nil
	}

	// Pattern 03
	// mempe => pe
	if s3.Is("m") && s4.Is("p") && s5.Is("e") {
		return word[3:], nil
	}

	// Pattern 04
	// mem(r?)[aiueo] => m(r?)[aiueo] OR p(r?)[aiueo]
	if s3.Is("m") && (s4.Is(vowel) || (s4.Is("r") && s5.Is(vowel))) {
		return word[3:], []string{"m", "p"}
	}

	// Pattern 05
	// men[cdjstz] => [cdjstz]
	if s3.Is("n") && s4.Is("cdjstz") {
		return word[3:], nil
	}

	// Pattern 06
	// men[aiueo] => n[aiueo] OR t[aiueo]
	if s3.Is("n") && s4.Is(vowel) {
		return word[3:], []string{"n", "t"}
	}

	// Pattern 07
	// meng[ghqk] => [ghqk]
	if s3.Is("n") && s4.Is("g") && s5.Is("ghqk") {
		return word[4:], nil
	}

	// Pattern 08
	// meng[aiueo] => [aiueo]
	if s3.Is("n") && s4.Is("g") && s5.Is(vowel) {
		return word[4:], []string{"k"}
	}

	// Pattern 09
	// meny[aiueo] => s[aiueo]
	if s3.Is("n") && s4.Is("y") && s5.Is(vowel) {
		return "s" + word[4:], nil
	}

	// Pattern 10
	// memp[aiuo] => p[aiuo]
	if s3.Is("m") && s4.Is("p") && s5.Is("e") {
		return word[3:], nil
	}

	return word, nil
}

func (stemmer *Stemmer) removePePrefix(word string) (string, []string) {
	s3 := NewChar(word, 2)
	s4 := NewChar(word, 3)
	s5 := NewChar(word, 4)
	s6 := NewChar(word, 5)
	s7 := NewChar(word, 6)
	s8 := NewChar(word, 7)

	// Pattern 01
	// pe[wy][aiueo] => [wy][aiueo]
	if s3.Is("wy") && s4.Is(vowel) {
		return word[2:], nil
	}

	// Pattern 02
	// per[aiueo] => [aiueo] OR r[aiueo]
	if s3.Is("r") && s4.Is(vowel) {
		return word[3:], []string{"r"}
	}

	// Pattern 03
	// per[^aiueor][a-z][^e] => [^aiueor][a-z][^e]
	if s3.Is("r") && s4.Is(consonant) && s4.IsNot("r") && s5.IsNot("") && s6.IsNot("e") {
		return word[3:], nil
	}

	// Pattern 4
	// per[^aiueor][a-z]er[aiueo] => [^aiueor][a-z]er[aiueo]
	if s3.Is("r") && s4.Is(consonant) && s4.IsNot("r") && s5.IsNot("") && s6.Is("e") && s7.Is("r") && s8.Is(vowel) {
		return word[3:], nil
	}

	// Pattern 05
	// pem[bfv] => [bfv]
	if s3.Is("m") && s4.Is("bfv") {
		return word[3:], nil
	}

	// Pattern 06
	// pem(r?)[aiueo] => m(r?)[aiueo] OR p(r?)[aiueo]
	if s3.Is("m") && (s4.Is(vowel) || (s4.Is("r") && s5.Is(vowel))) {
		return word[3:], []string{"m", "p"}
	}

	// Pattern 07
	// pen[cdjz] => [cdjz]
	if s3.Is("n") && s4.Is("cdjz") {
		return word[3:], nil
	}

	// Pattern 08
	// pen[aiueo] => n[aiueo] OR t[aiueo]
	if s3.Is("n") && s4.Is(vowel) {
		return word[3:], []string{"n", "t"}
	}

	// Pattern 09
	// peng[^aiueo] => [^aiueo]
	if s3.Is("n") && s4.Is("g") && s5.Is(consonant) {
		return word[4:], nil
	}

	// Pattern 10
	// peng[aiueo] => [aiueo] OR k[aiueo]
	if s3.Is("n") && s4.Is("g") && s5.Is(vowel) {
		if s5.Is("e") {
			return word[5:], nil
		}

		return word[4:], []string{"k"}
	}

	// Pattern 11
	// peny[aiueo] => s[aiueo]
	if s3.Is("n") && s4.Is("y") && s5.Is(vowel) {
		return "s" + word[4:], nil
	}

	// Pattern 12
	// pel[aiueo] => l[aiueo] OR 'ajar' for 'pelajar'
	if s3.Is("l") && s4.Is(vowel) {
		if word == "pelajar" {
			return "ajar", nil
		}

		return word[2:], nil
	}

	// Pattern 13
	// pe[^aiueorwylmn]er[aiueo] => er[aiueo]
	if s3.Is(consonant) && s3.IsNot("rwylmn") && s4.Is("e") && s5.Is("r") && s6.Is(vowel) {
		return word[3:], nil
	}

	// Pattern 14
	// pe[^aiueorwylmn][^e] => [^aiueorwylmn][^e]
	if s3.Is(consonant) && s3.IsNot("rwylmn") && s4.IsNot("e") {
		return word[2:], nil
	}

	// Pattern 15
	// pe[^aiueorwylmn]er[^aiueo] => [^aiueorwylmn]er[^aiueo]
	if s3.Is(consonant) && s3.IsNot("rwylmn") && s4.Is("e") && s5.Is("r") && s6.Is(consonant) {
		return word[2:], nil
	}

	return word, nil
}

func (stemmer *Stemmer) removeBePrefix(word string) (string, []string) {
	s3 := NewChar(word, 2)
	s4 := NewChar(word, 3)
	s5 := NewChar(word, 4)
	s6 := NewChar(word, 5)
	s7 := NewChar(word, 6)
	s8 := NewChar(word, 7)

	// Pattern 01
	// ber[aiueo] => [aiueo] OR r[aiueo]
	if s3.Is("r") && s4.Is(vowel) {
		return word[3:], []string{"r"}
	}

	// Pattern 02
	// ber[^aiueor][a-z][^e] => [^aiueor][a-z][^e]
	if s3.Is("r") && s4.Is(consonant) && s4.IsNot("r") && s5.IsNot("") && s6.IsNot("e") {
		return word[3:], nil
	}

	// Pattern 3
	// ber[^aiueor][a-z]er[aiueo] => [^aiueor][a-z]er[aiueo]
	if s3.Is("r") && s4.Is(consonant) && s4.IsNot("r") && s5.IsNot("") && s6.Is("e") && s7.Is("r") && s8.Is(vowel) {
		return word[3:], nil
	}

	// Pattern 04
	// belajar => ajar
	if word == "belajar" {
		return word[3:], nil
	}

	// Pattern 5
	// be[^aiueorl]er[aiueo] => [^aiueorl]er[aiueo]
	if s3.Is(consonant) && s3.IsNot("r") && s3.IsNot("l") && s4.Is("e") && s5.Is("r") && s6.Is(vowel) {
		return word[2:], nil
	}

	return word, nil
}

func (stemmer *Stemmer) removeTePrefix(word string) (string, []string) {
	s3 := NewChar(word, 2)
	s4 := NewChar(word, 3)
	s5 := NewChar(word, 4)
	s6 := NewChar(word, 5)
	s7 := NewChar(word, 6)

	// Pattern 01
	// ter[aiueo] => [aiueo] OR r[aiueo]
	if s3.Is("r") && s4.Is(vowel) {
		return word[3:], []string{"r"}
	}

	// Pattern 02
	// ter[^aiueor]er[aiueo] => [^aiueor]er[aiueo]
	if s3.Is("r") && s4.Is(consonant) && s4.IsNot("r") && s5.Is("e") && s6.Is("r") && s7.Is(vowel) {
		return word[3:], nil
	}

	// Pattern 3
	// ter[^aiueor][^e] => [^aiueor][^e]
	if s3.Is("r") && s4.Is(consonant) && s4.IsNot("r") && s5.IsNot("e") {
		return word[3:], nil
	}

	// Pattern 04
	// te[^aiueor]er[^aiueo] => [^aiueor]er[^aiueo]
	if s3.Is(consonant) && s3.IsNot("r") && s4.Is("e") && s5.Is("r") && s6.Is(consonant) {
		return word[2:], nil
	}

	// Pattern 05
	// ter[^aiueor]er[^aiueo] => [^aiueor]er[^aiueo]
	if s3.Is("r") && s4.Is(consonant) && s4.IsNot("r") && s5.Is("e") && s6.Is("r") && s7.Is(consonant) {
		return word[3:], nil
	}

	return word, nil
}
