package io

import (
	"path/filepath"
	"testing"
)

func TestLocalParquetToGeoJSON(t *testing.T) {
	path := filepath.Join("..", "testdata", "sample.parquet")
	fc, err := ParquetToGeoJSON(path, "X", "Y")
	if err != nil {
		t.Fatal("failed to parse Parquet to GeoJSON:", err)
	}

	if len(fc.Features) != 6 {
		t.Fail()
	}

	if len(fc.Features[0].Properties) != 3 {
		t.Fail()
	}
}

func TestS3ParquetToGeoJSON(t *testing.T) {
	path := "s3://example/data/path.parquet"
	fc, err := ParquetToGeoJSON(path, "X", "Y")
	if err != nil {
		t.Fatal("failed to parse Parquet to GeoJSON:", err)
	}

	if len(fc.Features) != 20527 {
		t.Fail()
	}
}
