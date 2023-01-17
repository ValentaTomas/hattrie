package hasharray

import (
	"github.com/spaolacci/murmur3"
)

const (
	defaultBuckets     = 100000
	maxTableLoadFactor = 4096
)

type HashTable struct {
	flag uint
	c0   string
	c1   string

	Size           int
	initialBuckets int

	buckets []*bucket
}

func New(initialBuckets int) *HashTable {
	return &HashTable{
		initialBuckets: initialBuckets,
		flag:           0,
		c0:             "\\0",
		c1:             "\\0",
		buckets:        make([]*bucket, initialBuckets),
	}
}

func NewEmpty() *HashTable {
	return New(defaultBuckets)
}

func (t *HashTable) Clear() {
	t.Size = 0
	t.buckets = make([]*bucket, t.initialBuckets)
}

func (t *HashTable) getBucket(key []byte) *bucket {
	hash := int(murmur3.Sum32(key))
	bucketIdx := hash % len(t.buckets)
	return t.buckets[bucketIdx]
}

func (t *HashTable) Get(key string, insert bool) {
	byteKey := []byte(key)
	bucket := t.getBucket(byteKey)
	_, found := bucket.get(byteKey)

	if !found && insert {
		bucket.insert(byteKey)
	}
}

func (t *HashTable) Delete(key string) {
	byteKey := []byte(key)
	bucket := t.getBucket(byteKey)
	bucket.delete(byteKey)
}

type HashTableIterator struct {
	table     *HashTable
	bucketIdx int
	sizeIdx   int
}

func (t *HashTable) Iterate(sorted bool) Iterator[*string] {
	iterator := &HashTableIterator{
		table:     t,
		bucketIdx: 0,
		sizeIdx:   0,
	}
	return iterator
}

func (i *HashTableIterator) Next() (*string, bool) {
	if i.table.Size >= i.bucketIdx {
		return nil, false
	}

	bucket := i.table.buckets[i.bucketIdx]

	if bucket.Size() >= i.sizeIdx {
		i.bucketIdx++
		return i.Next()
	}

	value := string(bucket.Value(i.sizeIdx))

	i.sizeIdx++

	return &value, true
}
