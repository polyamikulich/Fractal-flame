package model

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"
)

// PixelBuffer - структура для хранения пикселей
// Вместо хранения массива пикселей, храним массивы для r, g, b и hits
type PixelBuffer struct {
	width, height int
	XMin, XMax    float64
	YMin, YMax    float64
	r, g, b       []float64 //длина = width * height
	hits          []int64
	mu            sync.RWMutex // Для потокобезопасности
}

// NewPixelBuffer создаёт новый буфер пикселей заданного размера
func NewPixelBuffer(width, height int) *PixelBuffer {
	// Рассчитываем диапазон координат на основе пропорций
	// Здесь можно изменить домножающий коэфф. (например, на 10) -> мы как бы отдалимся от изображения тогда
	aspect := float64(width) / float64(height)
	xMin := -1.0 * aspect
	xMax := 1.0 * aspect
	yMin := -1.0
	yMax := 1.0

	size := width * height
	return &PixelBuffer{

		width:  width,
		height: height,
		XMin:   xMin,
		XMax:   xMax,
		YMin:   yMin,
		YMax:   yMax,
		r:      make([]float64, size),
		g:      make([]float64, size),
		b:      make([]float64, size),
		hits:   make([]int64, size),
	}
}

// AddPoint добавляет точку в массив пикселей
// Знает только координаты и цвет
func (pb *PixelBuffer) AddPoint(x, y int, r, g, b float64) {
	if x < 0 || x >= pb.width || y < 0 || y >= pb.height {
		return
	}

	idx := y*pb.width + x

	// // Блокируем для записи
	// pb.mu.Lock()
	// defer pb.mu.Unlock()

	// Обновление цвета
	if pb.hits[idx] == 0 {
		pb.r[idx] = r
		pb.g[idx] = g
		pb.b[idx] = b
	} else {
		pb.r[idx] = (pb.r[idx] + r) / 2.0
		pb.g[idx] = (pb.g[idx] + g) / 2.0
		pb.b[idx] = (pb.b[idx] + b) / 2.0
	}

	pb.hits[idx]++
}

// Width возвращает ширину буфера
func (pb *PixelBuffer) Width() int {
	return pb.width
}

// Height возвращает высоту буфера
func (pb *PixelBuffer) Height() int {
	return pb.height
}

func (pb *PixelBuffer) Hits() []int64 {
	return pb.hits
}

// RenderToImage рендерит буфер в изображение
// Знает только про буфер и гамму, больше ему ничего и не надо
func (pb *PixelBuffer) RenderToImage(gamma float64, enableGamma bool) image.Image {
	if pb.width <= 0 || pb.height <= 0 {
		return nil
	}

	img := image.NewRGBA(image.Rect(0, 0, pb.width, pb.height))

	var maxHits int64
	for _, hits := range pb.hits {
		if hits > maxHits {
			maxHits = hits
		}
	}

	if maxHits == 0 {
		// Заполняем всё изображение чёрным непрозрачным цветом
		for y := 0; y < pb.height; y++ {
			for x := 0; x < pb.width; x++ {
				img.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
		return img // Чёрное изображение
	}

	logMaxHits := math.Log(float64(1 + maxHits))

	for y := 0; y < pb.height; y++ {
		for x := 0; x < pb.width; x++ {
			idx := y*pb.width + x
			hits := pb.hits[idx]

			if hits == 0 {
				img.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
				continue
			}

			density := math.Log(float64(1+hits)) / logMaxHits

			if enableGamma {
				density = math.Pow(density, 1.0/gamma)
			}

			r := uint8(math.Min(255, 255.0*density*pb.r[idx]))
			g := uint8(math.Min(255, 255.0*density*pb.g[idx]))
			b := uint8(math.Min(255, 255.0*density*pb.b[idx]))

			img.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

// MergeFrom объединяет данные из другого буфера в текущий
func (pb *PixelBuffer) MergeFrom(other *PixelBuffer) error {
	if pb.width != other.width || pb.height != other.height {
		return fmt.Errorf("buffer size mismatch")
	}

	size := pb.width * pb.height
	pb.mu.Lock()
	defer pb.mu.Unlock()

	for i := 0; i < size; i++ {
		otherHits := other.hits[i]
		if otherHits == 0 {
			continue
		}

		pb.hits[i] += otherHits

		totalHits := pb.hits[i]
		pb.r[i] = (pb.r[i]*float64(totalHits-otherHits) + other.r[i]*float64(otherHits)) / float64(totalHits)
		pb.g[i] = (pb.g[i]*float64(totalHits-otherHits) + other.g[i]*float64(otherHits)) / float64(totalHits)
		pb.b[i] = (pb.b[i]*float64(totalHits-otherHits) + other.b[i]*float64(otherHits)) / float64(totalHits)
	}

	return nil
}
