package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const sizeStep = 1024

var sizeUnits = []string{"B", "K", "M", "G", "T"}

func CalcSizeBytes(size string) uint64 {
	if size == "" {
		return 0
	}
	numStr := size[:len(size)-1]
	num, err := strconv.ParseFloat(numStr, 32)
	if HasErr(err) {
		return 0
	}
	unit := strings.ToUpper(size[len(size)-1:])
	unitIdx := SearchSliceIndex(unit, sizeUnits)
	if unitIdx < 0 {
		return 0
	}
	return uint64(math.Pow(num*sizeStep, float64(unitIdx)))
}

func CalcBytesSize(bytes uint64) string {
	if bytes < 1 {
		return "0"
	}
	num := float64(bytes)
	var unit string
	for _, unit = range sizeUnits {
		if num < sizeStep {
			break
		}
		num /= sizeStep
	}
	return fmt.Sprintf("%v%s", Round(num, 3), unit)
}
