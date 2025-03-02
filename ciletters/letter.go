//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed notification.tmpl
var tmplFile string // Встраиваем файл как строку

var countLinesGitlabRunner = 10
var countBytesHashCommit = 8
var indent = "\n            "

func MakeLetter(n *Notification) (string, error) {
	funcMap := template.FuncMap{
		"shortenHashCommit": func(s string) string {
			return s[:countBytesHashCommit]
		},
		"shortenLog": func(s string) string {
			lines := strings.Split(s, "\n")

			if len(lines) > countLinesGitlabRunner {
				return strings.Join(lines[len(lines)-countLinesGitlabRunner:], indent)
			}
			return strings.Join(lines, indent)
		},
	}

	var buf bytes.Buffer

	tmpl, err := template.New("notification").Funcs(funcMap).Parse(tmplFile)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, n)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
