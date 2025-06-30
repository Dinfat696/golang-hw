package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var result strings.Builder
	runes := []rune(s)
	n := len(runes)

	for i := 0; i < n; i++ {
		ch := runes[i]

		if ch == '\\' {
			if i+1 >= n {
				return "", ErrInvalidString
			}
			result.WriteRune(runes[i+1])
			i++
			continue
		}

		if !unicode.IsDigit(ch) {
			result.WriteRune(ch)
			continue
		}

		if i == 0 || (i > 0 && unicode.IsDigit(runes[i-1])) {
			return "", ErrInvalidString
		}

		count := int(ch - '0')
		if count == 0 {
			if result.Len() == 0 {
				return "", ErrInvalidString
			}
			current := result.String()
			result.Reset()
			result.WriteString(current[:len(current)-1])
			continue
		}

		if i == 0 {
			return "", ErrInvalidString
		}
		prevChar := runes[i-1]
		for j := 1; j < count; j++ {
			result.WriteRune(prevChar)
		}
	}

	return result.String(), nil
}