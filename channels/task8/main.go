package main

import (
	"fmt"
	"sync"
)

func Pars(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for data := range in {

			out <- fmt.Sprintf("parsed- %s", data)
		}

	}()
	return out
}

func Split(in <-chan string, n int) []<-chan string {
	if n <= 0 {
		return nil
	}
	out := make([]chan string, n)
	for i := range n {
		out[i] = make(chan string)
	}

	go func() {
		idx := 0
		for data := range in {
			out[idx] <- data
			idx = (idx + 1) % n
		}
		for _, ch := range out {
			close(ch)
		}
	}()

	res := make([]<-chan string, n)
	for i := range n {
		res[i] = out[i]
	}

	return res
}

func Send(in []<-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		wg := sync.WaitGroup{}
		for _, ch := range in {
			wg.Add(1)
			go func(ch <-chan string) {
				defer wg.Done()
				for data := range ch {
					out <- fmt.Sprintf("sent- %s", data)
				}
			}(ch)
		}
		wg.Wait()
	}()

	return out
}

func generator(name string, count int, out chan<- string) {

	defer close(out)
	for i := range count {
		out <- fmt.Sprintf("%s: %d", name, i)
	}
}

func main() {
	data := make(chan string)
	go generator("data", 10, data)

	parsed := Pars(data)
	split := Split(parsed, 5)
	sent := Send(split)

	for v := range sent {
		fmt.Println(v)
	}

}
