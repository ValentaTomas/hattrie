package hattrie

// Usign interface and struct assertion between trieNode and trieContainer.
// This way we don't have to embed trieNode and trieContainer and differentiate between them based on a nil pointer or
// have overhead from using interface methods.
type node interface{}

type trieNode struct {
	// We can directly use byte as an index into the array - the array then behaves like a map with a fixed size and no additional overhead.
	// TODO: Make the trie more compact by using arrays with size 128 - splitting the byte into two nodes each handling 4 bits.
	// The second node would be present only if necessary and would be a child of the first node.
	children   [byteMaxValue]node
	value      ValueType
	validValue bool
}

func newTrieNode(child node) *trieNode {
	t := &trieNode{}

	// https://groups.google.com/g/golang-dev/c/35W8LvT51vg
	// Prevent range over array copy with &.
	for i := range &t.children {
		t.children[i] = child
	}
	return t
}

func (n *trieNode) setValue(value ValueType) {
	n.value = value
	n.validValue = true
}

func (n *trieNode) findNearest(key string) (nearest node, parent *trieNode, prefixIdx int) {
	nearest = n
	parent = n
	for i, head := range []byte(key) {
		switch t := nearest.(type) {
		case *trieNode:
			parent = t
			maybeNearest := t.children[head]
			if maybeNearest != nil {
				nearest = maybeNearest
			} else {
				break
			}
		case *trieContainer:
			break
		}
		prefixIdx = i
	}
	return nearest, parent, prefixIdx
}

func (n *trieNode) splitContainer(child *trieContainer) (*trieNode, *trieContainer) {
	// Turn pure into hybrid
	if !child.hybrid {
		newParent := newTrieNode(child)
		n.children[child.splitStart] = newParent

		// TODO: Bucket empty key? -> move to trie node

		child.hybrid = true
		return newParent, child
	}

	var occurrences [byteMaxValue]int

	for k := range child.pairs {
		occurrences[k[0]]++
	}

	split := int(child.splitStart)
	totalSize := len(child.pairs)
	leftSize := occurrences[split]
	rightSize := totalSize - leftSize

	for i, o := range occurrences[1:] {
		delta := abs((leftSize + o) - (rightSize - o))
		if delta <= leftSize-rightSize && leftSize+o < totalSize {
			split = i
			leftSize += o
			rightSize += o
		} else {
			break
		}
	}

	// TODO: Handle the preallocation and special cases better
	left := newTrieContainer(leftSize)

	left.splitStart = child.splitStart
	left.splitEnd = byte(split)

	if left.splitStart != left.splitEnd {
		left.hybrid = true
	}

	right := newTrieContainer(rightSize)

	right.splitStart = byte(split + 1)
	right.splitEnd = child.splitEnd

	if right.splitStart != right.splitEnd {
		right.hybrid = true
	}

	for i := left.splitStart; i <= left.splitEnd; i++ {
		n.children[i] = left
	}

	for i := right.splitStart; i <= right.splitEnd; i++ {
		n.children[i] = right
	}

	for k, v := range child.pairs {
		if k[0] <= left.splitEnd {
			left.Put(k, 0, v)
		} else {
			right.Put(k, 0, v)
		}
	}

	if leftSize >= maxContainerSizeBeforeBurst {
		return n, left
	} else if rightSize >= maxContainerSizeBeforeBurst {
		return n, right
	}

	return n, left
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

type trieIteratorStack struct {
	c     byte
	level uint
	n     node
	next  *trieIteratorStack
}

type TrieIterator struct {
	key    string
	prefix string

	// TODO: Do we need empty key values at all?
	// emptyKey   bool
	// emptyValue ValueType

	stack *trieIteratorStack
}

// FSA needs to process all pairs at once, so we don't have to implement an iterator.
// TODO: Check if using the function (that can use closure) is affecting performance too much.
func (t *Trie) iterate(sorted bool, fn func(key string, value ValueType)) {
	i := &TrieIterator{}

	i.
}
