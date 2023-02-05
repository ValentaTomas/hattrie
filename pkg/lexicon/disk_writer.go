package lexicon

import (
	"encoding/binary"
	"fmt"
	"os"
)

const (
	writerFilePermissions = 0o600
	writerFileFlags       = os.O_WRONLY | os.O_CREATE
)

type DiskLexiconWriter struct {
	*diskLexicon
}

func NewDiskLexiconWriter(dirpath string) (*DiskLexiconWriter, error) {
	d, err := newDiskLexicon(dirpath, writerFileFlags, writerFilePermissions)
	if err != nil {
		return nil, fmt.Errorf("error creating disk lexicon writer from '%s': %+v", dirpath, err)
	}

	return &DiskLexiconWriter{
		diskLexicon: d,
	}, nil
}

func (w *DiskLexiconWriter) Write(word string) error {
	_, err := w.indices.WriteString(word + string(nullByte))
	if err != nil {
		return fmt.Errorf("error writing word '%s' to lexicon file: %+v", word, err)
	}

	binary.LittleEndian.PutUint32(w.positionBuffer[:], uint32(len(word)))
	_, err = w.indices.Write(w.positionBuffer[:])
	if err != nil {
		return fmt.Errorf("error writing word '%s' to lexicon indices file: %+v", word, err)
	}

	return nil
}

func (w *DiskLexiconWriter) WriteSorted(id ID) error {
	binary.LittleEndian.PutUint32(w.positionBuffer[:], id)
	_, err := w.indices.Write(w.positionBuffer[:])
	if err != nil {
		return fmt.Errorf("error writing id '%d' to lexicon sorted file: %+v", id, err)
	}

	return nil
}
