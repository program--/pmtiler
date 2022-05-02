package io

import (
	"context"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go-source/s3"
	"github.com/xitongsys/parquet-go/reader"
	"github.com/xitongsys/parquet-go/source"
)

func LocalParquetFile(path string) (source.ParquetFile, error) {
	return local.NewLocalFileReader(path)
}

func S3ParquetFile(ctx context.Context, bucket string, key string) (source.ParquetFile, error) {
	return s3.NewS3FileReader(ctx, bucket, key, aws.NewConfig().WithRegion("us-east-1"))
}

func CreateParquetReader(file source.ParquetFile, schema interface{}) (*reader.ParquetReader, error) {
	return reader.NewParquetReader(file, schema, 10)
}

func ParquetToGeoJSON(path string, x string, y string) (*geojson.FeatureCollection, error) {
	var (
		fr  source.ParquetFile
		err error
	)

	// Setup ParquetFile and ParquetReader
	if uri, _ := url.Parse(path); uri.Scheme == "s3" {
		ctx := context.Background()
		fr, err = S3ParquetFile(ctx, uri.Host, strings.TrimLeft(uri.Path, "/"))
	} else {
		fr, err = LocalParquetFile(path)
	}

	if err != nil {
		return nil, err
	}

	defer fr.Close()

	pr, err := reader.NewParquetColumnReader(fr, 4)
	if err != nil {
		return nil, err
	}

	// Read in parquet columns to map
	rows := int64(pr.GetNumRows())
	interval := rows
	schema := pr.SchemaHandler.IndexMap
	delete(schema, 0)

	data := make(map[string][]interface{})
	for start := int64(0); start < rows; start += interval {
		end := start + interval
		if end > rows {
			end = rows
		}

		for _, v := range schema {
			key := strings.Replace(v, "Schema\u0001", "", 1)
			column, _, _, err := pr.ReadColumnByPath(v, interval)
			if err != nil {
				return nil, err
			}

			data[key] = append(data[key], column...)
		}
	}

	features := make([]*geojson.Feature, rows)
	xmax := data[x][0].(float64)
	xmin := data[x][0].(float64)
	ymax := data[y][0].(float64)
	ymin := data[y][0].(float64)
	for i := int64(0); i < rows; i++ {
		xcurr := data[x][i].(float64)
		ycurr := data[y][i].(float64)
		features[i] = geojson.NewFeature(orb.Point{xcurr, ycurr})

		if xmax < xcurr {
			xmax = xcurr
		}

		if ymax < ycurr {
			ymax = ycurr
		}

		if xmin > xcurr {
			xmin = xcurr
		}

		if ymin > ycurr {
			ymin = ycurr
		}

		for _, v := range schema {
			key := strings.Replace(v, "Schema\u0001", "", 1)

			if key != x && key != y {
				features[i].Properties[key] = data[key][i]
			}
		}
	}

	fc := geojson.NewFeatureCollection()
	fc.Features = features
	fc.BBox = geojson.NewBBox(orb.Bound{Min: orb.Point{xmin, ymin}, Max: orb.Point{xmax, ymax}})
	return fc, nil
}
