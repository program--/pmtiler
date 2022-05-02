package tiles

import (
	"encoding/binary"
	"os"
)

const MAX_DIR_SIZE = 21845

type Coordinates struct {
	Z uint8
	X uint32
	Y uint32
}

type Range struct {
	Offset uint64
	Length uint32
}

type Entry struct {
	Coords Coordinates
	Rng    Range
}

func (e *Entry) Write(f *os.File) error {
	if err := binary.Write(f, binary.LittleEndian, uint8(e.Coords.Z)); err != nil {
		return err
	}

	if _, err := f.Write(writeUint24(e.Coords.X)); err != nil {
		return err
	}

	if _, err := f.Write(writeUint24(e.Coords.Y)); err != nil {
		return err
	}

	if _, err := f.Write(writeUint48(e.Rng.Offset)); err != nil {
		return err
	}

	if err := binary.Write(f, binary.LittleEndian, uint32(e.Rng.Length)); err != nil {
		return err
	}

	return nil
}

type EntryPointers map[Coordinates]Range
type EntryAscending []Entry

func (e EntryAscending) Len() int {
	return len(e)
}

func (e EntryAscending) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e EntryAscending) Less(i, j int) bool {
	if e[i].Coords.Z != e[j].Coords.Z {
		return e[i].Coords.Z < e[j].Coords.Z
	}

	if e[i].Coords.X != e[j].Coords.X {
		return e[i].Coords.X < e[j].Coords.X
	}

	return e[i].Coords.Y < e[j].Coords.Y
}

func (e EntryAscending) Write(f *os.File) error {
	for _, entry := range e {
		if err := entry.Write(f); err != nil {
			return err
		}
	}

	return nil
}

type Directory struct {
	Z       uint8
	Entries EntryPointers
	Leaves  EntryPointers
}

func (d *Directory) SizeBytes() int {
	return 21*(len(d.Entries)+len(d.Leaves)) + 1
}

func (d *Directory) Parse(dir_bytes []byte) {
	var max uint8
	for i := 0; i < len(dir_bytes)/17; i++ {
		leaf, coords, rng := ParseEntry(dir_bytes[i*17 : i*17+17])

		if leaf == 0 {
			d.Entries[coords] = rng
		} else {
			max = leaf
			d.Leaves[coords] = rng
		}
	}
	d.Z = max
}

func ParseEntry(b []byte) (uint8, Coordinates, Range) {
	z := uint8(b[0])
	x := uint32(readUint24(b[1:4]))
	y := uint32(readUint24(b[4:7]))
	leaf := uint8(0)

	if b[0]&0b10000000 != 0 {
		leaf = b[0] & 0b01111111
		z = leaf
	}

	coords := Coordinates{Z: z, X: x, Y: y}
	rng := Range{Offset: readUint48(b[7:13]), Length: binary.LittleEndian.Uint32(b[13:17])}

	return leaf, coords, rng
}
