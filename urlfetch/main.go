//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func panic(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	urls := os.Args[1:]
	for _, url := range urls {
		resp, err := http.Get(url)
		panic(err)

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		panic(err)

		fmt.Println(string(body))
	}
}
