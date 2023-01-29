package hattrie

// Usign interface and struct assertion between trieNode and trieContainer.
// This way we don't have to embed trieNode and trieContainer and differentiate between them based on a nil pointer or
// have overhead from using interface methods.
type node interface{}

type trieNode struct {
	// We can directly use byte as an index into the array - the array then behaves like a map with a fixed size and no additional overhead.
	// TODO: Make the trie more compact by using arrays with size 128 - splitting the byte into two nodes each handling 4 bits.
	// The second node would be present only if necessary and would be a child of the first node.
	children   [numberOfByteValues]node
	value      Value
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

func (n *trieNode) setValue(value Value) {
	if !n.validValue {
		n.value = value
		n.validValue = true
	}
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
		prefixIdx = i + 1
	}
	return nearest, parent, prefixIdx
}

func (n *trieNode) splitContainer(child *trieContainer) (*trieNode, *trieContainer) {
	// Turn pure into hybrid
	if !child.hybrid {
		newParent := newTrieNode(child)
		n.children[child.splitStart] = newParent

		child.hybrid = true
		return newParent, child
	}

	var occurrences [numberOfByteValues]int

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
			left.Insert(k, 0, v)
		} else {
			right.Insert(k, 0, v)
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

type stackItem struct {
	prefix  *byte
	item    node
	visited bool
}

// FSA needs to process all pairs, so we don't have to implement an iterator.
// TODO: Check if using the function (that can use closure) is affecting performance too much.
func (t *Trie) ForEach(fn func(key string, value Value)) {
	// TODO: What is the ideal size for the preallocated stack?
	stack := make([]*stackItem, 0, numberOfByteValues+1)
	stack = append(stack, &stackItem{
		prefix:  nil,
		visited: false,
		item:    t.trieNode,
	})

	prefix := make([]byte, 0, t.longestKeySize)

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
			for _, k := range t.SortedKeys() {
				subkey := t.getKey(k, 0)
				if t.hybrid {
					// ERR: [1:] slice out of bounds
					fn(string(append(prefix, subkey[1:]...)), t.pairs[subkey])
				} else {
					fn(string(append(prefix, subkey...)), t.pairs[subkey])
				}
			}
			stack = stack[:len(stack)-1]
		}
	}
}
