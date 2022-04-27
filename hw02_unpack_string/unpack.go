package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func isNextDigit(str []rune, currentPos int) bool {
	if currentPos != len(str)-1 {
		if unicode.IsDigit(str[currentPos+1]) {
			return true
		}
	}
	return false
}

func Unpack(s string) (string, error) {
	var sBuilder strings.Builder
	str := []rune(s)
	symbol := ""

	for i := 0; i < len(str); i++ {
		if !unicode.IsDigit(str[i]) {
			if symbol != "" {
				sBuilder.WriteString(symbol)
			}
			symbol = string(str[i])
		} else {
			if i == 0 {
				return "", ErrInvalidString
			}
			if next := isNextDigit(str, i); next {
				return "", ErrInvalidString
			}
			y, _ := strconv.Atoi(string(str[i]))
			sBuilder.WriteString(strings.Repeat(symbol, y))
			symbol = ""
		}
	}
	sBuilder.WriteString(symbol)

	return sBuilder.String(), nil
}
