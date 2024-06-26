package hw03frequencyanalysis

import (
	"sort"
)

type Word struct {
	Value string
	Count uint
}

func Top10(s string) []string {
	var result []string
	m := make(map[string]uint)
	var word []rune

	if s == "" {
		return nil
	}

	for _, r := range s {
		if (string(r) == " " || string(r) == "\n" || string(r) == "\t") && word != nil {
			m = addWordToMap(m, string(word))
			word = nil
		}

		if string(r) != " " && string(r) != "\n" && string(r) != "\t" {
			word = append(word, r)
		}
	}

	m = addWordToMap(m, string(word))

	words := make([]Word, 0, 100)

	for k, v := range m {
		words = append(words, Word{k, v})
	}

	sort.Slice(words, func(i, j int) bool {
		if words[i].Count == words[j].Count {
			return words[i].Value < words[j].Value
		}

		return words[i].Count > words[j].Count
	})

	for k, w := range words {
		result = append(result, w.Value)

		if k >= 9 {
			break
		}
	}

	return result
}

func addWordToMap(m map[string]uint, w string) map[string]uint {
	_, ok := m[w]

	if ok {
		m[w]++
	} else {
		m[w] = 1
	}

	return m
}
