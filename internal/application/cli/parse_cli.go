package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// parseConfig парсит аргументы коммандной строки, не валидируя при этом
// По сути это просто проверка корректности формата
func parseConfig(cmd *cobra.Command, config *model.FlameConfig) error {
	if w, err := cmd.Flags().GetInt("width"); err != nil {
		return fmt.Errorf("invalid width value: %w", err)
	} else if w != 0 {
		config.Size.Width = w
	}

	if h, err := cmd.Flags().GetInt("height"); err != nil {
		return fmt.Errorf("invalid height value: %w", err)
	} else if h != 0 {
		config.Size.Height = h
	}

	if seed, err := cmd.Flags().GetInt64("seed"); err != nil {
		return fmt.Errorf("invalid seed value: %w", err)
	} else if seed != 0 {
		config.Seed = seed
	}

	if iterations, err := cmd.Flags().GetInt("iteration-count"); err != nil {
		return fmt.Errorf("invalid iterations count: %w", err)
	} else if iterations != 0 {
		config.IterationCount = iterations
	}

	if threads, err := cmd.Flags().GetInt("threads"); err != nil {
		return fmt.Errorf("invalid threads count: %w", err)
	} else if threads != 0 {
		config.Threads = threads
	}

	if out, err := cmd.Flags().GetString("output_path"); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	} else if out != "" {
		config.Output = out
	}

	if cmd.Flags().Changed("gamma-correction") {
		if cg, err := cmd.Flags().GetBool("gamma-correction"); err != nil {
			return fmt.Errorf("invalid gamma-correction value: %w", err)
		} else {
			config.EnableGamma = cg
		}
	}

	if gamma, err := cmd.Flags().GetFloat64("gamma"); err != nil {
		return fmt.Errorf("invalid gamma value: %w", err)
	} else if gamma != 0 {
		config.Gamma = gamma
	}

	if sym, err := cmd.Flags().GetInt("symmetry-level"); err != nil {
		return fmt.Errorf("invalid symmetry level: %w", err)
	} else if sym != 0 {
		config.SymmetryLevel = sym
	}

	if funcs, err := cmd.Flags().GetString("functions"); err != nil {
		return fmt.Errorf("invalid functions string: %w", err)
	} else if funcs != "" {
		variations, err := parseFunctions(funcs)
		if err != nil {
			return fmt.Errorf("failed to parse functions: %w", err)
		}
		if len(variations) != 0 {
			config.Variations = variations
		}
	}

	if affineStr, err := cmd.Flags().GetString("affine-params"); err != nil {
		return fmt.Errorf("invalid affine string: %w", err)
	} else if affineStr != "" {
		affineTransform, err := parseAffine(affineStr, config.Seed)
		if err != nil {
			return fmt.Errorf("failed to parse affine: %w", err)
		}
		if len(affineTransform) != 0 {
			config.AffineParams = affineTransform
		}
	}

	return nil
}

// parseFunctions парсит функции с весами, но не валидирует
// не проверяет функции на существование, а веса на положительность
func parseFunctions(s string) ([]model.Variation, error) {
	parts := strings.Split(s, ",")
	vars := make([]model.Variation, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		nameWeight := strings.Split(p, ":")
		if len(nameWeight) != 2 {
			return nil, fmt.Errorf("invalid function format '%s' (expected 'name:weight')", p)
		}

		name := strings.TrimSpace(nameWeight[0])
		weight, err := strconv.ParseFloat(strings.TrimSpace(nameWeight[1]), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid weight '%s' for function '%s': %w", nameWeight[1], name, err)
		}

		vars = append(vars, model.Variation{
			Name:   name,
			Weight: weight,
			Apply:  nil,
		})

	}

	// не возвращаем ошибку при нулевой длине, потому что у нас есть дефолтные параметры
	// if len(vars) == 0 {
	// 	return nil, fmt.Errorf("no valid functions found in: %s", s)
	// }

	return vars, nil
}

// parseAffine парсит афинные коэффициенты, но не валидирует их
// не проверяет коэффы на сжимаемость матрицы
func parseAffine(s string, seed int64) ([]model.AffineTransform, error) {
	transforms := strings.Split(s, "/")
	result := make([]model.AffineTransform, 0, len(transforms))

	for i, t := range transforms {
		t = strings.TrimSpace(t)
		if t == "" {
			continue
		}

		coeffs := strings.Split(t, ",")
		if len(coeffs) != 6 {
			return nil, fmt.Errorf("expected 6 coefficients, got %d in: %s", len(coeffs), t)
		}

		var floats [6]float64
		for j, c := range coeffs {
			val, err := strconv.ParseFloat(strings.TrimSpace(c), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid coefficient '%s' in: %s", c, t)
			}
			floats[j] = val
		}

		affine := model.AffineTransform{
			A: floats[0], B: floats[1], C: floats[2],
			D: floats[3], E: floats[4], F: floats[5],
		}

		// Используем уникальный seed для каждого трансформа: seed + индекс
		ft := model.NewAffineTransform(affine, seed+int64(i))
		result = append(result, ft)

	}

	// не возвращаем ошибку при нулевой длине, потому что у нас есть дефолтные параметры
	// if len(params) == 0 {
	// 	return nil, fmt.Errorf("no valid affine transforms found in: %s", s)
	// }

	return result, nil
}
