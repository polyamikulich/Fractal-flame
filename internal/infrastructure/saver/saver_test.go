package saver

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// mockVariationFunc - заглушка для VariationFunc
func mockVariationFunc(x, y float64) (float64, float64) {
	return x, y // тождественное преобразование
}

func TestSaveImageToFile(t *testing.T) {
	// Подготовим буфер с данными
	buffer := model.NewPixelBuffer(10, 10)
	// Добавим немного данных
	buffer.AddPoint(5, 5, 1.0, 0.0, 0.5) // Красно-синий пиксель

	// Создадим временный файл
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test_output.png")

	// Вызовем функцию сохранения
	err := SaveImageToFile(buffer, filename, 2.2, true) // с гаммой
	require.NoError(t, err)

	// Проверим, что файл существует
	_, err = os.Stat(filename)
	assert.NoError(t, err, "Ожидалось, что файл будет создан")

	// Проверим, что файл можно открыть как PNG
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()

	img, err := png.Decode(file)
	require.NoError(t, err)

	// Проверим размер изображения
	bounds := img.Bounds()
	assert.Equal(t, image.Rect(0, 0, 10, 10), bounds)

	// Проверим цвет в точке (5,5)
	expectedColor := color.RGBA{255, 0, 127, 255} // При гамме и плотности = 1
	actualColor := img.At(5, 5)
	assert.Equal(t, expectedColor, actualColor)
}

func TestSaveImageToFile_WithoutGamma(t *testing.T) {
	buffer := model.NewPixelBuffer(5, 5)
	buffer.AddPoint(2, 2, 0.8, 0.6, 0.4)

	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test_no_gamma.png")

	err := SaveImageToFile(buffer, filename, 2.2, false) // без гаммы
	require.NoError(t, err)

	_, err = os.Stat(filename)
	assert.NoError(t, err)
}

func TestSaveImageToFile_EmptyBuffer(t *testing.T) {
	buffer := model.NewPixelBuffer(5, 5) // пустой буфер

	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test_empty.png")

	err := SaveImageToFile(buffer, filename, 2.2, false)
	require.NoError(t, err)

	_, err = os.Stat(filename)
	assert.NoError(t, err)

	// Проверим, что изображение чёрное
	file, err := os.Open(filename)
	require.NoError(t, err)
	defer file.Close()

	img, err := png.Decode(file)
	require.NoError(t, err)

	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			c := img.At(x, y)
			assert.Equal(t, color.RGBA{0, 0, 0, 255}, c, "Пиксель (%d, %d) должен быть чёрным", x, y)
		}
	}
}
