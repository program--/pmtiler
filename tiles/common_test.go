package tiles

import (
	"bytes"
	"testing"
)

// Attribution:
// https://github.com/protomaps/go-pmtiles/blob/master/pmtiles/writer_test.go
// https://github.com/protomaps/go-pmtiles/blob/master/pmtiles/reader_test.go

func TestReadUint24(t *testing.T) {
	b := []byte{255, 255, 255}
	result := readUint24(b)
	if result != 16777215 {
		t.Errorf("result did not match, was %d", result)
	}
	b = []byte{255, 0, 0}
	result = readUint24(b)
	if result != 255 {
		t.Errorf("result did not match, was %d", result)
	}
}

func TestWriteUint24(t *testing.T) {
	var i uint32
	i = 16777215
	result := writeUint24(i)
	if !bytes.Equal(result, []byte{255, 255, 255}) {
		t.Errorf("result did not match, was %d", result)
	}
	i = 255
	result = writeUint24(i)
	if !bytes.Equal(result, []byte{255, 0, 0}) {
		t.Errorf("result did not match, was %d", result)
	}
}

func TestReadUint48(t *testing.T) {
	b := []byte{255, 255, 255, 255, 255, 255}
	result := readUint48(b)
	if result != 281474976710655 {
		t.Errorf("result did not match, was %d", result)
	}
	b = []byte{255, 0, 0, 0, 0, 0}
	result = readUint48(b)
	if result != 255 {
		t.Errorf("result did not match, was %d", result)
	}
}

func TestWriteUint48(t *testing.T) {
	var i uint64
	i = 281474976710655
	result := writeUint48(i)
	if !bytes.Equal(result, []byte{255, 255, 255, 255, 255, 255}) {
		t.Errorf("result did not match, was %d", result)
	}
	i = 255
	result = writeUint48(i)
	if !bytes.Equal(result, []byte{255, 0, 0, 0, 0, 0}) {
		t.Errorf("result did not match, was %d", result)
	}
}
