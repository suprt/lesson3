package main

import (
	"encoding/json"
	"testing"
)

func TestAllocations(t *testing.T) {
	const runs = 1000

	input := []byte(`{
		"name": "Alice",
		"items": {
			"foo": "bar",
			"hello": "world"
		}
	}`)

	poolAllocs := testing.AllocsPerRun(100, func() {
		for i := 0; i < runs; i++ {
			data := dataPool.Get().(*RequestData)

			_ = json.Unmarshal(input, data)

			data.Reset()
			dataPool.Put(data)
		}
	})

	noPoolAllocs := testing.AllocsPerRun(100, func() {
		for i := 0; i < runs; i++ {
			var data RequestData

			_ = json.Unmarshal(input, &data)
		}
	})

	t.Logf("Pool: %.0f allocations", poolAllocs)
	t.Logf("No Pool: %.0f allocations", noPoolAllocs)
}

func TestRequestDataReset(t *testing.T) {
	data := dataPool.Get().(*RequestData)
	defer dataPool.Put(data)

	data.Name = "Alice"
	data.Items["secret"] = "123"

	data.Reset()

	if data.Name != "" {
		t.Fatalf("Name was not reset: %q", data.Name)
	}

	if len(data.Items) != 0 {
		t.Fatalf("Items was not cleared: %v", data.Items)
	}
}
