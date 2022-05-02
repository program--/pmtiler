package geojson

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/mvt"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	"github.com/paulmach/orb/maptile/tilecover"

	tiles "partiles/tiles"
)

func GeoJSONToTiles(path string, fc *geojson.FeatureCollection, layer_name string, max_zoom maptile.Zoom) error {
	var err error = nil
	writer := tiles.NewWriter(path)

	for z := maptile.Zoom(0); z <= max_zoom; z++ {
		covering := tilecover.Bound(fc.BBox.Bound(), z)

		for tile := range covering {
			bound := tile.Bound(0.5)

			fcclone := geojson.NewFeatureCollection()
			for _, g := range fc.Features {
				if bound.Contains(g.Geometry.(orb.Point)) {
					newg := new(geojson.Feature)
					*newg = *g
					fcclone.Append(newg)
				}
			}

			current_layer := mvt.NewLayers(map[string]*geojson.FeatureCollection{layer_name: fcclone})
			current_layer.ProjectToTile(tile)
			tileb, err := mvt.Marshal(current_layer)
			if err != nil {
				return err
			}

			writer.WriteTile(tiles.Coordinates{Z: uint8(tile.Z), X: tile.X, Y: tile.Y}, tileb)
		}
	}

	err = writer.Finalize(&tiles.Metadata{
		Bounds:  boundToString(fc.BBox.Bound()),
		Format:  "pbf",
		Maxzoom: fmt.Sprint(max_zoom),
		Minzoom: "0",
	})

	return err
}

func wgs84ToXY(lat float64, lon float64, z maptile.Zoom) (uint32, uint32) {
	z_pow := float64(uint64(1) << z)
	lat_pi := lat * math.Pi / 180
	lat_log := math.Log(math.Tan(lat_pi) + 1/math.Cos(lat_pi))

	x := uint32(math.Floor((lon + 180) / 360 * z_pow))
	y := uint32(math.Floor((1 - lat_log/math.Pi) / 2 * z_pow))

	return x, y
}

func floatToString(x float64) string {
	return strconv.FormatFloat(x, 'f', -1, 64)
}

func floatSliceToString(x []float64) []string {
	s := make([]string, len(x))
	for i, v := range x {
		s[i] = floatToString(v)
	}
	return s
}

func boundToString(b orb.Bound) string {
	return strings.Join(floatSliceToString([]float64{b.Left(), b.Bottom(), b.Right(), b.Top()}), ",")
}
