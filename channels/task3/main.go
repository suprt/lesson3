package main

import (
	"fmt"
	"math/rand"
)

func main() {
	naturals := make(chan int)
	squares := make(chan int)
	v := 10
	go func() {
		defer close(naturals)
		for range v {

			naturals <- rand.Intn(100)
		}
	}()

	go func() {
		defer close(squares)
		for n := range naturals {
			squares <- n * n

		}
	}()

	for s := range squares {
		fmt.Println(s)
	}

}
