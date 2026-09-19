package api

import "math"

const nearZero = 1e-12

// TrigonometryResult contains trigonometric values for an angle in degrees.
// Nil reciprocal values represent mathematically undefined results.
type TrigonometryResult struct {
	Sin   float64  `json:"sin"`
	Cos   float64  `json:"cos"`
	Tan   *float64 `json:"tan"`
	Cosec *float64 `json:"cosec"`
	Sec   *float64 `json:"sec"`
	Cot   *float64 `json:"cot"`
}

// Trigonometry evaluates standard and reciprocal trigonometric functions for
// theta degrees. Values within nearZero of zero are normalized to zero so
// cardinal angles do not expose floating-point artifacts.
func Trigonometry(theta float64) TrigonometryResult {
	radians := theta * math.Pi / 180
	sin := normalizeNearZero(math.Sin(radians))
	cos := normalizeNearZero(math.Cos(radians))

	result := TrigonometryResult{Sin: sin, Cos: cos}
	if cos != 0 {
		tan := sin / cos
		sec := 1 / cos
		result.Tan = &tan
		result.Sec = &sec
	}
	if sin != 0 {
		cosec := 1 / sin
		cot := cos / sin
		result.Cosec = &cosec
		result.Cot = &cot
	}
	return result
}

func normalizeNearZero(value float64) float64 {
	if math.Abs(value) < nearZero {
		return 0
	}
	return value
}
