# FSA3
Dynamically creating a FSA that can then be used to effectively map between words in a big corpus and a number index.
This can be also used for things like keyword extraction.

## Usage

### Test
```sh
make test-fast
```

```sh
make benchmark
```

## Articles
- https://nlp.fi.muni.cz/raslan/raslan16.pdf#page=151
- https://nlp.fi.muni.cz/raslan/raslan13.pdf#page=71
- http://www.jandaciuk.pl/fsa.html#FSApack
- https://en.wikipedia.org/wiki/HAT-trie
- https://tessil.github.io/2017/06/22/hat-trie.html
- https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
- https://www.practical-go-lessons.com/chap-7-hexadecimal-octal-ascii-utf8-unicode-runes
- https://www.practical-go-lessons.com/chap-34-benchmarks
- https://homepages.dcc.ufmg.br/~nivio/papers/cikm07.pdf
- https://arxiv.org/pdf/2006.09973.pdf
- https://go.dev/blog/strings
- https://www.joelonsoftware.com/2003/10/08/the-absolute-minimum-every-software-developer-absolutely-positively-must-know-about-unicode-and-character-sets-no-excuses/
- https://cs.wikipedia.org/wiki/Trie
- https://en.wikipedia.org/wiki/Finite-state_machine
- https://blog.gopheracademy.com/advent-2014/bloom-filters/
- https://www.academia.edu/45458959/HAT_Trie_A_Cache_Conscious_Trie_Based_Data_Structure_For_Strings
- https://abhinavg.net/2020/03/12/pointers-as-map-keys/
- https://go.dev/blog/slices
- https://go.dev/blog/normalization
- https://go101.org/article/memory-layout.html
- https://planetscale.com/blog/generics-can-make-your-go-code-slower
- https://go101.org/article/memory-block.html
- https://stackoverflow.com/questions/41030545/are-we-overusing-pass-by-pointer-in-go
- https://pthevenet.com/posts/programming/go/bytesliceindexedmaps/
- https://en.m.wikipedia.org/wiki/Open_addressing
- https://en.m.wikipedia.org/wiki/Linear_probing
- https://en.m.wikipedia.org/wiki/MurmurHash
- https://github.com/golang/go/wiki/CompilerOptimizations

## Code Reference
- Burst trie - https://github.com/nlfiedler/sortingo/blob/master/sort/burstsort.go
- Tests - https://github.com/raviqqe/hamt/blob/v2/entry.go
- Tests - from https://github.com/dghubble/trie/blob/main/bench_test.go
- Hamt trie - https://github.com/raviqqe/hamt
- HAT-trie used in FSA3 - https://github.com/dcjones/hat-trie
- HAT-trie - https://github.com/Tessil/hat-trie
- Array hash - https://github.com/Tessil/array-hash
- Bloom filter? - https://github.com/GuillaumeHolley/BloomFilterTrie
- Go map (hash array) - https://github.com/golang/go/blob/master/src/runtime/map.go