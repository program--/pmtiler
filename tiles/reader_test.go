package tiles

import "testing"

func TestGetParentTile(t *testing.T) {
	a := Coordinates{Z: 8, X: 125, Y: 69}
	result := GetParentTile(a, 7)
	if (result != Coordinates{Z: 7, X: 62, Y: 34}) {
		t.Errorf("result did not match, was %d", result)
	}
}
