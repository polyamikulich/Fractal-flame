package model

import (
	"image/color"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPixelBuffer(t *testing.T) {
	width, height := 100, 50
	pb := NewPixelBuffer(width, height)

	assert.Equal(t, width, pb.Width())
	assert.Equal(t, height, pb.Height())

	// Проверка диапазонов
	aspect := float64(width) / float64(height)
	assert.Equal(t, -aspect, pb.XMin)
	assert.Equal(t, aspect, pb.XMax)
	assert.Equal(t, -1.0, pb.YMin)
	assert.Equal(t, 1.0, pb.YMax)

	// Проверка размеров слайсов
	expectedSize := width * height
	assert.Len(t, pb.r, expectedSize)
	assert.Len(t, pb.g, expectedSize)
	assert.Len(t, pb.b, expectedSize)
	assert.Len(t, pb.hits, expectedSize)
}

func TestAddPoint(t *testing.T) {
	pb := NewPixelBuffer(10, 10)

	// Точка внутри диапазона
	pb.AddPoint(5, 5, 1.0, 0.5, 0.0)
	idx := 5*pb.width + 5
	assert.Equal(t, int64(1), pb.hits[idx])
	assert.Equal(t, 1.0, pb.r[idx])
	assert.Equal(t, 0.5, pb.g[idx])
	assert.Equal(t, 0.0, pb.b[idx])

	// Второй раз в ту же точку - усреднение
	pb.AddPoint(5, 5, 0.0, 0.0, 1.0)
	assert.Equal(t, int64(2), pb.hits[idx])
	assert.Equal(t, 0.5, pb.r[idx])  // (1.0 + 0.0) / 2
	assert.Equal(t, 0.25, pb.g[idx]) // (0.5 + 0.0) / 2
	assert.Equal(t, 0.5, pb.b[idx])  // (0.0 + 1.0) / 2

	// Точка за пределами диапазона - не должна добавиться
	pb.AddPoint(15, 5, 0.0, 0.0, 0.0)
	assert.Equal(t, int64(2), pb.hits[idx]) // Не изменилось
}

func TestAddPoint_OutOfBounds(t *testing.T) {
	pb := NewPixelBuffer(10, 10)

	// Проверка отрицательных индексов
	pb.AddPoint(-1, 5, 1.0, 0.0, 0.0)
	pb.AddPoint(5, -1, 1.0, 0.0, 0.0)

	// Проверка индексов за пределами
	pb.AddPoint(10, 5, 1.0, 0.0, 0.0)
	pb.AddPoint(5, 10, 1.0, 0.0, 0.0)

	// Никакие хиты не должны измениться
	for _, h := range pb.hits {
		assert.Equal(t, int64(0), h)
	}
}

func TestRenderToImage_Empty(t *testing.T) {
	pb := NewPixelBuffer(10, 10)
	img := pb.RenderToImage(2.2, false)

	// Изображение должно быть чёрным
	black := color.RGBA{0, 0, 0, 255}
	for x := 0; x < pb.width; x++ {
		for y := 0; y < pb.height; y++ {
			assert.Equal(t, black, img.At(x, y))
		}
	}
}

func TestRenderToImage_SinglePoint(t *testing.T) {
	pb := NewPixelBuffer(10, 10)
	pb.AddPoint(5, 5, 1.0, 0.0, 0.5) // r=1, g=0, b=0.5

	img := pb.RenderToImage(2.2, false)

	// Проверяем, что точка (5,5) окрашена, остальные - чёрные
	for x := 0; x < pb.width; x++ {
		for y := 0; y < pb.height; y++ {
			if x == 5 && y == 5 {
				assert.Equal(t, color.RGBA{255, 0, 127, 255}, img.At(x, y))
			} else {
				assert.Equal(t, color.RGBA{0, 0, 0, 255}, img.At(x, y))
			}
		}
	}
}

func TestRenderToImage_WithGamma(t *testing.T) {
	pb := NewPixelBuffer(10, 10)
	// Добавим точку, чтобы maxHits > hits в этой точке -> density < 1
	pb.AddPoint(5, 5, 0.8, 0.4, 0.2) // hits = 1
	pb.AddPoint(6, 6, 1.0, 1.0, 1.0) // hits = 1, maxHits = 1
	pb.AddPoint(6, 6, 1.0, 1.0, 1.0) // hits = 2, maxHits = 2

	imgNoGamma := pb.RenderToImage(2.2, false)
	imgWithGamma := pb.RenderToImage(2.0, true) // gamma = 2.0

	cNoGamma := imgNoGamma.At(5, 5)
	cWithGamma := imgWithGamma.At(5, 5)

	// Проверим, что цвета отличаются
	assert.NotEqual(t, cNoGamma, cWithGamma, "Цвета в (5,5) должны отличаться с и без гаммы")
}

func TestRenderToImage_MultiplePoints(t *testing.T) {
	pb := NewPixelBuffer(5, 5)
	pb.AddPoint(0, 0, 1.0, 0.0, 0.0) // Red
	pb.AddPoint(0, 0, 0.0, 1.0, 0.0) // Green -> усреднение: 0.5, 0.5, 0.0
	pb.AddPoint(1, 1, 0.0, 0.0, 1.0) // Blue

	img := pb.RenderToImage(2.2, false)

	// Проверка (0,0)
	c00 := img.At(0, 0)
	require.Equal(t, color.RGBA{127, 127, 0, 255}, c00)

	// Проверка (1,1)
	c11 := img.At(1, 1)
	density := math.Log(2) / math.Log(3)
	expectedR := uint8(0.0)
	expectedG := uint8(0.0)
	expectedB := uint8(density * 1.0 * 255)
	require.Equal(t, color.RGBA{expectedR, expectedG, expectedB, 255}, c11)
	// Остальные чёрные
	c22 := img.At(2, 2)
	require.Equal(t, color.RGBA{0, 0, 0, 255}, c22)
}

func TestWidth(t *testing.T) {
	pb := NewPixelBuffer(123, 456)
	assert.Equal(t, 123, pb.Width())
}

func TestHeight(t *testing.T) {
	pb := NewPixelBuffer(123, 456)
	assert.Equal(t, 456, pb.Height())
}

func TestHits(t *testing.T) {
	pb := NewPixelBuffer(123, 456)
	assert.Equal(t, 56088, len(pb.Hits()))
}
