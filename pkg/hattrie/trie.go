// Package hattrie implements HAT-trie which is an optimized cache-friendly data structure
// that allows fast access to values by their associated keys.
// The HAT-trie also allows us to iterate over the stored key-value pairs in a lexical order.
//
// The provided HAT-trie is customized for the FSA use case — for example we don't implement the get and remove methods.
//
// https://en.wikipedia.org/wiki/HAT-trie
package hattrie

import (
	"math"
)

const (
	byteMaxValue                = math.MaxUint8
	numberOfByteValues          = byteMaxValue + 1
	maxContainerSizeBeforeBurst = 1 << 14
	initialContainerSize        = 1 << 12
)

// We don't need the delete and get methods for the FSA use case.
// The exposed iterator can also always interate in a sorted order.
// TODO: Do we need to handle empty/nil key pairs?
type Trie struct {
	*trieNode
	// We can easily keep track of the longest key size because we don't need to delete from the HAT-trie.
	// We use this value to prevent reallocating slices when iterating through the HAT-trie.
	longestKeySize int
}

func New() *Trie {
	c := newTrieContainer(initialContainerSize)
	c.hybrid = true
	return &Trie{
		trieNode: newTrieNode(c),
	}
}

// We ignore the empty key "" to simplify the code. That should be ok in the context of FSA.
func (t *Trie) Put(key string, value Value) {
	// TODO: Check if this optimization is worth it.
	size := len(key)
	if t.longestKeySize < size {
		t.longestKeySize = size
	}

	nearest, parent, prefixIdx := t.trieNode.findNearest(key)

	remainingKey := key[prefixIdx:]

	switch n := nearest.(type) {
	case *trieNode:
		if len(remainingKey) == 0 {
			n.setValue(value)
			return
		}
	case *trieContainer:
		if parent != t.trieNode && len(remainingKey) == 0 && n.hybrid {
			parent.setValue(value)
			return
		}

		if n.hybrid {
			n.Insert(key[prefixIdx-1:], value)
		} else {
			n.Insert(remainingKey, value)
		}

		for len(n.pairs) >= maxContainerSizeBeforeBurst {
			parent, n = parent.splitContainer(n)
		}
	}
}
