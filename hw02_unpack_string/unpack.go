package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func isNextNotShieldedDigit(str []rune, currentPos int) bool {
	if currentPos != len(str)-1 {
		if unicode.IsDigit(str[currentPos+1]) && string(str[currentPos-1]) != "\\" {
			return true
		}
	}
	return false
}

func isValidShielding(str []rune, currentPos int) bool {
	if currentPos != len(str)-1 {
		if !unicode.IsDigit(str[currentPos+1]) && string(str[currentPos+1]) != "\\" {
			return false
		}
	}
	return true
}

func Unpack(s string) (string, error) {
	var sBuilder strings.Builder
	str := []rune(s)
	symbol := ""
	shielding := false

	for i := 0; i < len(str); i++ {
		switch {
		case unicode.IsDigit(str[i]):
			if i == 0 {
				return "", ErrInvalidString
			}
			if next := isNextNotShieldedDigit(str, i); next {
				return "", ErrInvalidString
			}
			if shielding {
				symbol = string(str[i])
				shielding = false
			} else {
				y, _ := strconv.Atoi(string(str[i]))
				sBuilder.WriteString(strings.Repeat(symbol, y))
				symbol = ""
			}
		case string(str[i]) == "\\":
			if i == len(str)-1 {
				return "", ErrInvalidString
			}
			if next := isValidShielding(str, i); !next {
				return "", ErrInvalidString
			}
			if symbol != "" {
				sBuilder.WriteString(symbol)
				symbol = ""
			}
			if shielding {
				symbol = string(str[i])
				shielding = false
			} else {
				shielding = true
			}
		default:
			if symbol != "" {
				sBuilder.WriteString(symbol)
			}
			symbol = string(str[i])
		}
	}
	sBuilder.WriteString(symbol)

	return sBuilder.String(), nil
}
