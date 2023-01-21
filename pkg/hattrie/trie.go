// Package hattrie implements HAT-trie which is an optimized cache-friendly data structure
// that allows fast access to values by their associated keys.
//
// https://en.wikipedia.org/wiki/HAT-trie
package hattrie

const maxHashSizeBeforeBurst = 1 << 14

type Trie struct {
	root node
}

func New() *Trie {
	return &Trie{
		root: node{
			arrayHash: &arrayHash{},
		},
	}
}

func (t *Trie) Get(key string) (ValueType, bool) {
	return t.root.findValue(key)
}

// TODO: Should we normalize the byte representation when inserting?
// https://go.dev/blog/normalization
func (t *Trie) Put(key string, value ValueType) bool {
	return t.root.insert(key, value)
}
