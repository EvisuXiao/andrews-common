package utils

import "math"

func CeilInt(divider int, dividend int) int {
	return int(math.Ceil(float64(divider / dividend)))
}

func FloorInt(divider int, dividend int) int {
	return int(math.Floor(float64(divider / dividend)))
}
