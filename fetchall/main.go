//go:build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	urls := os.Args[1:]
	wg := sync.WaitGroup{}

	currentTime := time.Now()
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)

			if err != nil {
				fmt.Println(err)
				return
			}

			defer resp.Body.Close()
			deltaTime := time.Since(currentTime)
			fmt.Printf("url: %s, time: %s, ContentLength: %d\n", url, deltaTime.String(), resp.ContentLength)

		}(url)
	}

	wg.Wait()
	fmt.Printf("total time: %s", time.Since(currentTime).String())
}
