package hattrie

import (
	"log"
	"testing"
)

func TestTriePut(t *testing.T) {
	tr := New()

	tr.Put("cqq", 2)
	tr.Put("ca", 4)
	tr.Put("c", 5)

	tr.ForEach(func(key string, value Value) {
		log.Printf("foreach - %s, %d", key, value)
	})
}
