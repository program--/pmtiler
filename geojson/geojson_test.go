package geojson

import (
	"path/filepath"
	ptio "pmtiler/io"
	"testing"

	"github.com/paulmach/orb/maptile"
)

func TestWritingPMTiles(t *testing.T) {
	tiles_path := filepath.Join("..", "testdata", "sample.pmtiles")
	parquet_path := "s3://example/data/path.parquet"
	if parquet_path == "s3://example/data/path.parquet" {
		t.Skip("S3 file path needs to be modified")
	}

	fc, err := ptio.ParquetToGeoJSON(parquet_path, "X", "Y")
	if err != nil {
		t.Fatal("failed to parse Parquet to GeoJSON:", err)
	}
	err = GeoJSONToTiles(tiles_path, fc, "test_layer", maptile.Zoom(20))
	if err != nil {
		t.Fatal("failed to write GeoJSON to PMTiles:", err)
	}
}
