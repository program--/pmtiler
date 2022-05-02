package tiles

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"os"
	"sort"
)

var ErrDirSizeTooLarge error = errors.New("leaf directories not supported")

type HashOffsetMap map[uint64]uint64

type Writer struct {
	file         *os.File
	offset       uint64
	tiles        EntryAscending
	hashToOffset HashOffsetMap
	header       *Header
}

func (w *Writer) WriteTile(coords Coordinates, data []byte) {
	hash := fnv.New64a()
	hash.Write(data)
	tileHash := hash.Sum64()

	existingOffset, ok := w.hashToOffset[tileHash]

	if ok {
		w.tiles = append(w.tiles, Entry{
			Coords: coords,
			Rng: Range{
				Offset: existingOffset,
				Length: uint32(len(data)),
			},
		})
	} else {
		w.file.Write(data)
		w.tiles = append(w.tiles, Entry{
			Coords: coords,
			Rng: Range{
				Offset: w.offset,
				Length: uint32(len(data)),
			},
		})
		w.hashToOffset[tileHash] = w.offset
		w.offset += uint64(len(data))
	}
}

func (w *Writer) Finalize(metadata *Metadata) error {
	defer w.file.Close()

	if len(w.tiles) > MAX_DIR_SIZE {
		return ErrDirSizeTooLarge
	}

	if _, err := w.file.Seek(0, 0); err != nil {
		return err
	}

	meta_bytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	w.header.MetadataLength = uint32(len(meta_bytes))
	w.header.HeaderMetadata = metadata
	w.header.RootDirLength = uint16(w.tiles.Len())
	if err := w.header.Write(w.file); err != nil {
		return err
	}

	sort.Sort(EntryAscending(w.tiles))

	if err := w.tiles.Write(w.file); err != nil {
		return err
	}

	return nil
}

func NewWriter(path string) Writer {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	empty := make([]byte, HEADER_SIZE)
	_, err = f.Write(empty)
	if err != nil {
		panic(err)
	}

	header := new(Header)
	header.Version = uint16(2)

	return Writer{
		file:         f,
		offset:       HEADER_SIZE,
		tiles:        nil,
		hashToOffset: make(HashOffsetMap),
		header:       header,
	}
}
