package saver

import (
	"image/png"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// SaveImageToFile сохраняет PixelBuffer в PNG-файл
func SaveImageToFile(buffer *model.PixelBuffer, filename string, gamma float64, enableGamma bool) error {
	img := buffer.RenderToImage(gamma, enableGamma)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
