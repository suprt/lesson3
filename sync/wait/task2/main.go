package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func fetchUrl(url string) error {
	_, err := http.Get(url)
	return err
}
func main() {
	urls := []string{
		"https://www.lamoda.ru",
		"https://www.yandex.ru",
		"https://www.mail.ru",
		"https://www.google.ru",
	}
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fmt.Printf("Fetching %s....\n", url)
			err := fetchUrl(url)
			if err != nil {
				fmt.Printf("Error feaching %s: %v\n", url, err)
				return
			}
			fmt.Printf("Fetched %s\n", url)
		}(url)
	}

	fmt.Println("All request launched!")
	wg.Wait()
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Program finished")
}
