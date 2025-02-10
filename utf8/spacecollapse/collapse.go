//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var sb strings.Builder
	symbolsToChange := map[rune]struct{}{
		'\r': {},
		'\n': {},
		'\t': {},
		' ':  {},
	}

	for i := 0; i < len(input); {
		needSpace := false
		for j := i; j < len(input); {
			symbol, size := utf8.DecodeRuneInString(input[i:])
			if _, ok := symbolsToChange[symbol]; !ok {
				break
			}
			needSpace = true
			j += size
			i += size
		}

		if needSpace {
			sb.WriteRune(' ')
		}

		if i < len(input) {
			symbol, size := utf8.DecodeRuneInString(input[i:])

			sb.WriteRune(symbol)
			i += size
		}

	}
	return sb.String()
}
