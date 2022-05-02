package tiles

import (
	"math"
)

func GetParentTile(tile Coordinates, level uint8) Coordinates {
	diff := tile.Z - level
	x := math.Floor(float64(tile.X) / float64(int(1)<<diff))
	y := math.Floor(float64(tile.Y) / float64(int(1)<<diff))
	return Coordinates{Z: level, X: uint32(x), Y: uint32(y)}
}
