package hattrie

const byteMaxValue = ^byte(0)

// TODO: Implement the hybrid trie.
type trieNode struct {
	prefix byte
	// We can directly use byte as an index into the array - the array then behaves like a map with a fixed size and no additional overhead.
	// TODO: Make the trie more compact by using arrays with size 128 - splitting the byte into two nodes each handling 4 bits.
	// The second node would be present only if necessary and would be a child of the first node.
	children [byteMaxValue]node
}

// Usign interface and struct assertion between trieNode and arrayHash.
// This way we don't have to embed trieNode and arrayHash and differentiate between them based on a nil pointer or
// have overhead from using interface methods.
type node interface{}

type Apply func(key string, value ValueType)

func findValue(start node, key string) (ValueType, bool) {
	nearest, idx := findNearestTrieNode(start, key)

	if t, ok := nearest.(*arrayHash); ok {
		// TODO: We may need to check for nil values in type assert switches
		v, success := (*t)[key[idx:]]
		return v, success
	}

	return 0, false
}

// We are using iteration instead of tail recursion for performance.
func findNearestTrieNode(start node, key string) (nearest node, idx int) {
	nearest = start
	// []byte conversion (optimized by the compiler) is necessary so we don't iterate over runes.
	for i, head := range []byte(key) {
		var newNearest node

		if t, ok := nearest.(*trieNode); ok {
			newNearest = t.children[head]
		}

		if _, ok := newNearest.(*trieNode); ok {
			nearest = newNearest
			idx = i
		} else {
			break
		}
	}
	return nearest, idx
}

func burst(n node, head byte) {
	if t, ok := n.(*trieNode); ok {
		child := t.children[head]

		if oldHash, ok := child.(*arrayHash); ok {
			container := &trieNode{}
			t.children[head] = container

			for k, v := range *oldHash {
				newChild := container.children[k[0]]

				if subHash, ok := newChild.(*arrayHash); ok {
					if subHash == nil {
						subHash = &arrayHash{}
						container.children[k[0]] = subHash
					}
					(*subHash)[k[1:]] = v
				}
			}
		}
	}
}

// TODO: Use iteration instead of a recursion.
func put(start node, key string, value ValueType) bool {
	nearest, idx := findNearestTrieNode(start, key)

	if t, ok := nearest.(*arrayHash); ok {
		if len(*t) < maxHashSizeBeforeBurst {
			(*t)[key[idx:]] = value
			return true
		}
	}

	burst(nearest, key[0])
	return put(start, key[idx:], value)
}

// TODO: Try using iterator struct instead of passing a function to apply.
func forEach(start node, apply Apply, sorted bool) {
	stack := []node{start}
	// var prefix []byte

	for len(stack) > 0 {
		visiting := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		// TODO: Check if we are iterating nodes in the lexicographic order

		switch t := visiting.(type) {
		case *trieNode:
			stack = append(stack, t.children[:]...)

		case *arrayHash:
			if sorted {
				for _, key := range t.SortedKeys() {
					value := (*t)[key]
					// TODO: Assemble the whole key
					apply(key, value)
				}
			} else {
				for key, value := range *t {
					// TODO: Assemble the whole key
					apply(key, value)
				}
			}
		}
	}
}
