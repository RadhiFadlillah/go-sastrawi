package sastrawi

import (
	"regexp"
	"strings"
)

const (
	vowel     = "aiueo"
	consonant = "bcdfghjklmnpqrstvwxyz"
)

type char string

func newChar(word string, index int) char {
	if index >= len(word) {
		return char("")
	}

	return char(word[index])
}

func (c char) is(chars string) bool {
	for _, char := range strings.Split(chars, "") {
		if string(c) == char {
			return true
		}
	}

	return false
}

func (c char) isNot(chars string) bool {
	return !c.is(chars)
}

// Stemmer is object for stemming word
type Stemmer struct {
	dictionary    Dictionary
	rxPrefixFirst *regexp.Regexp
	rxParticle    *regexp.Regexp
	rxPossesive   *regexp.Regexp
	rxSuffix      *regexp.Regexp
}

// NewStemmer returns new Stemmer using dict as its dictionary
func NewStemmer(dict Dictionary) Stemmer {
	return Stemmer{
		dictionary:    dict,
		rxPrefixFirst: regexp.MustCompile(`^(be.+lah|be.+an|me.+i|di.+i|pe.+i|ter.+i)$`),
		rxParticle:    regexp.MustCompile(`-?(lah|kah|tah|pun)$`),
		rxPossesive:   regexp.MustCompile(`-?(ku|mu|nya)$`),
		rxSuffix:      regexp.MustCompile(`-?(is|isme|isasi|i|kan|an)$`),
	}
}

// Stem reduces inflected or derived word to its root form
func (stemmer Stemmer) Stem(word string) string {
	word = strings.ToLower(word)

	var (
		rootFound    = false
		originalWord = word
		particle     string
		possesive    string
		suffix       string
	)

	if len(word) < 3 {
		return word
	}

	if stemmer.dictionary.Find(word) {
		return word
	}

	// Check if prefix must be removed first
	if stemmer.rxPrefixFirst.MatchString(word) {
		// Remove prefix
		rootFound, word = stemmer.removePrefixes(word)
		if rootFound {
			return word
		}

		// Remove particle
		particle, word = stemmer.removeParticle(word)
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove possesive
		possesive, word = stemmer.removePossesive(word)
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove suffix
		suffix, word = stemmer.removeSuffix(word)
		if stemmer.dictionary.Find(word) {
			return word
		}
	} else {
		// Remove particle
		particle, word = stemmer.removeParticle(word)
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove possesive
		possesive, word = stemmer.removePossesive(word)
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove suffix
		suffix, word = stemmer.removeSuffix(word)
		if stemmer.dictionary.Find(word) {
			return word
		}

		// Remove prefix
		rootFound, word = stemmer.removePrefixes(word)
		if rootFound {
			return word
		}
	}

	// If no root found, do loopPengembalianAkhiran
	removedSuffixes := []string{"", suffix, possesive, particle}
	if suffix == "kan" {
		removedSuffixes = []string{"", "k", "an", possesive, particle}
	}

	rootFound, word = stemmer.loopPengembalianAkhiran(originalWord, removedSuffixes)
	if rootFound {
		return word
	}

	// When EVERYTHING failed, return original word
	return originalWord
}

func (stemmer Stemmer) removeParticle(word string) (string, string) {
	result := stemmer.rxParticle.ReplaceAllString(word, "")
	particle := strings.Replace(word, result, "", 1)
	return particle, result
}

func (stemmer Stemmer) removePossesive(word string) (string, string) {
	result := stemmer.rxPossesive.ReplaceAllString(word, "")
	possesive := strings.Replace(word, result, "", 1)
	return possesive, result
}

func (stemmer Stemmer) removeSuffix(word string) (string, string) {
	result := stemmer.rxSuffix.ReplaceAllString(word, "")
	suffix := strings.Replace(word, result, "", 1)
	return suffix, result
}

