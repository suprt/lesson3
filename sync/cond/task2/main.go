package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Restaurant struct {
	mu       sync.Mutex
	cond     *sync.Cond
	capacity int
	occupied int
}

func NewRestaurant(capacity int) *Restaurant {
	r := &Restaurant{
		capacity: capacity,
	}
	r.cond = sync.NewCond(&r.mu)
	return r
}

func (r *Restaurant) OccupyTable() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for r.occupied >= r.capacity {
		r.cond.Wait()
	}
	r.occupied++
	fmt.Println("occupied: ", r.occupied)
}

func (r *Restaurant) ReleaseTable() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.occupied--
	fmt.Println("Table released, occupied: ", r.occupied)
	r.cond.Signal()
}

func main() {
	r := NewRestaurant(5)
	var wg sync.WaitGroup

	for i := 1; i <= 15; i++ {
		wg.Add(1)
		//time.Sleep(2 * time.Second)
		go func(id int) {

			defer wg.Done()
			fmt.Printf("Guest %d is waiting for a table\n", id)
			r.OccupyTable()
			fmt.Printf("Guest %d occupied a table\n", id)
			time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
			fmt.Printf("Guest %d is leaving\n", id)
			r.ReleaseTable()
		}(i)

	}
	wg.Wait()

}
