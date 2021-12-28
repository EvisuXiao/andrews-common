package utils

import "math"

func Round(num float64, scale int) float64 {
	p := math.Pow10(scale)
	return math.Round(num*p) / p
}

func RoundInt(num float64) int {
	return int(math.Round(num))
}

func Ceil(num float64, scale int) float64 {
	p := math.Pow10(scale)
	return math.Ceil(num*p) / p
}

func CeilInt(num float64) int {
	return int(math.Ceil(num))
}

func Floor(num float64, scale int) float64 {
	p := math.Pow10(scale)
	return math.Floor(num*p) / p
}

func FloorInt(num float64) int {
	return int(math.Floor(num))
}
