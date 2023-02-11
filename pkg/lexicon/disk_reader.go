package lexicon

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sort"
)

const (
	readerFilePermissions = 0o600
	readerFileFlags       = os.O_RDONLY
)

type DiskLexiconReader struct {
	*diskLexicon
}

func NewDiskLexiconReader(dirpath string) (*DiskLexiconReader, error) {
	d, err := newDiskLexicon(dirpath, readerFileFlags, readerFilePermissions)
	if err != nil {
		return nil, fmt.Errorf("error creating disk lexicon reader from '%s': %+v", dirpath, err)
	}

	return &DiskLexiconReader{
		diskLexicon: d,
	}, nil
}

func (r *DiskLexiconReader) sortedID(position int) (id ID, err error) {
	_, err = r.indices.ReadAt(r.positionBuffer[:], int64(position*positionByteSize))
	if err != nil {
		return id, fmt.Errorf("error reading position '%d' in sorted lexicon file: %+v", position, err)
	}
	return binary.LittleEndian.Uint32(r.positionBuffer[:]), nil
}

func (r *DiskLexiconReader) Word(id ID) (word string, err error) {
	_, err = r.indices.ReadAt(r.positionBuffer[:], int64(id*positionByteSize))
	if err != nil {
		return word, fmt.Errorf("error reading position '%d' in lexicon indices file: %+v", id, err)
	}

	offset := binary.LittleEndian.Uint32(r.positionBuffer[:])
	_, err = r.lexicon.Seek(int64(offset), 0)
	if err != nil {
		return word, fmt.Errorf("error seeking offset '%d' in lexicon file: %+v", offset, err)
	}

	reader := bufio.NewReader(r.lexicon)
	word, err = reader.ReadString(nullByte)
	if err != nil {
		return word, fmt.Errorf("error reading after offset '%d' in lexicon file: %+v", offset, err)
	}
	return word, nil
}

func (l *DiskLexiconReader) ID(word string) (id ID, err error) {
	fi, err := l.sorted.Stat()
	if err != nil {
		return id, fmt.Errorf("error obtaining stats for sorted lexicon file: %+v", err)
	}

	numberOfIndices := int(fi.Size() / positionByteSize)

	positionInIndices := sort.Search(numberOfIndices, func(currentPosition int) bool {
		currentID, idErr := l.sortedID(currentPosition)
		if idErr != nil {
			return false
		}
		currentWord, wordErr := l.Word(currentID)
		if wordErr != nil {
			return false
		}
		return currentWord >= word
	})

	if positionInIndices == numberOfIndices {
		return id, fmt.Errorf("cannot find word '%s' in the sorted lexicon file", word)
	}

	id, err = l.sortedID(positionInIndices)
	if err != nil {
		return id, fmt.Errorf("error retrieving id from position '%d' in lexicon indices file: %+v", positionInIndices, err)
	}
	return id, nil
}

func (r *DiskLexiconReader) Close() error {
	return r.diskLexicon.Close()
}
