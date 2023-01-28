package hattrie

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Use assert

const containerSize = maxContainerSizeBeforeBurst

// TODO: Test several lengths of keys or just use a proper dictionary.
func getKey(i ValueType) string {
	return fmt.Sprintf("%d-%d", i, i)
	// return fmt.Sprintf("%d%d%d%d%d%d", i, i, i, i, i, i)
}

// func TestTrieContainerPut(t *testing.T) {
// 	c := newTrieContainer(initialContainerSize)

// 	c.Put(getKey(1), 0, 1)

// 	c.Put(getKey(2), 1, 2)

// 	if c.pairs[getKey(1)[0+1:]] != 1 {
// 		t.Errorf("Wrong value retrieved for prefix 0: key %v", getKey(1)[0:])
// 	}

// 	if c.pairs[getKey(2)[1+1:]] != 2 {
// 		t.Errorf("Wrong value retrieved for prefix 1: key %v", getKey(2)[1:])
// 	}
// }

func TestTrieContainerPutHybrid(t *testing.T) {
	c := newTrieContainer(initialContainerSize)

	c.hybrid = true

	c.Put(getKey(1), 0, 1)

	c.Put(getKey(2), 1, 2)

	if c.pairs[getKey(1)[0:]] != 1 {
		t.Errorf("Wrong value retrieved for prefix 0: key %v", getKey(1)[0:])
	}

	if c.pairs[getKey(2)[1:]] != 2 {
		t.Errorf("Wrong value retrieved for prefix 1: key %v", getKey(2)[1:])
	}
}

func TestTrieContainerSortedKeys(t *testing.T) {
	c := newTrieContainer(initialContainerSize)
	pairs := []struct {
		key   string
		value ValueType
	}{
		{"c", 3},
		{"a", 1},
		{"b", 2},
	}

	for _, pair := range pairs {
		value := pair.value
		c.pairs[pair.key] = value
	}

	sortedKeys := c.SortedKeys()
	sorted := make([]struct {
		key   string
		value ValueType
	}, 0, len(sortedKeys))

	for _, key := range sortedKeys {
		value := c.pairs[key]

		sorted = append(sorted, struct {
			key   string
			value ValueType
		}{
			key:   key,
			value: value,
		})
	}

	if !sort.SliceIsSorted(sorted, func(i, j int) bool {
		return sorted[i].value < sorted[j].value
	}) {
		t.Errorf("SortedKeys didn't sort the keys correctly: %+v", sorted)
	}
}

func BenchmarkTrieContainerPut(b *testing.B) {
	b.StopTimer()
	c := newTrieContainer(initialContainerSize)
	// Start with the half of the maximum key-value pairs
	for i := 0; i < containerSize/2; i++ {
		value := ValueType(i)
		// c.pairs[getKey(value)] = value
		c.Put(getKey(value), 0, value)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		value := ValueType(i)
		// c.pairs[getKey(value)] = value
		c.Put(getKey(value), 0, value)
	}

	b.ReportAllocs()
}

func BenchmarkTrieContainerSortedIterate(b *testing.B) {
	b.StopTimer()
	c := newTrieContainer(initialContainerSize)

	for i := 0; i < containerSize; i++ {
		value := ValueType(i)
		c.pairs[getKey(value)] = value
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, k := range c.SortedKeys() {
			_ = c.pairs[k]
		}
	}

	b.ReportAllocs()
}

func BenchmarkTrieContainerIterate(b *testing.B) {
	b.StopTimer()
	c := newTrieContainer(initialContainerSize)

	for i := 0; i < containerSize; i++ {
		value := ValueType(i)
		c.pairs[getKey(value)] = value
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range c.pairs {
			_ = v
		}
	}

	b.ReportAllocs()
}

func BenchmarkClosure(b *testing.B) {
	ls := make([]int, 300)

	for i := 0; i < b.N; i++ {
		var val int

		fn := func(v int) {
			val = v
		}

		for i := range ls {
			fn(i)
		}

		assert.GreaterOrEqual(b, val, 0)
	}

	b.ReportAllocs()
}

func BenchmarkLoop(b *testing.B) {
	ls := make([]int, 300)

	for i := 0; i < b.N; i++ {
		var val int

		for i := range ls {
			val = i
		}

		assert.GreaterOrEqual(b, val, 0)
	}

	b.ReportAllocs()
}
