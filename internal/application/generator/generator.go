package generator

import (
	"log/slog"
	"sync"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/application/worker"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

type Generator struct {
	config *model.FlameConfig
	buffer *model.PixelBuffer
}

// NewGenerator создает новый генератор
func NewGenerator(cfg *model.FlameConfig) *Generator {
	return &Generator{
		config: cfg,
		buffer: model.NewPixelBuffer(cfg.Size.Width, cfg.Size.Height),
	}
}

// Generate запускает генерацию фрактала
func (g *Generator) Generate() (*model.PixelBuffer, error) {
	var err error = nil

	slog.Info("Starting generation", "threads", g.config.Threads, "iterations", g.config.IterationCount)

	if g.config.Threads == 1 {
		// Однопоточный режим
		slog.Info("Running in single-threaded mode")
		worker := worker.NewWorker(g.config, g.buffer, 0, g.config.IterationCount)
		worker.Run()
	} else {
		// Многопоточный режим
		slog.Info("Running in multi-threaded mode", "thread_count", g.config.Threads)
		err = g.generateParallel()
		if err != nil {
			slog.Error("Parallel generation failed", "error", err)
			return nil, err
		}
	}

	slog.Info("Generation completed successfully")

	return g.buffer, err
}

func (g *Generator) generateParallel() error {
	threadBuffers := make([]*model.PixelBuffer, g.config.Threads)
	var wg sync.WaitGroup
	errCh := make(chan error, g.config.Threads)

	// Распределяем итерации между потоками
	baseIterations := g.config.IterationCount / g.config.Threads
	extraIterations := g.config.IterationCount % g.config.Threads

	for t := 0; t < g.config.Threads; t++ {
		iterations := baseIterations
		if t < extraIterations {
			iterations++
		}

		// Создаем буфер для потока
		threadBuffers[t] = model.NewPixelBuffer(g.config.Size.Width, g.config.Size.Height)

		wg.Add(1)

		go func(threadID int, iters int, buf *model.PixelBuffer) {
			defer wg.Done()

			slog.Debug("Starting worker thread", "thread_id", threadID, "iterations", iters)
			worker := worker.NewWorker(g.config, buf, threadID, iterations)
			worker.Run()
			slog.Debug("Worker thread finished", "thread_id", threadID)
		}(t, iterations, threadBuffers[t])
	}

	wg.Wait()
	close(errCh)

	// Объединяем результаты в финальный буфер
	for _, buf := range threadBuffers {
		err := g.buffer.MergeFrom(buf)
		if err != nil {
			slog.Error("Failed to merge buffer from thread", "error", err)
			return err
		}
	}

	return nil
}
