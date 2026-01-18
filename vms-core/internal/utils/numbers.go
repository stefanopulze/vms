package utils

import "strconv"

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
