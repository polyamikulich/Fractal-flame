package variations

import "math"

func Linear(x, y float64) (float64, float64) {
	return x, y
}

func Sinusoidal(x, y float64) (float64, float64) {
	new_x := math.Sin(x)
	new_y := math.Sin(y)
	return new_x, new_y
}

func Spherical(x, y float64) (float64, float64) {
	r2 := x*x + y*y
	if r2 == 0 {
		return 0, 0
	}
	invR2 := 1.0 / r2
	return x * invR2, y * invR2
}

func Swirl(x, y float64) (float64, float64) {
	r := r(x, y)
	new_x := x*math.Sin(r*r) - y*math.Cos(r*r)
	new_y := x*math.Cos(r*r) + y*math.Sin(r*r)
	return new_x, new_y
}

func Horseshoe(x, y float64) (float64, float64) {
	r := r(x, y)
	if r == 0 {
		return 0, 0
	}
	new_x := (x - y) * (x + y) / r
	new_y := 2 * x * y / r
	return new_x, new_y
}

func Polar(x, y float64) (float64, float64) {
	if y == 0 {
		return 0, 0
	}
	theta := theta(x, y)
	r := r(x, y)
	new_x := theta / math.Pi
	new_y := r - 1.0
	return new_x, new_y
}

func Handkerchief(x, y float64) (float64, float64) {
	r := r(x, y)
	if y == 0 {
		return 0, 0
	}
	theta := theta(x, y)
	new_x := r * (math.Sin(theta + r))
	new_y := r * (math.Cos(theta - r))
	return new_x, new_y
}

func Heart(x, y float64) (float64, float64) {
	r := r(x, y)
	if y == 0 {
		return 0, 0
	}
	theta := theta(x, y)
	new_x := r * (math.Sin(theta * r))
	new_y := r * (math.Cos(theta * r)) * (-1.0)
	return new_x, new_y
}

func Disc(x, y float64) (float64, float64) {
	r := r(x, y)
	if y == 0 {
		return 0, 0
	}
	theta := theta(x, y)
	new_x := theta * math.Sin(r*math.Pi) / math.Pi
	new_y := theta * math.Cos(r*math.Pi) / math.Pi
	return new_x, new_y
}

func Spiral(x, y float64) (float64, float64) {
	r := r(x, y)
	if r == 0 || y == 0 {
		return 0, 0
	}
	theta := theta(x, y)
	new_x := (math.Cos(theta) + math.Sin(r)) / r
	new_y := (math.Sin(theta) - math.Cos(r)) / r
	return new_x, new_y
}
