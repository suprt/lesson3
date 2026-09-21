package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Result struct {
	File  string
	Count int
	Err   error
}

func countWords(input string) int {
	return len(strings.Fields(input))
}
func worker(path string, out chan<- Result) {
	data, err := os.ReadFile(path)
	if err != nil {
		out <- Result{path, 0, err}
	}
	out <- Result{path, countWords(string(data)), nil}
}

/*func altWorker(path string) Result {
	data, err := os.ReadFile(path)
	if err != nil {
		return Result{
			File: path,
			Err:  err,
		}
	}

	return Result{
		File:  path,
		Count: countWords(string(data)),
	}
}*/

func SplitJobs(input <-chan string) <-chan Result {
	results := make(chan Result)
	wg := sync.WaitGroup{}
	for path := range input {
		wg.Add(1)

		go func() {
			defer wg.Done()
			worker(path, results)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func main() {
	start := time.Now()
	jobs := make(chan string)
	go func() {
		defer close(jobs)
		jobs <- "file1.txt"
		jobs <- "file2.txt"
		jobs <- "file3.txt"
	}()

	results := SplitJobs(jobs)
	total := 0
	for result := range results {
		if result.Err != nil {
			fmt.Printf("%s: %v\n", result.File, result.Err)
			continue
		}

		fmt.Printf("%s: %d words\n", result.File, result.Count)
		total += result.Count
	}

	fmt.Printf("Total: %d\n", total)
	fmt.Println(time.Since(start))
	/*
		start = time.Now()
		files := []string{
			"file1.txt",
			"file2.txt",
			"file3.txt",
		}

		total = 0

		for _, path := range files {
			result := altWorker(path)

			if result.Err != nil {
				fmt.Printf("%s: %v\n", result.File, result.Err)
				continue
			}

			fmt.Printf("%s: %d words\n", result.File, result.Count)
			total += result.Count
		}

		fmt.Println("Total:", total)
		fmt.Println(time.Since(start))*/
}
