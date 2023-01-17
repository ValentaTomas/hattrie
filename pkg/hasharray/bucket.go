package hasharray

import (
	"bytes"

	"golang.org/x/exp/slices"
)

type bucket struct {
	sizes  []int
	values []byte
}

type bucketIndices struct {
	sizes      int
	valueStart int
	valueEnd   int
}

func (b *bucket) insert(key []byte) {
	b.sizes = append(b.sizes, len(key))
	b.values = append(b.values, key...)
}

func (b *bucket) delete(key []byte) {
	indices, ok := b.get(key)
	if !ok {
		return
	}

	slices.Delete(b.values, indices.valueStart, indices.valueEnd)
	slices.Delete(b.sizes, indices.sizes, indices.sizes+1)
}

func (b *bucket) get(key []byte) (indices *bucketIndices, ok bool) {
	keySize := len(key)

	for sizesIdx, size := range b.sizes {
		if size != keySize {
			indices.valueStart += size
			continue
		}

		indices.valueEnd = indices.valueStart + size

		if bytes.Compare(b.values[indices.valueStart:indices.valueEnd], key) == 0 {
			indices.sizes = sizesIdx
			return indices, true
		}
	}
	return nil, false
}

func (b *bucket) Size() int {
	return len(b.sizes)
}

func (b *bucket) getValuePosition(idx int) (offset, size int) {
	for sizesIdx, currentSize := range b.sizes {
		size = currentSize
		if sizesIdx < idx {
			offset += size
		}
	}
	return offset, size
}

func (b *bucket) Value(idx int) []byte {
	offset, size := b.getValuePosition(idx)
	return b.values[offset : offset+size]
}
