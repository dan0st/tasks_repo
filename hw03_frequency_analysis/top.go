package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func GetTop10(words map[string]int) []string {
	type wordAmountStruct struct {
		amount int
		word   string
	}

	wordsSlice := make([]wordAmountStruct, 0, len(words))
	for word, amount := range words {
		wordsSlice = append(wordsSlice, wordAmountStruct{amount, word})
	}

	sort.Slice(wordsSlice, func(i, j int) bool {
		if wordsSlice[i].amount != wordsSlice[j].amount {
			return wordsSlice[i].amount > wordsSlice[j].amount
		}
		return wordsSlice[i].word < wordsSlice[j].word
	})

	var topWordsSlice []string
	for i := 0; i < 10 && i < len(wordsSlice); i++ {
		topWordsSlice = append(topWordsSlice, wordsSlice[i].word)
	}

	return topWordsSlice
}

func Top10(s string) []string {
	splittedStr := strings.Fields(s)
	words := make(map[string]int)

	for _, word := range splittedStr {
		words[word]++
	}

	return GetTop10(words)
}
