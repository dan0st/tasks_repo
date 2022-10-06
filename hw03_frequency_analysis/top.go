package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`[^0-9a-z_а-я-]+`)

func Filter(word string, reg *regexp.Regexp) string {
	lowerWord := strings.ToLower(word)
	result := reg.ReplaceAll([]byte(lowerWord), []byte(""))
	return string(result)
}

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
		if word == "-" {
			continue
		}
		filteredWord := Filter(word, reg)
		words[filteredWord]++
	}

	return GetTop10(words)
}
