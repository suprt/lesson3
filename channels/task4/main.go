package main

import (
	"fmt"
	"sync"
	"time"
)

func mergeChannels(channels ...<-chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}
	for _, channel := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range channel {
				out <- v
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
