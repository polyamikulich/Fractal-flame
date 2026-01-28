package variations

import (
	"math"
	"testing"
)

func TestLinear(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "positive input",
			x:    1.0, y: 2.0,
			expX: 1.0, expY: 2.0,
		},
		{
			name: "negative input",
			x:    -1.0, y: -2.0,
			expX: -1.0, expY: -2.0,
		},
		{
			name: "mixed input 1",
			x:    3.5, y: -2.1,
			expX: 3.5, expY: -2.1,
		},
		{
			name: "mixed input 2",
			x:    -3.5, y: 2.1,
			expX: -3.5, expY: 2.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Linear(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-9 || math.Abs(gotY-tt.expY) > 1e-9 {
				t.Errorf("Linear(%f, %f) = (%f, %f), want (%f, %f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestSinusoidal(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0, // sin(0) = 0
		},
		{
			name: "positive input",
			x:    1.0, y: 1.0,
			expX: 0.8414709848078975, expY: 0.8414709848078975, // sin(1) = 0.8415
		},
		{
			name: "negative input",
			x:    -1.0, y: -1.0,
			expX: -0.8414709848078975, expY: -0.8414709848078975, // sin(-1) = -0.8415
		},
		{
			name: "large input",
			x:    100.0, y: -100.0,
			expX: -0.5063656411, expY: 0.50636564115, // sin(100) = 0.5063
		},
		{
			name: "small input",
			x:    0.1, y: 0.1,
			expX: 0.09983341664682805, expY: 0.09983341664682805, // sin(0.1) = 0.0998
		},
		{
			name: "pi input",
			x:    math.Pi, y: -math.Pi,
			expX: 0.0, expY: 0.0, // sin(+-pi) = 0
		},
		{
			name: "pi/2 input",
			x:    math.Pi / 2, y: -math.Pi / 2,
			expX: 1.0, expY: -1.0, // sin(pi/2) = 1, sin(-pi/2) = -1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Sinusoidal(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-9 || math.Abs(gotY-tt.expY) > 1e-9 {
				t.Errorf("Sinusoidal(%.3f, %.3f) = (%.3f, %.3f), want (%.3f, %.3f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestSpherical(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    1.0, y: 0.0,
			expX: 1.0, expY: 0.0,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: 0.0, expY: -1.0,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: 0.5, expY: -0.5,
		},
		{
			name: "large input",
			x:    100.0, y: -100.0,
			expX: 0.005, expY: -0.005,
		},
		{
			name: "small input",
			x:    0.05, y: -0.05,
			expX: 10.0, expY: -10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Spherical(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-9 || math.Abs(gotY-tt.expY) > 1e-9 {
				t.Errorf("Spherical(%.3f, %.3f) = (%.3f, %.3f), want (%.3f, %.3f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestSwirl(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    1.0, y: 0.0,
			expX: 0.841470985, expY: 0.540302306,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: 0.540302306, expY: -0.841470985,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: 0.49315059, expY: -1.32544426,
		},
		{
			name: "large input",
			x:    100.0, y: -100.0,
			expX: 139.51844526, expY: 23.12149286,
		},
		{
			name: "small input",
			x:    0.05, y: -0.05,
			expX: 0.05024937, expY: 0.04974938,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Swirl(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-7 || math.Abs(gotY-tt.expY) > 1e-7 {
				t.Errorf("Swirl(%.3f, %.3f) = (%.10f, %.10f), want (%.10f, %.10f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestHorseshoe(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    1.0, y: 0.0,
			expX: 1.0, expY: 0.0,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: -1.0, expY: 0.0,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: 0.0, expY: -1.41421356,
		},
		{
			name: "large input",
			x:    100.0, y: -100.0,
			expX: 0.0, expY: -141.42135624,
		},
		{
			name: "small input",
			x:    0.05, y: -0.05,
			expX: 0.0, expY: -0.07071068,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Horseshoe(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-7 || math.Abs(gotY-tt.expY) > 1e-7 {
				t.Errorf("Horseshoe(%.3f, %.3f) = (%.8f, %.8f), want (%.8f, %.8f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestPolar(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    0.0, y: 1.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: -0.25, expY: 0.41421356,
		},
		{
			name: "large input",
			x:    100.0, y: -100.0,
			expX: -0.25, expY: 140.42135624,
		},
		{
			name: "small input",
			x:    0.08, y: -0.05,
			expX: -0.32219232, expY: -0.90566019,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Polar(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-7 || math.Abs(gotY-tt.expY) > 1e-7 {
				t.Errorf("Polar(%.3f, %.3f) = (%.8f, %.8f), want (%.8f, %.8f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestHandkerchief(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    0.0, y: 1.0,
			expX: 0.84147098, expY: 0.54030231,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: 0.84147098, expY: 0.54030231,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: 0.83182225, expY: -0.83182225,
		},
		{
			name: "large input",
			x:    100.0, y: -200.0,
			expX: -20.11343028, expY: -117.52947121,
		},
		{
			name: "small input",
			x:    0.08, y: -0.05,
			expX: -0.07493427, expY: 0.04224167,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Handkerchief(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-7 || math.Abs(gotY-tt.expY) > 1e-7 {
				t.Errorf("Handkerchief(%.3f, %.3f) = (%.8f, %.8f), want (%.8f, %.8f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestHeart(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    0.0, y: 1.0,
			expX: 0.0, expY: -1.0,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: 0.0, expY: -1.0,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: -1.26716213, expY: -0.62793322,
		},
		{
			name: "large input",
			x:    100.0, y: -200.0,
			expX: 0.49183721, expY: 223.60625684,
		},
		{
			name: "small input",
			x:    0.08, y: -0.05,
			expX: -0.00899487, expY: -0.09391002,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Heart(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-7 || math.Abs(gotY-tt.expY) > 1e-7 {
				t.Errorf("Heart(%.3f, %.3f) = (%.8f, %.8f), want (%.8f, %.8f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestDisc(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    0.0, y: 1.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: 0.24097563, expY: 0.06656384,
		},
		{
			name: "large input",
			x:    100.0, y: -200.0,
			expX: 0.13935448, expY: -0.04859272,
		},
		{
			name: "small input",
			x:    0.08, y: -0.05,
			expX: -0.09409863, expY: -0.30814499,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Disc(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-7 || math.Abs(gotY-tt.expY) > 1e-7 {
				t.Errorf("Disc(%.3f, %.3f) = (%.8f, %.8f), want (%.8f, %.8f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}

func TestSpiral(t *testing.T) {
	tests := []struct {
		name       string
		x, y       float64
		expX, expY float64
	}{
		{
			name: "zero input",
			x:    0.0, y: 0.0,
			expX: 0.0, expY: 0.0,
		},
		{
			name: "pos + null input",
			x:    0.0, y: 1.0,
			expX: 1.84147098, expY: -0.54030231,
		},
		{
			name: "neg + null input",
			x:    0.0, y: -1.0,
			expX: 1.84147098, expY: -0.54030231,
		},
		{
			name: "pos + neg input",
			x:    1.0, y: -1.0,
			expX: 1.19845600, expY: -0.61026884,
		},
		{
			name: "large input",
			x:    100.0, y: -200.0,
			expX: 0.00164831, expY: 0.00180389,
		},
		{
			name: "small input",
			x:    0.08, y: -0.05,
			expX: 6.61649485, expY: -19.54160791,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := Spiral(tt.x, tt.y)

			if math.Abs(gotX-tt.expX) > 1e-7 || math.Abs(gotY-tt.expY) > 1e-7 {
				t.Errorf("Spiral(%.3f, %.3f) = (%.8f, %.8f), want (%.8f, %.8f)",
					tt.x, tt.y, gotX, gotY, tt.expX, tt.expY)
			}
		})
	}
}
