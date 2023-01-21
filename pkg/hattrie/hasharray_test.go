package hattrie

import (
	"fmt"
	"sort"
	"testing"
)

const arrayHashSize = maxHashSizeBeforeBurst

// TODO: Test several lengths of keys or just use a proper dictionary.
func getKey(i ValueType) string {
	return fmt.Sprintf("%d-%d", i, i)
	// return fmt.Sprintf("%d%d%d%d%d%d", i, i, i, i, i, i)
}

func TestArrayHashSortedKeys(t *testing.T) {
	h := newArrayHash()
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
		h[pair.key] = value
	}

	sortedKeys := h.SortedKeys()
	sorted := make([]struct {
		key   string
		value ValueType
	}, 0, len(sortedKeys))

	for _, key := range sortedKeys {
		value := h[key]

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

func BenchmarkArrayHashInsert(b *testing.B) {
	h := newArrayHash()
	// b.StopTimer()

	// // Start with the half of the maximum key-value pairs
	// for i := 0; i < arrayHashSize/2; i++ {
	// 	value := ValueType(i)
	// 	h[getKey(value)] = value
	// }
	// b.StartTimer()

	for i := 0; i < b.N; i++ {
		value := ValueType(i)
		h[getKey(value)] = value
	}
	b.ReportAllocs()
}

func BenchmarkArrayHashGet(b *testing.B) {
	b.StopTimer()
	h := arrayHash{}

	for i := 0; i < arrayHashSize; i++ {
		value := ValueType(i)
		h[getKey(value)] = value
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = h[getKey(1)]
	}

	b.ReportAllocs()
}

func BenchmarkArrayHashSortedIterate(b *testing.B) {
	b.StopTimer()
	h := newArrayHash()

	for i := 0; i < arrayHashSize; i++ {
		value := ValueType(i)
		h[getKey(value)] = value
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, k := range h.SortedKeys() {
			_ = h[k]
		}
	}

	b.ReportAllocs()
}

func BenchmarkArrayHashIterate(b *testing.B) {
	b.StopTimer()
	h := newArrayHash()

	for i := 0; i < arrayHashSize; i++ {
		value := ValueType(i)
		h[getKey(value)] = value
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for _, v := range h {
			_ = v
		}
	}

	b.ReportAllocs()
}
