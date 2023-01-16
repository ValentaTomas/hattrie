package hasharray

const (
	initialTableSize   = 100000.0
	maxTableLoadFactor = 4096
)

type HashTable struct {
	slots uint64
	pairs uint64
}
