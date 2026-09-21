package main

import (
	"fmt"
	"sync"
	"time"
)

// Реплика БД (имитация)
func dbReplica(name string, in <-chan int) {
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(100 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

func Tee(in <-chan int, out []chan int) {

	for data := range in {
		for _, outChan := range out {
			outChan <- data
		}
	}
	for _, outChan := range out {
		close(outChan)
	}
}

func main() {
	input := make(chan int) // Канал для входящих данных
	replicas := []chan int{ // Реплики БД (каналы)
		make(chan int),
		make(chan int),
		make(chan int),
	}

	wg := &sync.WaitGroup{}
	go Tee(input, replicas)

	for i, outChan := range replicas {
		wg.Add(1)
		go func() {
			defer wg.Done()
			dbReplica(fmt.Sprintf("%d", i), outChan)
		}()
	}

	for i := range 5 {
		input <- i
	}
	close(input)
	wg.Wait()
}
