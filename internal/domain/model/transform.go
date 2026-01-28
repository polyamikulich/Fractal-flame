package model

import "math/rand"

// FlameTransform - аф.преобразование + его цвета
type AffineTransform struct {
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	C      float64 `json:"c"`
	D      float64 `json:"d"`
	E      float64 `json:"e"`
	F      float64 `json:"f"`
	ColorR float64 `json:"-"`
	ColorG float64 `json:"-"`
	ColorB float64 `json:"-"`
}

// NewFlameTransform создаёт новый FlameTransform с заданным аффинным преобразованием
// и генерирует случайный RGB цвет в диапазоне [0, 1].
func NewAffineTransform(affine AffineTransform, seed int64) AffineTransform {
	// Используем seed для детерминированной генерации цветов
	rng := rand.New(rand.NewSource(seed))

	return AffineTransform{
		A:      affine.A,
		B:      affine.B,
		C:      affine.C,
		D:      affine.D,
		E:      affine.E,
		F:      affine.F,
		ColorR: rng.Float64(),
		ColorG: rng.Float64(),
		ColorB: rng.Float64(),
	}
}
