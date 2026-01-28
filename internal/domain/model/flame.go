package model

// SizeConfig — структура для вложенного объекта "size" в JSON
type SizeConfig struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type FlameConfig struct {
	Size SizeConfig `json:"size"`
	//Width          int               `json:"width"`
	//Height         int               `json:"height"`
	Seed           int64             `json:"seed"`
	IterationCount int               `json:"iteration_count"`
	Output         string            `json:"output_path"`
	Threads        int               `json:"threads"`
	EnableGamma    bool              `json:"gamma_correction,omitempty"`
	Gamma          float64           `json:"gamma,omitempty"`
	SymmetryLevel  int               `json:"symmetry_level,omitempty"`
	Variations     []Variation       `json:"functions"`
	AffineParams   []AffineTransform `json:"affine_params"`
}

// SetDefaults устанавливает значения по умолчанию, если они не заданы.
// Вызывается после парсинга, до валидации.
func (c *FlameConfig) SetDefaults() {
	if c.Size.Width == 0 {
		c.Size.Width = 1920
	}
	if c.Size.Height == 0 {
		c.Size.Height = 1080
	}
	if c.IterationCount == 0 {
		c.IterationCount = 2500
	}
	if c.Seed == 0 {
		c.Seed = 5
	}
	if c.Output == "" {
		c.Output = "result.png"
	}
	if c.Threads == 0 {
		c.Threads = 1
	}
	if c.Gamma == 0 {
		c.Gamma = 2.2 // Дефолтное значение гаммы
	}
	if c.SymmetryLevel == 0 {
		c.SymmetryLevel = 1 // Без симметрии
	}
	if c.AffineParams == nil {
		defaultAffines := []AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
			{A: 0.5, B: 0.0, C: -0.5, D: 0.0, E: 0.5, F: 0.0},
			{A: 0.5, B: 0.0, C: 0.5, D: 0.0, E: 0.5, F: 0.0},
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.5},
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: -0.5},
		}
		c.AffineParams = make([]AffineTransform, len(defaultAffines))
		for i, aff := range defaultAffines {
			// Используем фиксированный seed для детерминированности
			c.AffineParams[i] = NewAffineTransform(aff, c.Seed+int64(i))
		}
	}
	if c.Variations == nil {
		c.Variations = []Variation{
			{Name: "linear", Weight: 1.0, Apply: nil},
		}
	}
}
