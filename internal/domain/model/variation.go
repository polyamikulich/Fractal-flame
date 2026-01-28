package model

// Variation - структура для хранения информации о вариации
type Variation struct {
	Name   string        `json:"name"`
	Weight float64       `json:"weight"`
	Apply  VariationFunc `json:"-"`
}

type VariationFunc func(x, y float64) (float64, float64)
