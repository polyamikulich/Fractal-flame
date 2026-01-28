package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// parseJSON парсит конфигурационный файл и заполняет поля структуры FlameConfig
func parseJSON(path string, cfg *model.FlameConfig) error {
	//Проверка существования файла
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", path)
	}

	// Проверка расширения файла
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		return fmt.Errorf("config file must have .json extension: %s", path)
	}

	// Попытка прочитать файл
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}
