package hattrie

type trieNode struct {
	// TODO: Should we prealloc the map size for all possible byte values?
	children map[byte]*node
}

type node struct {
	*arrayHash
	*trieNode
}

func (n *node) findValue(key string) (ValueType, bool) {
	nearestNode, idx := n.findNearestNode(key)

	if nearestNode.arrayHash == nil {
		return 0, false
	}

	v, ok := (*n.arrayHash)[key[idx:]]
	return v, ok
}

func (n *node) findNearestNode(key string) (nearest *node, idx int) {
	nearest = n
	for i, head := range []byte(key) {
		newNearest := nearest.trieNode.children[head]
		idx = i

		if newNearest == nil {
			break
		}
		nearest = newNearest
	}
	return nearest, idx
}

func (n *node) burst() {
	if n.arrayHash != nil {
		n.trieNode = &trieNode{}

		for k, v := range *n.arrayHash {
			head := k[0]

			container, ok := n.trieNode.children[head]
			if !ok {
				container = &node{
					arrayHash: &arrayHash{},
				}
				n.trieNode.children[head] = container
			}

			(*container.arrayHash)[k[1:]] = v
		}
		n.arrayHash = nil
	}
}

func (n *node) insert(key string, value ValueType) bool {
	nearest, idx := n.findNearestNode(key)

	if len(*nearest.arrayHash) > maxHashSizeBeforeBurst {
		n.burst()
		return n.insert(key[idx:], value)
	}

	(*n.arrayHash)[key[idx:]] = value
	return true
}
