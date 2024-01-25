package util

import "math"

func Round(val float64, precision int) float64 {
	return math.Round(val*float64(precision)) / float64(precision)
}
