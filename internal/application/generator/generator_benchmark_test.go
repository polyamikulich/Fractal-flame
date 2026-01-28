package generator

import (
	"testing"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/variations"
)

func benchmarkConfig() *model.FlameConfig {
	return &model.FlameConfig{
		Size:           model.SizeConfig{Width: 800, Height: 600},
		Seed:           12345,
		IterationCount: 500000, // БОЛЬШЕ итераций для бенчмарков
		Output:         "",
		Threads:        1,
		EnableGamma:    true,
		Gamma:          2.2,
		SymmetryLevel:  1,
		AffineParams: []model.AffineTransform{
			{
				A: 0.5, B: 0.0, C: 0.0,
				D: 0.0, E: 0.5, F: 0.0,
				ColorR: 1.0, ColorG: 0.5, ColorB: 0.2,
			},
		},
		Variations: []model.Variation{
			{
				Name:   "linear",
				Weight: 1.0,
				Apply:  variations.Linear,
			},
		},
	}
}

func BenchmarkGenerate_1Thread(b *testing.B) {
	config := benchmarkConfig()
	config.Threads = 1

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator := NewGenerator(config)
		generator.Generate()
	}
}

func BenchmarkGenerate_2Thread(b *testing.B) {
	config := benchmarkConfig()
	config.Threads = 2

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator := NewGenerator(config)
		generator.Generate()
	}
}

func BenchmarkGenerate_4Thread(b *testing.B) {
	config := benchmarkConfig()
	config.Threads = 4

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator := NewGenerator(config)
		generator.Generate()
	}
}

func BenchmarkGenerate_8Thread(b *testing.B) {
	config := benchmarkConfig()
	config.Threads = 8

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generator := NewGenerator(config)
		generator.Generate()
	}
}
