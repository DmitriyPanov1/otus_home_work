package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var prev rune
	var builder strings.Builder
	var shield bool
	var isDigitIgnorPrev bool

	for _, r := range s {
		switch {
		case string(prev) == "\\" && string(r) == "n":
			return "", ErrInvalidString
		case string(r) == "\\":
			if shield {
				builder.WriteString("\\")
				shield = false
			} else {
				shield = true
			}
		case unicode.IsDigit(r) && unicode.IsDigit(prev) && !isDigitIgnorPrev:
			return "", ErrInvalidString
		case unicode.IsDigit(r):
			if prev == 0 {
				return "", ErrInvalidString
			}

			if shield {
				builder.WriteString(string(r))
				shield = false
				isDigitIgnorPrev = true
			} else {
				var res string

				count, _ := strconv.Atoi(string(r))

				if count > 0 {
					res = strings.Repeat(string(prev), count-1)
				} else {
					res = builder.String()
					res = res[:len(res)-1]
					builder.Reset()
				}

				builder.WriteString(res)

				isDigitIgnorPrev = false
			}
		default:
			builder.WriteString(string(r))
			shield = false
		}

		prev = r
	}

	return builder.String(), nil
}
