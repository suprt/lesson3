package main

import (
	"fmt"
)

type ServerMetric struct {
	Name  string  // Название метрики (например, "memory_usage")
	Value float64 // Значение в байтах
}

func Transformer(in <-chan ServerMetric) <-chan ServerMetric {
	out := make(chan ServerMetric)
	go func() {
		defer close(out)
		for v := range in {
			out <- ServerMetric{Name: v.Name, Value: v.Value / 1024}
		}
	}()

	return out
}

func generator(name string, count int, out chan<- ServerMetric) {

	defer close(out)
	for i := 0; i < count; i++ {
		out <- ServerMetric{Name: name, Value: float64(1024 * (i + 1))}
	}

}
func main() {
	metrics := make(chan ServerMetric)

	transformed := Transformer(metrics)

	go generator("memory_usage", 10, metrics)

	for v := range transformed {
		fmt.Println(v)
	}

}
