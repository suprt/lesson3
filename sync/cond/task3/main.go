package main

import (
	"fmt"
	"sync"
	"time"
)

type ConnectionPool struct {
	mu        sync.Mutex
	cond      *sync.Cond
	available []*Connection
}
type Connection struct {
	ID   int
	addr string
}

func NewConnectionPool(limit int) *ConnectionPool {
	pool := &ConnectionPool{
		available: make([]*Connection, 0, limit),
	}
	for i := range limit {
		pool.available = append(pool.available, &Connection{
			ID:   i,
			addr: "localhost:5432",
		})
	}
	pool.cond = sync.NewCond(&pool.mu)
	return pool
}

func (cp *ConnectionPool) Get() *Connection {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for len(cp.available) == 0 {
		cp.cond.Wait()
	}

	conn := cp.available[len(cp.available)-1]
	cp.available = cp.available[:len(cp.available)-1]

	return conn
}
func (cp *ConnectionPool) Release(conn *Connection) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.available = append(cp.available, conn)

	cp.cond.Signal()
}

func main() {
	pool := NewConnectionPool(3) // Пул на 3 подключения

	for i := 0; i < 10; i++ {
		go func(id int) {
			conn := pool.Get()
			defer pool.Release(conn)

			fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)
			time.Sleep(2 * time.Second) // Имитация работы
			fmt.Printf("Горутина %d: подключение %d освобождено\n", id, conn.ID)
		}(i)
	}

	time.Sleep(10 * time.Second)
}
