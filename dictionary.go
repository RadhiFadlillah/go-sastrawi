package sastrawi

import (
	"fmt"
)

// Dictionary is map[string]struct{} that used as root words database
type Dictionary map[string]struct{}

// NewDictionary creates new Dictionary with words as its content
func NewDictionary(words ...string) Dictionary {
	dict := make(map[string]struct{})
	for _, word := range words {
		dict[word] = struct{}{}
	}

	return Dictionary(dict)
}

// Count returns the size of dictionary
func (dictionary Dictionary) Count() int {
	return len(dictionary)
}

// Contains is used for to check if word exists within dictionary
func (dictionary Dictionary) Contains(word string) bool {
	_, found := dictionary[word]
	return found
}

// Add is used to append new words to dictionary
func (dictionary Dictionary) Add(words ...string) {
	for _, word := range words {
		dictionary[word] = struct{}{}
	}
}

// Remove is used to remove some words from dictionary
func (dictionary Dictionary) Remove(words ...string) {
	for _, word := range words {
		delete(dictionary, word)
	}
}

// Print is used for printing content of dictionary, where each word is separated by separator
func (dictionary Dictionary) Print(separator string) {
	if separator == "" {
		separator = ", "
	}

	index := 0
	lenDictionary := len(dictionary)
	for word := range dictionary {
		index++
		if index >= lenDictionary {
			fmt.Println(word)
		} else {
			fmt.Print(word, separator)
		}
	}
}
