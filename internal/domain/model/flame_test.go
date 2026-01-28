package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlameConfig_SetDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    *FlameConfig
		expected *FlameConfig
	}{
		{
			name:  "all fields zero/nil",
			input: &FlameConfig{},
			expected: &FlameConfig{
				Size:           SizeConfig{Width: 1920, Height: 1080},
				Seed:           5,
				IterationCount: 2500,
				Output:         "result.png",
				Threads:        1,
				Gamma:          2.2,
				SymmetryLevel:  1,
				AffineParams: []AffineTransform{
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0}, // Пример цвета
					{A: 0.5, B: 0.0, C: -0.5, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.5, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.5},
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: -0.5},
				},
				Variations: []Variation{
					{Name: "linear", Weight: 1.0, Apply: nil},
				},
			},
		},
		{
			name: "some fields set",
			input: &FlameConfig{
				Size: SizeConfig{Width: 800},
				Seed: 123,
				// Остальные поля оставляем нулевыми
			},
			expected: &FlameConfig{
				Size:           SizeConfig{Width: 800, Height: 1080},
				Seed:           123,
				IterationCount: 2500,
				Output:         "result.png",
				Threads:        1,
				Gamma:          2.2,
				SymmetryLevel:  1,
				AffineParams: []AffineTransform{
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: -0.5, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.5, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.5},
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: -0.5},
				},
				Variations: []Variation{
					{Name: "linear", Weight: 1.0, Apply: nil},
				},
			},
		},
		{
			name: "all fields set",
			input: &FlameConfig{
				Size:           SizeConfig{Width: 800, Height: 600},
				Seed:           123,
				IterationCount: 10000,
				Output:         "res.png",
				Threads:        4,
				Gamma:          2,
				EnableGamma:    true,
				SymmetryLevel:  3,
				AffineParams: []AffineTransform{
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.5, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.5},
				},
				Variations: []Variation{
					{Name: "linear", Weight: 0.7, Apply: nil},
					{Name: "swirl", Weight: 0.9, Apply: nil},
				},
			},
			expected: &FlameConfig{
				Size:           SizeConfig{Width: 800, Height: 600},
				Seed:           123,
				IterationCount: 10000,
				Output:         "res.png",
				Threads:        4,
				Gamma:          2,
				EnableGamma:    true,
				SymmetryLevel:  3,
				AffineParams: []AffineTransform{
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.5, D: 0.0, E: 0.5, F: 0.0},
					{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.5},
				},
				Variations: []Variation{
					{Name: "linear", Weight: 0.7, Apply: nil},
					{Name: "swirl", Weight: 0.9, Apply: nil},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём копию, чтобы не изменять оригинальный input
			cfg := *tt.input
			cfg.SetDefaults()

			// Проверяем основные поля
			assert.Equal(t, tt.expected.Size, cfg.Size)
			assert.Equal(t, tt.expected.Seed, cfg.Seed)
			assert.Equal(t, tt.expected.IterationCount, cfg.IterationCount)
			assert.Equal(t, tt.expected.Output, cfg.Output)
			assert.Equal(t, tt.expected.Threads, cfg.Threads)
			assert.Equal(t, tt.expected.Gamma, cfg.Gamma)
			assert.Equal(t, tt.expected.SymmetryLevel, cfg.SymmetryLevel)
			assert.Equal(t, tt.expected.Variations, cfg.Variations)

			// Проверяем AffineParams (длина и коэффициенты)
			require.Len(t, cfg.AffineParams, len(tt.expected.AffineParams))
			for i := range cfg.AffineParams {
				assert.Equal(t, tt.expected.AffineParams[i].A, cfg.AffineParams[i].A)
				assert.Equal(t, tt.expected.AffineParams[i].B, cfg.AffineParams[i].B)
				assert.Equal(t, tt.expected.AffineParams[i].C, cfg.AffineParams[i].C)
				assert.Equal(t, tt.expected.AffineParams[i].D, cfg.AffineParams[i].D)
				assert.Equal(t, tt.expected.AffineParams[i].E, cfg.AffineParams[i].E)
				assert.Equal(t, tt.expected.AffineParams[i].F, cfg.AffineParams[i].F)
			}
		})
	}
}
