package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	words := strings.Fields(text)

	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}

	type wordCount struct {
		word  string
		count int
	}

	var wordCounts []wordCount
	for word, count := range freq {
		wordCounts = append(wordCounts, wordCount{word, count})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].count == wordCounts[j].count {
			return wordCounts[i].word < wordCounts[j].word
		}
		return wordCounts[i].count > wordCounts[j].count
	})

	result := make([]string, 0, 10)
	for i := 0; i < len(wordCounts) && i < 10; i++ {
		result = append(result, wordCounts[i].word)
	}

	return result
}