func (stemmer Stemmer) loopPengembalianAkhiran(originalWord string, suffixes []string) (bool, string) {
	lenSuffixes := 0
	for _, suffix := range suffixes {
		lenSuffixes += len(suffix)
	}
	wordWithoutSuffix := originalWord[:len(originalWord)-lenSuffixes]

	for i := range suffixes {
		suffixCombination := ""
		for j := 0; j <= i; j++ {
			suffixCombination += suffixes[j]
		}

		word := wordWithoutSuffix + suffixCombination
		if stemmer.dictionary.Find(word) {
			return true, word
		}

		rootFound, word := stemmer.removePrefixes(word)
		if rootFound {
			return true, word
		}
	}

	return false, originalWord
}

func (stemmer Stemmer) removePrefixes(word string) (bool, string) {
	originalWord := word
	currentPrefix := ""
	removedPrefix := ""
	recodingChar := []string{}

	for i := 0; i < 3; i++ {
		if len(word) < 3 {
			return false, originalWord
		}

		currentPrefix = word[:2]
		if currentPrefix == removedPrefix {
			break
		}

		removedPrefix, word, recodingChar = stemmer.removePrefix(word)
		if stemmer.dictionary.Find(word) {
			return true, word
		}

		for _, char := range recodingChar {
			if stemmer.dictionary.Find(char + word) {
				return true, char + word
			}
		}
	}

	return false, word
}

func (stemmer Stemmer) removePrefix(word string) (string, string, []string) {
	var (
		prefix   string
		result   string
		recoding []string
	)

	if strings.HasPrefix(word, "di") || strings.HasPrefix(word, "ke") || strings.HasPrefix(word, "se") || strings.HasPrefix(word, "ku") {
		prefix = word[:2]
		result = word[2:]
	} else if strings.HasPrefix(word, "kau") {
		prefix = "kau"
		result = word[3:]
	} else if strings.HasPrefix(word, "me") {
		prefix = "me"
		result, recoding = stemmer.removeMePrefix(word)
	} else if strings.HasPrefix(word, "pe") {
		prefix = "pe"
		result, recoding = stemmer.removePePrefix(word)
	} else if strings.HasPrefix(word, "be") {
		prefix = "be"
		result, recoding = stemmer.removeBePrefix(word)
	} else if strings.HasPrefix(word, "te") {
		prefix = "te"
		result, recoding = stemmer.removeTePrefix(word)
	} else {
		result, recoding = stemmer.removeInfix(word)
	}

	return prefix, result, recoding
}

func (stemmer Stemmer) removeMePrefix(word string) (string, []string) {
	s3 := newChar(word, 2)
	s4 := newChar(word, 3)
	s5 := newChar(word, 4)

	// Pattern 01
	// me{l|r|w|y}V => me-{l|r|w|y}V
	if s3.is("lrwy") && s4.is(vowel) {
		return word[2:], nil
	}

	// Pattern 02
	// mem{b|f|v} => mem-{b|f|v}
	if s3.is("m") && s4.is("bfv") {
		return word[3:], nil
	}

	// Pattern 03
	// mempe => mem-pe
	if s3.is("m") && s4.is("p") && s5.is("e") {
		return word[3:], nil
	}

	// Pattern 04
	// mem{rV|V} => mem-{rV|V} OR me-p{rV|V}
	if s3.is("m") && (s4.is(vowel) || (s4.is("r") && s5.is(vowel))) {
		return word[3:], []string{"m", "p"}
	}

	// Pattern 05
	// men{c|d|j|s|t|z} => men-{c|d|j|s|t|z}
	if s3.is("n") && s4.is("cdjstz") {
		return word[3:], nil
	}

	// Pattern 06
	// menV => nV OR tV
	if s3.is("n") && s4.is(vowel) {
		return word[3:], []string{"n", "t"}
	}

	// Pattern 07
	// meng{g|h|q|k} => meng-{g|h|q|k}
	if s3.is("n") && s4.is("g") && s5.is("ghqk") {
		return word[4:], nil
	}

	// Pattern 08
	// mengV => meng-V OR meng-kV OR me-ngV OR mengV- where V = 'e'
	if s3.is("n") && s4.is("g") && s5.is(vowel) {
		if s5.is("e") {
			return word[5:], nil
		}

		return word[4:], []string{"ng", "k"}
	}

	// Pattern 09
	// menyV => meny-sV OR me-nyV to stem menyala
	if s3.is("n") && s4.is("y") && s5.is(vowel) {
		if s5.is("a") {
			return word[2:], nil
		}

		return "s" + word[4:], nil
	}

	// Pattern 10
	// mempV => mem-pV where V != 'e'
	if s3.is("m") && s4.is("p") && s5.isNot("e") {
		return word[3:], nil
	}

	return word, nil
}

