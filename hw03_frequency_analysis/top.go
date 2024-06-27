package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type Word struct {
	Value string
	Count uint
}

func Top10(s string) []string {
	if s == "" {
		return nil
	}

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '-'
	}

	m := make(map[string]uint)

	for _, word := range strings.FieldsFunc(s, f) {
		word = strings.ToLower(word)

		if word != "-" {
			m[word]++
		}
	}

	words := make([]Word, 0, len(m))

	for k, v := range m {
		words = append(words, Word{k, v})
	}

	words = sortWords(words)

	result := make([]string, 0, 10)
	for i := 0; i < len(words) && i < 10; i++ {
		result = append(result, words[i].Value)
	}

	return result
}

func sortWords(words []Word) []Word {
	sort.Slice(words, func(i, j int) bool {
		if words[i].Count == words[j].Count {
			return words[i].Value < words[j].Value
		}
		return words[i].Count > words[j].Count
	})

	return words
}
