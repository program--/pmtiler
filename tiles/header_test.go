package tiles

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHeader(t *testing.T) {
	var header Header
	f, err := os.Open(filepath.Join("..", "testdata", "tiles.pmtiles"))
	if err != nil {
		t.Fatal("failed to read tiles.pmtiles", err)
	}

	if err := header.Parse(f, true); err != nil {
		t.Error("failed to parse tiles.pmtiles", err)
	}

	if header.Version != 2 {
		t.Fail()
	}

	if header.RootDirLength != 26 {
		t.Fail()
	}

	if header.HeaderMetadata.Format != "pbf" {
		t.Fail()
	}

	if header.HeaderMetadata.Maxzoom != "14" {
		t.Fail()
	}

	if header.HeaderMetadata.Minzoom != "0" {
		t.Fail()
	}

	if header.HeaderMetadata.Bounds != "-121.5111923,38.5617238,-121.4681053,38.5815868" {
		t.Fail()
	}
}

func TestWriteHeader(t *testing.T) {
	var header Header
	golden, err := os.Open(filepath.Join("..", "testdata", "tiles.pmtiles"))
	if err != nil {
		t.Fatal("failed to read tiles.pmtiles", err)
	}

	if err := header.Parse(golden, true); err != nil {
		t.Error("failed to parse tiles.pmtiles", err)
	}

	testfile, err := os.Create("../testdata/test.pmtiles")
	if err != nil {
		t.Fatal("failed to create testfile", err)
	}

	if err := header.Write(testfile); err != nil {
		t.Error("failed to write tiles to testfile", err)
	}

	testfile.Sync()
	testfile.Close()
	testfile, _ = os.Open("../testdata/test.pmtiles")

	var testheader Header
	if err := testheader.Parse(testfile, true); err != nil {
		t.Error("failed to parse testfile", err)
	}

	defer testfile.Close()
	defer os.Remove(testfile.Name())

	if header.Version != testheader.Version {
		t.Fail()
	}

	if header.RootDirLength != testheader.RootDirLength {
		t.Fail()
	}

	if header.HeaderMetadata.Format != testheader.HeaderMetadata.Format {
		t.Fail()
	}

	if header.HeaderMetadata.Maxzoom != testheader.HeaderMetadata.Maxzoom {
		t.Fail()
	}

	if header.HeaderMetadata.Minzoom != testheader.HeaderMetadata.Minzoom {
		t.Fail()
	}

	if header.HeaderMetadata.Bounds != testheader.HeaderMetadata.Bounds {
		t.Fail()
	}
}
