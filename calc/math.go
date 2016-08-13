package calc

import (
	"fmt"
	"math"
)

const rads = 3960
const dToRads = math.Pi / 180.0
const radsToD = 180.0 / math.Pi

func GetLatDelta(dist float64) float64 {

	return float64(dist/rads) * dToRads
}

func GetLongDelta(lat float64, dist float64) float64 {

	r := rads * math.Cos(lat*dToRads)
	return (float64(dist) / r) * radsToD
}

func Intersect(lats []int64, longs []int64) []int64 {
	fmt.Print("hi\n")
	m := make(map[int64]bool)
	vendors := make([]int64, 0)
	for _, k := range lats {
		m[k] = true
	}

	for _, k := range longs {
		if m[k] {
			vendors = append(vendors, k)
		}
	}

	return vendors
}
