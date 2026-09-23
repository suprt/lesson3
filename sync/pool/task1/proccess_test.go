package main

import (
	"strings"
	"testing"
)

var result string

func TestAllocations(t *testing.T) {
	s := strings.Repeat("hello world ", 1000)

	allocsPool := testing.AllocsPerRun(100, func() {
		for i := 0; i < 1000; i++ {
			result = ProcessString(s)
		}
	})

	allocsNoPool := testing.AllocsPerRun(100, func() {
		for i := 0; i < 1000; i++ {
			result = AltProcessString(s)
		}
	})

	allocsNotPointer := testing.AllocsPerRun(100, func() {
		for i := 0; i < 1000; i++ {
			result = ProcessStringPoolValue(s)
		}
	})
	t.Logf("Pool: %.0f allocations", allocsPool)
	t.Logf("No Pool: %.0f allocations", allocsNoPool)
	t.Logf("Not Pointer: %.0f allocations", allocsNotPointer)
}
