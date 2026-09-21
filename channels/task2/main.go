package main

import (
	"sync"
	"time"
)

func worker() chan int {
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 42
	}()
	return ch
}
func main() {
	timeStart := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = <-worker()
		}()
	}
	wg.Wait()

	println(int(time.Since(timeStart).Seconds()))
}

/*
Будет выполняться порядка 6 секунд.
receive операции вычисляются последовательно слева направо, поэтому и worker() будут
запускаться последовательно, из-за чего итоговое время составит 3+3=6с
*/
