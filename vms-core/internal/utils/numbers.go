package utils

import (
	"math"
	"strconv"
)

func ParseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func ParseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func ParseFloatAsInt16(s string) uint16 {
	return uint16(ParseFloat(s))
}

func ConvertAndRoundWatt(f any) int32 {
	switch v := f.(type) {
	case float64:
		return int32(math.Round(v))
	case int64:
		return int32(v)
	case int32:
		return v
	default:
		return 0
	}
}

func NonZero(v int32) int32 {
	if v < 0 {
		return 0
	}

	return v
}
