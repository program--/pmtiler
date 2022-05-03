package io

import (
	"path/filepath"
	"testing"
)

func TestGDALRead(t *testing.T) {
	gdal_file := filepath.Join("..", "testdata", "sample.gpkg")
	x, err := GDALFile(gdal_file, 0)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(
		"Num Features: %d",
		len(x.Features),
	)
	if len(x.Features) != 147 {
		t.Fail()
	}

	t.Logf(
		"Num Fields: %d",
		len(x.Features[0].Properties),
	)

	t.Log(x.Features[0].Properties)
	if len(x.Features[0].Properties) != 33 {
		t.Fail()
	}
}
