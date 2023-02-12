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
	maxContainerSizeBeforeBurst = 1 << 1
	initialContainerSize        = 1 << 12
)

type Trie struct {
	*trieNode
}

func New() *Trie {
	c := newTrieContainer(initialContainerSize)
	c.hybrid = true
	return &Trie{
		trieNode: newTrieNode(c),
	}
}

func (t *Trie) Put(key string, value Value) {
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

type stackItem struct {
	prefix  *byte
	item    node
	visited bool
}

func (t *Trie) ForEach(fn func(key string, value Value)) {
	stack := make([]*stackItem, 0, numberOfByteValues+1)
	stack = append(stack, &stackItem{
		prefix:  nil,
		visited: false,
		item:    t.trieNode,
	})

	prefix := make([]byte, 0)

	for len(stack) > 0 {
		n := stack[len(stack)-1]

		switch t := n.item.(type) {
		case *trieNode:
			if n.visited {
				stack = stack[:len(stack)-1]
				if len(prefix) > 0 {
					prefix = prefix[:len(prefix)-1]
				}
				continue
			}

			if n.prefix != nil {
				prefix = append(prefix, *n.prefix)
			}

			if t.validValue {
				fn(string(prefix), t.value)
			}

			var previousChild node
			for i := byteMaxValue; i >= 0; i-- {
				if t.children[i] != nil && t.children[i] != previousChild {
					p := byte(i)
					stack = append(stack, &stackItem{
						prefix:  &p,
						visited: false,
						item:    t.children[i],
					})
				}
				previousChild = t.children[i]
			}
			n.visited = true
		case *trieContainer:
			for _, key := range t.SortedKeys() {
				if t.hybrid {
					hybridPrefix := prefix
					if len(prefix) > 0 {
						hybridPrefix = prefix[:len(prefix)-1]
					}

					fn(string(append(hybridPrefix, key...)), t.pairs[key])
				} else {
					fn(string(append(prefix, key...)), t.pairs[key])
				}
			}
			stack = stack[:len(stack)-1]
		}
	}
}
