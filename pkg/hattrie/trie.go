package hattrie

const (
	maxBucketSize = 16384
)

type node interface{}

type trieNode struct{}

type Trie struct {
	root       node
	KeysStored int
}

func New(size int) *Trie {
	return &Trie{}
}
