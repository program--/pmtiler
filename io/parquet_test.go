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
		t.Logf("Num Features: %d", len(fc.Features))
		t.Fail()
	}

	if len(fc.Features[0].Properties) != 1 {
		t.Logf("Num Properties: %d", len(fc.Features[0].Properties))
		t.Logf("%v", fc.Features[0].Properties)
		t.Fail()
	}
}

func TestS3ParquetToGeoJSON(t *testing.T) {
	path := "s3://example/data/path.parquet"

	if path == "s3://example/data/path.parquet" {
		t.Skip("S3 file path needs to be modified")
	}

	fc, err := ParquetToGeoJSON(path, "X", "Y")
	if err != nil {
		t.Fatal("failed to parse Parquet to GeoJSON:", err)
	}

	if len(fc.Features) != 20527 {
		t.Fail()
	}
}
