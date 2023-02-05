package lexicon

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/valentatomas/fsa3/pkg/hattrie"
)

type ID = hattrie.Value

const (
	lexiconFileName        = "hat.lex"
	lexiconIndicesFileName = "hat.lex.idx"
	sortedFileName         = "hat.lex.srt"

	nullByte = byte(0)

	positionByteSize = 32 / 8
)

type diskLexicon struct {
	lexicon *os.File
	indices *os.File
	sorted  *os.File

	positionBuffer [positionByteSize]byte
}

func newDiskLexicon(dirpath string, flags int, permissions fs.FileMode) (*diskLexicon, error) {
	cleanPath := filepath.Clean(dirpath)

	lexiconPath := filepath.Join(cleanPath, lexiconFileName)
	lexicon, err := os.OpenFile(lexiconPath, flags, permissions)
	if err != nil {
		return nil, fmt.Errorf("error opening lexicon file '%s': %+v", lexiconPath, err)
	}

	indicesPath := filepath.Join(cleanPath, lexiconIndicesFileName)
	indices, err := os.OpenFile(indicesPath, flags, permissions)
	if err != nil {
		return nil, fmt.Errorf("error opening lexicon indices file '%s': %+v", indicesPath, err)
	}

	sortedPath := filepath.Join(cleanPath, sortedFileName)
	sorted, err := os.OpenFile(sortedPath, flags, permissions)
	if err != nil {
		return nil, fmt.Errorf("error opening sorted lexicon file '%s': %+v", sortedPath, err)
	}

	return &diskLexicon{
		lexicon: lexicon,
		indices: indices,
		sorted:  sorted,
	}, nil
}
