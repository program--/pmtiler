package io

import (
	"github.com/lukeroth/gdal"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/geojson"
)

type GDALSource struct {
	ds *gdal.DataSource
}

func (gs *GDALSource) LayerDefinition(index int) gdal.FeatureDefinition {
	return gs.ds.LayerByIndex(index).Definition()
}

func (gs *GDALSource) ToGeoJSON(layer int) (*geojson.FeatureCollection, error) {
	// Getting Fields Information
	def := gs.LayerDefinition(layer)
	nfields := def.FieldCount()
	fieldDefs := []gdal.FieldDefinition{}
	for i := 0; i < nfields; i++ {
		fieldDef := def.FieldDefinition(i)
		if !fieldDef.IsIgnored() {
			fieldDefs = append(fieldDefs, fieldDef)
		}
	}

	// Getting Features
	fc := geojson.NewFeatureCollection()
	wgs84 := gdal.CreateSpatialReference("")
	wgs84.FromEPSG(4326)

	gdalLayer := gs.ds.LayerByIndex(layer)
	feature := gdalLayer.NextFeature()
	for feature != nil {
		gdalGeometry := feature.Geometry()
		err := gdalGeometry.TransformTo(wgs84)
		if err != nil {
			return nil, err
		}

		fwkb, err := gdalGeometry.ToWKB()
		if err != nil {
			return nil, err
		}

		geom, err := wkb.Unmarshal(fwkb)
		if err != nil {
			return nil, err
		}

		gjs := geojson.NewFeature(geom)
		gjs.Properties = make(geojson.Properties)
		for i := 0; i < nfields; i++ {
			var ret any

			field_name := fieldDefs[i].Name()
			field_type := fieldDefs[i].Type()

			switch field_type {
			// Binary
			case gdal.FT_Binary:
				ret = feature.FieldAsBinary(i)
				break

			// DateTime
			case gdal.FT_Date:
				fallthrough
			case gdal.FT_Time:
				fallthrough
			case gdal.FT_DateTime:
				ret, _ = feature.FieldAsDateTime(i)
				break

			// Integers
			case gdal.FT_Integer:
				ret = feature.FieldAsInteger(i)
				break
			case gdal.FT_IntegerList:
				ret = feature.FieldAsIntegerList(i)
				break
			case gdal.FT_Integer64:
				ret = feature.FieldAsInteger64(i)
				break
			case gdal.FT_Integer64List:
				ret = feature.FieldAsInteger64List(i)
				break

			// Floats
			case gdal.FT_Real:
				ret = feature.FieldAsFloat64(i)
				break
			case gdal.FT_RealList:
				ret = feature.FieldAsFloat64List(i)
				break
			case gdal.FT_String:
				ret = feature.FieldAsString(i)
				break
			case gdal.FT_StringList:
				ret = feature.FieldAsStringList(i)
				break

			default:
				ret = -1
			}
			gjs.Properties[field_name] = ret
		}

		fc = fc.Append(gjs)
		feature = gdalLayer.NextFeature()
	}

	env, err := gdalLayer.Extent(false)
	if err != nil {
		return nil, err
	}

	ring := gdal.Create(gdal.GT_LinearRing)
	ring.AddPoint2D(env.MinX(), env.MinY())
	ring.AddPoint2D(env.MaxX(), env.MinY())
	ring.AddPoint2D(env.MaxX(), env.MaxY())
	ring.AddPoint2D(env.MinX(), env.MaxY())
	ring.AddPoint2D(env.MinX(), env.MinY())
	bbox := gdal.Create(gdal.GT_Polygon)
	bbox.AddGeometry(ring)
	bbox.SetSpatialReference(gdalLayer.SpatialReference())
	bbox.TransformTo(wgs84)
	env = bbox.Boundary().Envelope()

	fc.BBox = geojson.NewBBox(orb.Bound{
		Min: orb.Point{env.MinX(), env.MinY()},
		Max: orb.Point{env.MaxX(), env.MaxY()},
	})

	return fc, nil
}

func NewGDALSource(path string) *GDALSource {
	ds := gdal.OpenDataSource(path, int(gdal.ReadOnly))
	gs := GDALSource{ds: &ds}
	return &gs
}

func GDALFile(path string, layer int) (*geojson.FeatureCollection, error) {
	src := NewGDALSource(path)
	return src.ToGeoJSON(layer)
}
