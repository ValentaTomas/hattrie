package hattrie

import (
	"fmt"
	"sort"
	"testing"
)

const containerSize = maxContainerSizeBeforeBurst

func getKey(i Value) string {
	return fmt.Sprintf("%d-%d", i, i)
}

func TestTrieContainerSortedKeys(t *testing.T) {
	c := newTrieContainer(initialContainerSize)
	pairs := []struct {
		key   string
		value Value
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
		value Value
	}, 0, len(sortedKeys))

	for _, key := range sortedKeys {
		value := c.pairs[key]

		sorted = append(sorted, struct {
			key   string
			value Value
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
	c := newTrieContainer(initialContainerSize)
	// Start with the half of the maximum key-value pairs
	for i := 0; i < containerSize/2; i++ {
		value := Value(i)
		// c.pairs[getKey(value)] = value
		c.Insert(getKey(value), value)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		value := Value(i)
		// c.pairs[getKey(value)] = value
		c.Insert(getKey(value), value)
	}
	b.ReportAllocs()
}

func BenchmarkTrieContainerSortedIterate(b *testing.B) {
	c := newTrieContainer(initialContainerSize)

	for i := 0; i < containerSize; i++ {
		value := Value(i)
		c.pairs[getKey(value)] = value
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, k := range c.SortedKeys() {
			_ = c.pairs[k]
		}
	}
	b.ReportAllocs()
}

func BenchmarkTrieContainerIterate(b *testing.B) {
	c := newTrieContainer(initialContainerSize)

	for i := 0; i < containerSize; i++ {
		value := Value(i)
		c.pairs[getKey(value)] = value
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range c.pairs {
			_ = v
		}
	}
	b.ReportAllocs()
}
