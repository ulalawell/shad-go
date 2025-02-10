//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var b strings.Builder
	b.Grow(len(input))

	for len(input) > 0 {
		symbol, size := utf8.DecodeLastRuneInString(input)
		input = input[:len(input)-size]
		b.WriteRune(symbol)

	}

	return b.String()
}
