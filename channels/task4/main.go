package main

import (
	"fmt"
	"sync"
	"time"
)

func mergeChannels(channels ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, channel := range channels {
		go func() {
			defer wg.Done()
			for n := range channel {
				out <- n
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	go func() {
		defer close(a)
		a <- 5
		a <- 10
		time.Sleep(time.Millisecond * 500)
		a <- 3
	}()

	go func() {
		defer close(b)
		b <- 1
		b <- 2
		b <- 4
	}()

	go func() {
		defer close(c)
		c <- 11
		c <- 22
		c <- 44
	}()
	d := mergeChannels(a, b, c)
	for n := range d {
		fmt.Println(n)
	}
}
