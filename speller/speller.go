//go:build !solution

package speller

import (
	"strings"
)

func Spell(number int64) string {
	nameByNumber := map[int64]string{
		0:          "zero",
		1:          "one",
		2:          "two",
		3:          "three",
		4:          "four",
		5:          "five",
		6:          "six",
		7:          "seven",
		8:          "eight",
		9:          "nine",
		10:         "ten",
		11:         "eleven",
		12:         "twelve",
		13:         "thirteen",
		14:         "fourteen",
		15:         "fifteen",
		16:         "sixteen",
		17:         "seventeen",
		18:         "eighteen",
		19:         "nineteen",
		20:         "twenty",
		30:         "thirty",
		40:         "forty",
		50:         "fifty",
		60:         "sixty",
		70:         "seventy",
		80:         "eighty",
		90:         "ninety",
		100:        "hundred",
		1000:       "thousand",
		1000000:    "million",
		1000000000: "billion",
	}

	absNumber := abs(number)
	if absNumber == 0 {
		return nameByNumber[0]
	}

	var sb strings.Builder

	groups := []struct {
		value int64
		name  string
	}{
		{absNumber / 1_000_000_000, nameByNumber[1_000_000_000]},  // billions
		{(absNumber / 1_000_000) % 1000, nameByNumber[1_000_000]}, // millions
		{(absNumber / 1000) % 1000, nameByNumber[1000]},           // thousands
		{absNumber % 1000, ""},                                    // last three digits
	}

	for _, group := range groups {
		if group.value != 0 {
			if sb.Len() > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(process3Digit(nameByNumber, group.value))
			if group.name != "" {
				sb.WriteByte(' ')
				sb.WriteString(group.name)
			}
		}
	}

	if number < 0 {
		return "minus " + sb.String()
	}
	return sb.String()
}

func process3Digit(nameByNumber map[int64]string, n int64) string {
	var sb strings.Builder

	if n >= 100 {
		sb.WriteString(nameByNumber[n/100])
		sb.WriteString(" hundred")
		n %= 100

		if n == 0 {
			return sb.String()
		}

		sb.WriteByte(' ')
	}

	if n > 0 && n < 20 {
		sb.WriteString(nameByNumber[n])
		return sb.String()
	}

	if n >= 20 {
		sb.WriteString(nameByNumber[n/10*10])
		if n%10 != 0 {
			sb.WriteByte('-')
			sb.WriteString(nameByNumber[n%10])
		}
	}

	return sb.String()
}

func abs(k int64) int64 {
	if k < 0 {
		return -k
	}
	return k
}
