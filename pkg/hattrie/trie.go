// Package hattrie implements HAT-trie which is an optimized cache-friendly data structure
// that allows fast access to values by their associated keys.
//
// https://en.wikipedia.org/wiki/HAT-trie
package hattrie

const maxHashSizeBeforeBurst = 1 << 14

// We don't need the delete and get methods for the FSA use case.
type Trie struct {
	root node
}

func New() *Trie {
	return &Trie{
		root: &arrayHash{},
	}
}

// TODO: Should we normalize the UNICODE representation when inserting?
// https://go.dev/blog/normalization
func (t *Trie) Put(key string, value ValueType) bool {
	return put(t.root, key, value)
}

// We only need sorted iteration, unsorted used only internally.
// TODO: Do we need sorted for FSA?
func (t *Trie) ForEach(apply Apply) {
	forEach(t.root, apply, true)
}
