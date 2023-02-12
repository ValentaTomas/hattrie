package hattrie

import (
	"github.com/yourbasic/radix"
)

// Value is an alias for values in the array hash.
type Value = uint32

// pairs is a type definition for string to uint64 map.
//
// The map in Go is implemented with buckets as arrays with a fixed length 8 and a factor of 6.5.
// According to the implemetation https://github.com/golang/go/blob/master/src/runtime/map.go,
// the key and values are saved in separate arrays, instead of alternating length-key-value in one array, as in the C implementation mentioned further down.
//
// I haven't tried using the exact implementation from the HAT-trie C lib (https://github.com/dcjones/hat-trie/blob/master/src/ahtable.h) yet.
// That would require using unsafe pointers directly and I want to avoid that until I know that this is a performance bottleneck.
//
// All other custom array hash implementations that I tried were slower that the default Go map for the 1 to 16384 entries that the array hash should handle before the trie node burst.
type pairs map[string]Value

type trieContainer struct {
	pairs
	hybrid     bool
	splitStart byte
	splitEnd   byte
}

func newTrieContainer(size int) *trieContainer {
	return &trieContainer{
		pairs:    make(pairs, size),
		splitEnd: byteMaxValue,
	}
}

// SortedKeys returns a slice of hash's keys, sorted lexicographically with the radix sort.
// This slice could be then used for iterating the hash in order.
func (c *trieContainer) SortedKeys() []string {
	keys := make([]string, 0, len(c.pairs))

	for k := range c.pairs {
		keys = append(keys, k)
	}

	// Radix sort should work well with the HAT-trie, because even though the sorting needs additional space,
	// we are only sorting one array hash at the time, so the space requirement increase should be constant (at worst for 16384 - max hash table size before burst).
	// Iterating over 2^14 elements in the map in order takes 2992016 ns for the radix sort vs 4705317 ns for the default sort.String.
	radix.Sort(keys)

	return keys
}

func (c *trieContainer) Insert(key string, value Value) {
	if _, ok := c.pairs[key]; !ok {
		c.pairs[key] = value
	}
}
