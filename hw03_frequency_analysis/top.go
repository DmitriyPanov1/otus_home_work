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

	m := make(map[string]uint)
	var word strings.Builder

	for _, r := range s {
		r = unicode.ToLower(r)

		if unicode.IsSpace(r) || r == '.' {
			word, m = addWordToMap(word, m)
		} else {
			word.WriteRune(r)
		}
	}

	_, m = addWordToMap(word, m)
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

func addWordToMap(word strings.Builder, m map[string]uint) (strings.Builder, map[string]uint) {
	if word.Len() > 0 {
		w := word.String()
		if w != "-" {
			m[w]++
		}
		word.Reset()
	}

	return word, m
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
