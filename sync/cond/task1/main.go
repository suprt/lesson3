package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrShutdown = errors.New("queue is shut down")

type BoundedQueue struct {
	mu       sync.Mutex
	cond     *sync.Cond
	tasks    []interface{}
	capacity int
	shutdown bool
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{
		capacity: capacity,
		tasks:    make([]interface{}, 0, capacity),
	}
	q.cond = sync.NewCond(&q.mu)
	return q
}

func (q *BoundedQueue) Put(task interface{}) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for len(q.tasks) >= q.capacity && !q.shutdown {
		q.cond.Wait()
	}
	if q.shutdown {
		return ErrShutdown
	}
	q.tasks = append(q.tasks, task)
	q.cond.Signal()
	return nil
}

func (q *BoundedQueue) Get() (interface{}, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.tasks) == 0 && !q.shutdown {
		q.cond.Wait()
	}
	if q.shutdown {
		return nil, ErrShutdown
	}
	task := q.tasks[0]
	q.tasks[0] = nil
	q.tasks = q.tasks[1:]
	q.cond.Signal()
	return task, nil
}

func (q *BoundedQueue) Shutdown() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.shutdown = true
	q.cond.Broadcast()
}

func main() {
	q := NewBoundedQueue(2)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 5; i++ {
			fmt.Printf("put: %d\n", i)
			//time.Sleep(2 * time.Second)
			if err := q.Put(i); err != nil {
				if errors.Is(err, ErrShutdown) {
					fmt.Println("shutdown...")
					return
				}
			}
			fmt.Printf("put done: %d\n", i)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 1; i <= 5; i++ {
			time.Sleep(1 * time.Second)

			task, err := q.Get()
			if err != nil {

				if errors.Is(err, ErrShutdown) {
					fmt.Println("shutdown...")
					return
				}
			}
			fmt.Printf("get: %v\n", task)
		}
	}()

	//time.Sleep(1 * time.Second)
	//q.Shutdown()

	wg.Wait()

	fmt.Println("done")
}
