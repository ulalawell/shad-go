//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var sb strings.Builder
	sb.Grow(len(format) + len(args)*2)
	usedArgs := 0

	for i := 0; i < len(format); i++ {
		if format[i] != '{' {
			sb.WriteByte(format[i])
			continue
		}

		j := i + 1
		for j < len(format) && format[j] != '}' {
			j++
		}

		if j >= len(format) {
			sb.WriteByte('{')
			continue
		}

		if i+1 == j {
			if usedArgs < len(args) {
				writeArg(&sb, args[usedArgs])
			}
		} else {
			argNumber, err := strconv.Atoi(format[i+1 : j])
			if err == nil && argNumber < len(args) {
				writeArg(&sb, args[argNumber])
			}
		}

		usedArgs++
		i = j
	}

	return sb.String()
}

func writeArg(sb *strings.Builder, arg interface{}) {
	switch v := arg.(type) {
	case string:
		sb.WriteString(v)
	case int:
		sb.WriteString(strconv.Itoa(v))
	default:
		sb.WriteString(fmt.Sprint(v))
	}
}
