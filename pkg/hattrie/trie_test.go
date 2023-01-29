package hattrie

import (
	"log"
	"testing"
)

func TestTriePut(t *testing.T) {
	tr := New()

	tr.Put("a", 1)

	tr.ForEach(func(key string, value Value) {
		log.Printf("foreach - %s, %d", key, value)
	})
}

// func TestTrieMultiplePut(t *testing.T) {
// 	tr := New()

// 	tr.Put("a", 1)
// 	tr.Put()
// 	tr.Put()
// 	tr.Put()
// }

// func TestTrieForEach(t *testing.T) {
// }
