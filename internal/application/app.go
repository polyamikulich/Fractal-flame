package application

import (
	"fmt"
	"log/slog"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/application/generator"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/infrastructure/saver"
)

// Run — основная бизнес-логика приложения: генерация фрактала и сохранение в файл
func Run(config *model.FlameConfig) error {
	if config == nil {
		err := fmt.Errorf("config is nil")
		slog.Error("Config validation failed", "error", err)
		return err
	}

	slog.Info("Starting fractal flame generation")

	gen := generator.NewGenerator(config)
	buffer, err := gen.Generate()
	if err != nil {
		slog.Error("Failed to generate fractal flame", "error", err)
		return err
	}

	// Сохраняем изображение
	err = saver.SaveImageToFile(buffer, config.Output, config.Gamma, config.EnableGamma)
	if err != nil {
		slog.Error("Failed to save image", "error", err)
		return fmt.Errorf("failed to save image: %w", err)
	}

	slog.Info("Fractal flame generation completed successfully", "output", config.Output)

	return nil
}
