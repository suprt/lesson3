package main

import (
	"fmt"
	"sync"
	"time"
)

type Database struct {
	conn string
	s    sync.Once
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) GetConnection() string {
	db.s.Do(func() {
		fmt.Println("connecting to database")
		time.Sleep(500 * time.Millisecond)
		db.conn = "localhost:5276"
		fmt.Println("connection established")

	})
	return db.conn
}

func main() {
	db := NewDatabase()
	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			conn := db.GetConnection()
			fmt.Println(conn)
		}(i)
	}
	wg.Wait()
}
