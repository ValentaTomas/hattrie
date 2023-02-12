package hattrie

type node interface{}

type trieNode struct {
	children   [numberOfByteValues]node
	value      Value
	validValue bool
}

func newTrieNode(child node) *trieNode {
	t := &trieNode{}

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

NearLoop:
	for _, head := range []byte(key) {
		switch t := nearest.(type) {
		case *trieNode:
			parent = t
			maybeNearest := t.children[head]
			if maybeNearest != nil {
				prefixIdx++
				nearest = maybeNearest
			} else {
				break NearLoop
			}
		case *trieContainer:
			break NearLoop
		}
	}
	return nearest, parent, prefixIdx
}

func (n *trieNode) splitContainer(child *trieContainer) (*trieNode, *trieContainer) {
	if !child.hybrid {
		newParent := newTrieNode(child)
		n.children[child.splitStart] = newParent

		if value, ok := child.pairs[""]; ok {
			newParent.setValue(value)
			delete(child.pairs, "")
		}

		child.hybrid = true
		child.splitStart = 0
		child.splitEnd = byteMaxValue
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

	// for split+1 < int(child.splitEnd) {
	// o := occurrences[split+1]
	for _, o := range occurrences[split+1 : child.splitEnd] {
		delta := abs((leftSize + o) - (rightSize - o))
		if delta <= leftSize-rightSize && leftSize+o < totalSize {
			split++
			leftSize += o
			rightSize -= o
		} else {
			break
		}
	}

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

	for i := left.splitStart; ; i++ {
		n.children[i] = left

		if i == left.splitEnd {
			break
		}
	}

	for i := right.splitStart; ; i++ {
		n.children[i] = right

		if i == right.splitEnd {
			break
		}
	}

	for k, v := range child.pairs {
		if k[0] <= left.splitEnd {
			if left.hybrid {
				left.Insert(k, v)
			} else {
				left.Insert(k[1:], v)
			}
		} else {
			if right.hybrid {
				right.Insert(k, v)
			} else {
				right.Insert(k[1:], v)
			}
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