func (stemmer Stemmer) removePePrefix(word string) (string, []string) {
	s3 := newChar(word, 2)
	s4 := newChar(word, 3)
	s5 := newChar(word, 4)
	s6 := newChar(word, 5)
	s7 := newChar(word, 6)
	s8 := newChar(word, 7)

	// Pattern 01
	// pe{w|y}V => pe-{w|y}V
	if s3.is("wy") && s4.is(vowel) {
		return word[2:], nil
	}

	// Pattern 02
	// perV => per-V OR pe-rV
	if s3.is("r") && s4.is(vowel) {
		return word[3:], []string{"r"}
	}

	// Pattern 03
	// perCAP => per-CAP where C != 'r' and P != 'er'
	if s3.is("r") && s4.is(consonant) && s4.isNot("r") && s5.isNot("") && s6.isNot("e") {
		return word[3:], nil
	}

	// Pattern 4
	// perCAerV => per-CAerV where C != 'r'
	if s3.is("r") && s4.is(consonant) && s4.isNot("r") && s5.isNot("") && s6.is("e") && s7.is("r") && s8.is(vowel) {
		return word[3:], nil
	}

	// Pattern 05
	// pem{b|f|v} => pem-{b|f|v}
	if s3.is("m") && s4.is("bfv") {
		return word[3:], nil
	}

	// Pattern 06
	// pem{rV|V} => pe-m{rV|V} OR pe-p{rV|V}
	if s3.is("m") && (s4.is(vowel) || (s4.is("r") && s5.is(vowel))) {
		return word[3:], []string{"m", "p"}
	}

	// Pattern 07
	// pen{c|d|j|z} => pen-{c|d|j|z}
	if s3.is("n") && s4.is("cdjz") {
		return word[3:], nil
	}

	// Pattern 08
	// penV => pe-nV OR pe-tV
	if s3.is("n") && s4.is(vowel) {
		return word[3:], []string{"n", "t"}
	}

	// Pattern 09
	// pengC => peng-C
	if s3.is("n") && s4.is("g") && s5.is(consonant) {
		return word[4:], nil
	}

	// Pattern 10
	// pengV => peng-V OR peng-kV OR pengV- where V = 'e'
	if s3.is("n") && s4.is("g") && s5.is(vowel) {
		if s5.is("e") {
			return word[5:], nil
		}

		return word[4:], []string{"k"}
	}

	// Pattern 11
	// penyV => peny-sV OR pe-nyV
	if s3.is("n") && s4.is("y") && s5.is(vowel) {
		return word[4:], []string{"s", "ny"}
	}

	// Pattern 12
	// pelV => pe-lV OR pel-V for pelajar
	if s3.is("l") && s4.is(vowel) {
		if word == "pelajar" {
			return "ajar", nil
		}

		return word[2:], nil
	}

	// Pattern 13
	// peCerV => per-erV where C != {r|w|y|l|m|n}
	if s3.is(consonant) && s3.isNot("rwylmn") && s4.is("e") && s5.is("r") && s6.is(vowel) {
		return word[3:], nil
	}

	// Pattern 14
	// peCP => pe-CP where C != {r|w|y|l|m|n} and P != 'er'
	if s3.is(consonant) && s3.isNot("rwylmn") && s4.isNot("e") {
		return word[2:], nil
	}

	// Pattern 15
	// peC1erC2 => pe-C1erC2 where C1 != {r|w|y|l|m|n}
	if s3.is(consonant) && s3.isNot("rwylmn") && s4.is("e") && s5.is("r") && s6.is(consonant) {
		return word[2:], nil
	}

	return word, nil
}

