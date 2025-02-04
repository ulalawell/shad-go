package main

import (
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	files := os.Args[1:]
	countByWord := map[string]int{}

	for _, file := range files {
		data, error := os.ReadFile(file)
		check(error)
		lines := strings.Split(string(data), "\n")

		for _, line := range lines {
			countByWord[line]++
		}
	}

	for word, count := range countByWord {
		if count > 1 {
			fmt.Printf("%d\t%s\n", count, word)
		}
	}
}
