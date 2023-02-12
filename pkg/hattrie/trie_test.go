package hattrie

import (
	"log"
	"testing"
)

func TestTriePut(t *testing.T) {
	tr := New()

	tr.Put("a", 2)
	tr.Put("aqwd", 2)
	tr.Put("ppol", 2)
	tr.Put("cqq", 2)
	tr.Put("ca", 4)
	tr.Put("c", 5)
	tr.Put("asd", 5)

	tr.ForEach(func(key string, value Value) {
		log.Printf("foreach - %s, %d", key, value)
	})
}

func TestTrieForEach(t *testing.T) {
	tr := New()
	for i := 0; i < 200; i++ {
		b := []byte{byte(i)}

		tr.Put(string(b), Value(i))
		if i == 127 {
			log.Default()
		}
		// log i values

		log.Printf("put - %x, %d", byte(i), Value(i))
	}

	i := 0
	tr.ForEach(func(key string, value Value) {
		log.Printf("foreach - %s, %d", key, value)
		i++
	})
	if i != 200 {
		t.Errorf("expected 200 iterations, got %d", i)
	}
}

func BenchmarkTriePut(b *testing.B) {
	tr := New()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tr.Put(string(rune(i)), Value(i))
	}
}

func BenchmarkTrieForEach(b *testing.B) {
	tr := New()
	for i := 0; i < 200; i++ {
		tr.Put(string(rune(i)), Value(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.ForEach(func(key string, value Value) {})
	}
}
