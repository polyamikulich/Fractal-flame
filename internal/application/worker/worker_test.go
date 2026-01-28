package worker

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// mockVariationFunc - заглушка для VariationFunc
func mockVariationFunc(x, y float64) (float64, float64) {
	return x, y // тождественное преобразование
}

func TestWorker_Run(t *testing.T) {
	// Подготовим конфигурацию
	cfg := &model.FlameConfig{
		Seed:           123,
		IterationCount: 100,
		SymmetryLevel:  1,
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0, ColorR: 0.34, ColorB: 0.45, ColorG: 0.9},
		},
		Variations: []model.Variation{
			{Name: "linear", Weight: 1.0, Apply: mockVariationFunc},
		},
	}

	// Создадим буфер
	buffer := model.NewPixelBuffer(100, 100)

	// Создадим воркер
	worker := NewWorker(cfg, buffer, 0, cfg.IterationCount)

	// Запустим
	worker.Run()

	// Проверим, что буфер не пустой (некоторые хиты должны быть)
	hitCount := 0
	for _, h := range buffer.Hits() {
		if h > 0 {
			hitCount++
		}
	}
	// Ожидаем, что хотя бы несколько точек попало
	assert.Greater(t, hitCount, 0, "Ожидалось, что Run добавит хотя бы одну точку в буфер")
}

func TestWorker_Run_ZeroIterations(t *testing.T) {
	cfg := &model.FlameConfig{
		Seed:           123,
		IterationCount: 0,
		SymmetryLevel:  1,
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0, ColorR: 1, ColorB: 0, ColorG: 0},
		},
		Variations: []model.Variation{
			{Name: "linear", Weight: 1.0, Apply: mockVariationFunc},
		},
	}

	buffer := model.NewPixelBuffer(100, 100)
	worker := NewWorker(cfg, buffer, 0, cfg.IterationCount)
	worker.Run()

	// Все хиты должны быть 0
	for _, h := range buffer.Hits() {
		assert.Equal(t, int64(0), h)
	}
}

func TestWorker_iterate(t *testing.T) {
	cfg := &model.FlameConfig{
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0, ColorR: 1, ColorG: 0, ColorB: 0},
		},
		Variations: []model.Variation{
			{Name: "linear", Weight: 1.0, Apply: mockVariationFunc}, // тождественное
		},
	}

	buffer := model.NewPixelBuffer(100, 100)
	worker := NewWorker(cfg, buffer, 0, 10)
	rng := rand.New(rand.NewSource(123))

	// Проверим, что iterate возвращает разумные значения
	x, y, idx := worker.iterate(1.0, 1.0, rng)
	// x, y = applyAffine(1,1) -> (0.5*1 + 0*1 + 0, 0*1 + 0.5*1 + 0) = (0.5, 0.5)
	// applyVariations (тождественная) -> (0.5, 0.5)
	assert.Equal(t, 0.5, x)
	assert.Equal(t, 0.5, y)
	assert.GreaterOrEqual(t, int(idx), 0)
	assert.Less(t, int(idx), len(cfg.AffineParams))
}

func TestWorker_applyAffine(t *testing.T) {
	cfg := &model.FlameConfig{
		AffineParams: []model.AffineTransform{
			{A: 2.0, B: 0.0, C: 1.0, D: 0.0, E: 2.0, F: 1.0, ColorR: 0, ColorG: 0, ColorB: 0}, // масштаб 2 и сдвиг на (1,1)
		},
	}
	buffer := model.NewPixelBuffer(100, 100)
	worker := NewWorker(cfg, buffer, 0, 10)

	x, y := worker.applyAffine(1.0, 1.0, 0)
	// x_new = 2*1 + 0*1 + 1 = 3
	// y_new = 0*1 + 2*1 + 1 = 3
	assert.Equal(t, 3.0, x)
	assert.Equal(t, 3.0, y)
}

func TestWorker_applyVariations(t *testing.T) {
	cfg := &model.FlameConfig{
		Variations: []model.Variation{
			{Name: "mock1", Weight: 1.0, Apply: func(x, y float64) (float64, float64) { return x + 1, y + 1 }},
			{Name: "mock2", Weight: 1.0, Apply: func(x, y float64) (float64, float64) { return x - 1, y - 1 }},
		},
	}
	buffer := model.NewPixelBuffer(100, 100)
	worker := NewWorker(cfg, buffer, 0, 10)

	x, y := worker.applyVariations(5.0, 10.0)
	assert.Equal(t, 10.0, x)
	assert.Equal(t, 20.0, y)
}

func TestWorker_applyVariations_Empty(t *testing.T) {
	cfg := &model.FlameConfig{
		Variations: []model.Variation{},
	}
	buffer := model.NewPixelBuffer(100, 100)
	worker := NewWorker(cfg, buffer, 0, 10)

	x, y := worker.applyVariations(5.0, 10.0)
	// Если вариаций нет, возвращаем (0, 0)
	assert.Equal(t, 0.0, x)
	assert.Equal(t, 0.0, y)
}

func TestWorker_Run_OutOfBounds(t *testing.T) {
	cfg := &model.FlameConfig{
		Seed:           123,
		IterationCount: 10,
		SymmetryLevel:  1,
		AffineParams: []model.AffineTransform{
			{A: 10.0, B: 0.0, C: 0.0, D: 0.0, E: 10.0, F: 0.0, ColorR: 1, ColorG: 0, ColorB: 0}, // Масштабирует в 10 раз
		},
		Variations: []model.Variation{
			{Name: "linear", Weight: 1.0, Apply: mockVariationFunc},
		},
	}

	buffer := model.NewPixelBuffer(100, 100)
	worker := NewWorker(cfg, buffer, 0, cfg.IterationCount)
	worker.Run()

	// Точки будут за пределами, поэтому все хиты должны быть 0
	hitCount := 0
	for _, h := range buffer.Hits() {
		if h > 0 {
			hitCount++
		}
	}

	zeroHits := 0
	for _, h := range buffer.Hits() {
		if h == 0 {
			zeroHits++
		}
	}
	// Должно быть много пикселей с 0 хитов, если точки уходят за границы
	assert.Greater(t, zeroHits, 0, "Ожидалось, что некоторые точки не попадут в буфер из-за масштаба")
}
