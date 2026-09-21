package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	alreadyStored := make(map[int]struct{})
	capacity := 1000
	doubles := make([]int, 0, capacity)
	for i := 0; i < capacity; i++ {
		doubles = append(doubles, rand.Intn(10))
	}
	uniqueIDs := make(chan int, capacity)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := 0; i < capacity; i++ {
		i := i //в современных версиях не требуется
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			//конкурентный доступ к map - защитил через mutex
			if _, ok := alreadyStored[doubles[i]]; !ok {
				alreadyStored[doubles[i]] = struct{}{}
				uniqueIDs <- doubles[i]
			}

		}()
	}
	wg.Wait()
	close(uniqueIDs)
	//range по не закрытому каналу приведёт к deadlock - закрыл канал
	for val := range uniqueIDs {
		fmt.Println(val)
	}
	fmt.Println(uniqueIDs)
}
