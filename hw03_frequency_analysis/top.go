package hw03frequencyanalysis
ipackage hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	if len(text) == 0 {
		return nil
	}

	words := strings.Fields(text)
	freq := make(map[string]int, len(words))

	for _, word := range words {
		freq[word]++
	}

	type wordCount struct {
		word  string
		count int
	}

	wordCounts := make([]wordCount, 0, len(freq))
	for word, count := range freq {
		wordCounts = append(wordCounts, wordCount{word, count})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].count == wordCounts[j].count {
			return wordCounts[i].word < wordCounts[j].word
		}
		return wordCounts[i].count > wordCounts[j].count
	})

	resultSize := 10
	if len(wordCounts) < resultSize {
		resultSize = len(wordCounts)
	}
	result := make([]string, 0, resultSize)
	for i := 0; i < resultSize; i++ {
		result = append(result, wordCounts[i].word)
	}

	return result
}

var taskWithAsteriskIsCompleted = false