func (stemmer Stemmer) removeBePrefix(word string) (string, []string) {
	s3 := newChar(word, 2)
	s4 := newChar(word, 3)
	s5 := newChar(word, 4)
	s6 := newChar(word, 5)
	s7 := newChar(word, 6)
	s8 := newChar(word, 7)

	// Pattern 01
	// berV => ber-V OR be-rV
	if s3.is("r") && s4.is(vowel) {
		return word[3:], []string{"r"}
	}

	// Pattern 02
	// berCAP => ber-CAP
	if s3.is("r") && s4.is(consonant) && s4.isNot("r") && s5.isNot("") && s6.isNot("e") {
		return word[3:], nil
	}

	// Pattern 3
	// berCAerV => ber-CAerV where C != 'r'
	if s3.is("r") && s4.is(consonant) && s4.isNot("r") && s5.isNot("") && s6.is("e") && s7.is("r") && s8.is(vowel) {
		return word[3:], nil
	}

	// Pattern 04
	// belajar => bel-ajar
	if word == "belajar" {
		return word[3:], nil
	}

	// Pattern 5
	// beC1erC2 => be-C1erC2 where C1 != {'r'|'l'}
	if s3.is(consonant) && s3.isNot("r") && s3.isNot("l") && s4.is("e") && s5.is("r") && s6.is(consonant) {
		return word[2:], nil
	}

	return word, nil
}

func (stemmer Stemmer) removeTePrefix(word string) (string, []string) {
	s3 := newChar(word, 2)
	s4 := newChar(word, 3)
	s5 := newChar(word, 4)
	s6 := newChar(word, 5)
	s7 := newChar(word, 6)

	// Pattern 01
	// terV => ter-V OR te-rV
	if s3.is("r") && s4.is(vowel) {
		return word[3:], []string{"r"}
	}

	// Pattern 02
	// terCerV => ter-CerV where C != 'r'
	if s3.is("r") && s4.is(consonant) && s4.isNot("r") && s5.is("e") && s6.is("r") && s7.is(vowel) {
		return word[3:], nil
	}

	// Pattern 3
	// terCP => ter-CP where C != 'r' and P != 'er'
	if s3.is("r") && s4.is(consonant) && s4.isNot("r") && s5.isNot("e") {
		return word[3:], nil
	}

	// Pattern 04
	// teC1erC2 => te-C1erC2 where C1 != 'r'
	if s3.is(consonant) && s3.isNot("r") && s4.is("e") && s5.is("r") && s6.is(consonant) {
		return word[2:], nil
	}

	// Pattern 05
	// terC1erC2 => ter-C1erC2 where C1 != 'r'
	if s3.is("r") && s4.is(consonant) && s4.isNot("r") && s5.is("e") && s6.is("r") && s7.is(consonant) {
		return word[3:], nil
	}

	return word, nil
}

func (stemmer Stemmer) removeInfix(word string) (string, []string) {
	s1 := newChar(word, 0)
	s2 := newChar(word, 1)
	s3 := newChar(word, 2)
	s4 := newChar(word, 3)

	// Pattern 01
	// CerV => CerV OR CV
	if s1.is(consonant) && s2.is("e") && s3.is("rlm") && s4.is(vowel) {
		return word[3:], []string{word[:3], word[:1]}
	}

	// Pattern 02
	// CinV => CinV OR CV
	if s1.is(consonant) && s2.is("i") && s3.is("n") && s4.is(vowel) {
		return word[3:], []string{word[:3], word[:1]}
	}

	return word, nil
}
