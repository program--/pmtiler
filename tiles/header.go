package tiles

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

const HEADER_SIZE = 512000

var MAGIC_NUMBER = []byte("PM")
var ErrInvalidMagic error = errors.New("read error; invalid magic number")
var ErrInvalidVersion error = errors.New("read error; invalid version number")

type Metadata struct {
	Bounds      string `json:"bounds"`
	Format      string `json:"format"`
	Maxzoom     string `json:"maxzoom"`
	Minzoom     string `json:"minzoom"`
	Compress    string `json:"compress,omitempty"`
	Attribution string `json:"attribution,omitempty"`
}

type Header struct {
	Version         uint16
	MetadataLength  uint32
	RootDirLength   uint16
	HeaderMetadata  *Metadata
	HeaderDirectory *Directory
}

// Read

func (h *Header) Parse(r io.Reader, metaOnly bool) error {
	if err := h.parseMagic(r); err != nil {
		return fmt.Errorf("[parseMagic] %w", err)
	}

	if err := h.parseVersion(r); err != nil {
		return fmt.Errorf("[parseVersion] %w", err)
	}

	if err := h.parseLengths(r); err != nil {
		return fmt.Errorf("[parseLengths] %w", err)
	}

	if err := h.parseMetadata(r); err != nil {
		return fmt.Errorf("[parseMetadata] %w", err)
	}

	if !metaOnly {
		if err := h.parseRootDirectory(r); err != nil {
			fmt.Println(fmt.Errorf("[parseRootDirectory] %w", err))
			return err
		}
	}

	return nil
}

func (h *Header) parseMagic(r io.Reader) error {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err == nil {
		for i, v := range buf {
			if v != MAGIC_NUMBER[i] {
				// TODO:
				//   Show read magic number and
				//   correct magic number in error
				return ErrInvalidMagic
			}
		}
	}

	return err
}

func (h *Header) parseVersion(r io.Reader) error {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf)
	if err == nil {
		h.Version = binary.LittleEndian.Uint16(buf)
		if h.Version != 2 {
			return ErrInvalidVersion
		}
	}
	return err
}

func (h *Header) parseLengths(r io.Reader) error {
	mbuf := make([]byte, 4)
	_, err := io.ReadFull(r, mbuf)
	if err == nil {
		h.MetadataLength = binary.LittleEndian.Uint32(mbuf)
	}

	rbuf := make([]byte, 2)
	_, err = io.ReadFull(r, rbuf)
	if err == nil {
		h.RootDirLength = binary.LittleEndian.Uint16(rbuf)
	}

	return err
}

func (h *Header) parseMetadata(r io.Reader) error {
	buf := make([]byte, h.MetadataLength)
	_, err := io.ReadFull(r, buf)
	if err == nil {
		err = json.Unmarshal(buf, &(h.HeaderMetadata))
	}

	return err
}

func (h *Header) parseRootDirectory(r io.Reader) error {
	buf := make([]byte, h.RootDirLength*17)
	_, err := io.ReadFull(r, buf)
	if err == nil {
		h.HeaderDirectory = new(Directory)
		h.HeaderDirectory.Z = 0
		h.HeaderDirectory.Entries = make(EntryPointers)
		h.HeaderDirectory.Leaves = make(EntryPointers)
		h.HeaderDirectory.Parse(buf)
	}

	return err
}

// Write

func (h *Header) Write(file *os.File) error {
	if _, err := file.Write(MAGIC_NUMBER); err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, h.Version); err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, h.MetadataLength); err != nil {
		return err
	}

	if err := binary.Write(file, binary.LittleEndian, h.RootDirLength); err != nil {
		return err
	}

	metadata, err := json.Marshal(h.HeaderMetadata)
	if err == nil {
		_, err = file.Write(metadata)
	}

	return err
}
