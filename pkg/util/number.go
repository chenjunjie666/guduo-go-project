package util

import "math"

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
// float64，保留 p 位小数
func ToFixedFloat(f float64, p int) float64 {
	output := math.Pow(10, float64(p))
	return float64(round(f * output)) / output
}
