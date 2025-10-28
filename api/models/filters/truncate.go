package filters

import "math"

func TruncateFloat(x float64, decimals int) float64 {
	pow := math.Pow(10, float64(decimals))
	return math.Trunc(x*pow) / pow
}
