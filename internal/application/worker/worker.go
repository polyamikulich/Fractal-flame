package worker

import (
	"log/slog"
	"math"
	"math/rand"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

type Worker struct {
	config     *model.FlameConfig
	buffer     *model.PixelBuffer
	threadID   int
	iterations int
}

func NewWorker(config *model.FlameConfig, buffer *model.PixelBuffer, threadID int, iterations int) *Worker {

	return &Worker{
		config:     config,
		buffer:     buffer,
		threadID:   threadID,
		iterations: iterations,
	}
}

func (w *Worker) Run() {
	source := rand.NewSource(w.config.Seed + int64(w.threadID)*12345)
	rng := rand.New(source)

	// Диапазон для начальной точки
	xMin := w.buffer.XMin
	xMax := w.buffer.XMax
	yMin := w.buffer.YMin
	yMax := w.buffer.YMax

	// Начальная точка
	x := rng.Float64()*(xMax-xMin) + xMin
	y := rng.Float64()*(yMax-yMin) + yMin

	// Логирование прогресса
	progressStep := w.iterations / 5 // 20% шаги
	if progressStep == 0 {
		progressStep = 1
	}

	for i := -20; i < w.iterations; i++ {
		xnew, ynew, transformIdx := w.iterate(x, y, rng)

		angleStep := 2 * math.Pi / float64(w.config.SymmetryLevel)
		for k := 0; k < w.config.SymmetryLevel; k++ {
			angle := float64(k) * angleStep
			cos_a := math.Cos(angle)
			sin_a := math.Sin(angle)

			sx_math := xnew*cos_a - ynew*sin_a
			sy_math := xnew*sin_a + ynew*cos_a

			if i > 0 && sx_math >= xMin && sx_math <= xMax && sy_math >= yMin && sy_math <= yMax {
				x1 := w.buffer.Width() - int(((xMax-sx_math)/(xMax-xMin))*float64(w.buffer.Width()))
				y1 := w.buffer.Height() - int(((yMax-sy_math)/(yMax-yMin))*float64(w.buffer.Height()))

				// Проверяем, попадает ли пиксель в границы
				if x1 >= 0 && x1 < w.buffer.Width() && y1 >= 0 && y1 < w.buffer.Height() {
					transform := w.config.AffineParams[transformIdx]
					r, g, b := transform.ColorR, transform.ColorG, transform.ColorB

					w.buffer.AddPoint(x1, y1, r, g, b)
				}
			}
		}

		x = xnew
		y = ynew

		// Логируем прогресс каждые 20%
		if i >= 0 && i%progressStep == 0 {
			percent := float64(i) / float64(w.iterations) * 100
			slog.Info("Worker progress", "thread_id", w.threadID, "iteration", i, "progress_percent", percent)
		}
	}
}

func (w *Worker) iterate(x, y float64, r *rand.Rand) (float64, float64, int64) {
	affineInd := r.Intn(len(w.config.AffineParams))

	xa, ya := w.applyAffine(x, y, affineInd)

	xv, yv := w.applyVariations(xa, ya)

	return xv, yv, int64(affineInd)
}

func (w *Worker) applyAffine(x, y float64, affineIdx int) (float64, float64) {
	at := w.config.AffineParams[affineIdx]

	newX := at.A*x + at.B*y + at.C
	newY := at.D*x + at.E*y + at.F

	return newX, newY
}

func (w *Worker) applyVariations(x, y float64) (float64, float64) {
	variations := w.config.Variations

	resultX := 0.0
	resultY := 0.0
	for _, variation := range variations {
		vx, vy := variation.Apply(x, y)

		resultX += vx * variation.Weight
		resultY += vy * variation.Weight
	}

	return resultX, resultY
}
