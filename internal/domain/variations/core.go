package variations

import "math"

func r(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func theta(x, y float64) float64 {
	return math.Atan(x / y)
}

func phi(x, y float64) float64 {
	return math.Atan(y / x)
}
