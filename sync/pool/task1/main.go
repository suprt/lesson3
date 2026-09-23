package main

import (
	"fmt"
	"sync"
)

/*
Сделал несколько реализаций ProccessString - тестировал разницу в количестве аллокаций
*/

var bufPool = sync.Pool{
	New: func() any {
		return new([]byte)
	},
}

func ProcessString(s string) string {
	buf := bufPool.Get().(*[]byte)
	*buf = (*buf)[:0]
	for _, v := range s {
		if v >= 'a' && v <= 'z' {
			v = v - ('a' - 'A')
		}
		*buf = append(*buf, byte(v))
	}
	res := string(*buf)

	bufPool.Put(buf)

	return res
}

var bufPoolSecond = sync.Pool{
	New: func() any {
		return []byte{}
	},
}

func ProcessStringPoolValue(s string) string {
	buf := bufPoolSecond.Get().([]byte)
	buf = (buf)[:0]
	for _, v := range s {
		if v >= 'a' && v <= 'z' {
			v = v - ('a' - 'A')
		}
		buf = append(buf, byte(v))
	}
	res := string(buf)

	bufPoolSecond.Put(buf)

	return res
}

func AltProcessString(s string) string {
	buf := []byte(s)
	for i, v := range s {
		if v >= 'a' && v <= 'z' {
			buf[i] -= 'a' - 'A'
		}

	}
	return string(buf)

}

func main() {
	examples := []string{
		"hello, world!",
		"gopher",
		"lorem ipsum dolor sit amet",
	}

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}

}
