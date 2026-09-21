package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {

		for {
			time.Sleep(time.Duration(rand.Intn(400)+100) * time.Millisecond)
			naturals <- rand.Intn(100)
		}
	}()

	go func() {

		for {

			n := <-naturals
			squares <- n * n

		}
	}()
	go func() {

		for {

			n := <-squares
			fmt.Println(n)

		}
	}()
	select {}
}
