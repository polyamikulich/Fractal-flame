package variations

import (
	"fmt"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// registry - лист регистрации (на самом деле мапа)
var registry = map[string]model.VariationFunc{
	"linear":       Linear,
	"sinusoidal":   Sinusoidal,
	"spherical":    Spherical,
	"swirl":        Swirl,
	"horseshoe":    Horseshoe,
	"polar":        Polar,
	"handkerchief": Handkerchief,
	"heart":        Heart,
	"disc":         Disc,
	"spiral":       Spiral,
}

// GetVariation - получение функции вариации по имени из листа регистрации
func GetVariation(name string) (model.VariationFunc, error) {
	fn, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown variation: %s", name)
	}
	return fn, nil
}
