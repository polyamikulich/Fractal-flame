package cli

import (
	"fmt"
	"log/slog"
	"math"
	"os"
	"path/filepath"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/variations"
)

func validateConfig(config *model.FlameConfig) error {
	if err := ValidateInt(config); err != nil {
		return err
	}

	if err := ValidateAffineParams(config.AffineParams); err != nil {
		return err
	}

	if err := ValidateFunctions(config.Variations); err != nil {
		return err
	}

	if err := ValidateOutput(config.Output); err != nil {
		return err
	}

	return nil
}

func ValidateInt(c *model.FlameConfig) error {
	if c.Size.Width <= 0 {
		return fmt.Errorf("width must be positive, got: %d", c.Size.Width)
	}
	if c.Size.Height <= 0 {
		return fmt.Errorf("height must be positive, got: %d", c.Size.Height)
	}
	if c.IterationCount <= 0 {
		return fmt.Errorf("iteration count must be positive, got: %d", c.IterationCount)
	}
	if c.Threads <= 0 {
		return fmt.Errorf("threads must be positive, got: %d", c.Threads)
	}
	if c.Gamma <= 0 {
		return fmt.Errorf("gamma must be positive, got: %f", c.Gamma)
	}
	if c.SymmetryLevel < 1 {
		return fmt.Errorf("symmetry level must be >= 1, got: %d", c.SymmetryLevel)
	}

	return nil
}

func ValidateAffineParams(affineTransforms []model.AffineTransform) error {
	if len(affineTransforms) == 0 {
		return fmt.Errorf("affine params must be non-empty")
	}
	for i, at := range affineTransforms {
		//aff := at.Affine
		if math.IsNaN(at.A) || math.IsInf(at.A, 0) ||
			math.IsNaN(at.B) || math.IsInf(at.B, 0) ||
			math.IsNaN(at.C) || math.IsInf(at.C, 0) ||
			math.IsNaN(at.D) || math.IsInf(at.D, 0) ||
			math.IsNaN(at.E) || math.IsInf(at.E, 0) ||
			math.IsNaN(at.F) || math.IsInf(at.F, 0) {
			return fmt.Errorf("affine param #%d contains invalid coefficients (NaN or Inf): a=%f, b=%f, c=%f, d=%f, e=%f, f=%f", i, at.A, at.B, at.C, at.D, at.E, at.F)
		}

		if at.A*at.A+at.D*at.D >= 1 {
			slog.Warn("Affine transform may not be contractive",
				"transform_index", i,
				"condition", "a^2 + d^2 < 1",
				"value", at.A*at.A+at.D*at.D,
			)
		}

		if at.B*at.B+at.E*at.E > 1 {
			slog.Warn("Affine transform may not be contractive",
				"transform_index", i,
				"condition", "b^2 + e^2 < 1",
				"value", at.B*at.B+at.E*at.E,
			)
		}

		det := at.A*at.E - at.B*at.D
		rightSide := 1 + det*det
		if at.A*at.A+at.D*at.D+at.B*at.B+at.E*at.E > 1+rightSide {
			slog.Warn("Affine transform may violate convergence condition",
				"transform_index", i,
				"condition", "a^2 + d^2 + b^2 + e^2 < 1 + (a*e - b*d)^2",
				"value_left_side", at.A*at.A+at.D*at.D+at.B*at.B+at.E*at.E,
				"value_right_side", rightSide,
			)
		}
	}

	return nil
}

func ValidateFunctions(functions []model.Variation) error {
	if len(functions) == 0 {
		return fmt.Errorf("variations must be non-empty")
	}
	for i := range functions {
		v := &functions[i]
		if math.IsNaN(v.Weight) || math.IsInf(v.Weight, 0) || v.Weight <= 0 {
			return fmt.Errorf("variation #%d (%s) has invalid weight (NaN or Inf or non-positive): %f", i, v.Name, v.Weight)
		}

		if variation, err := variations.GetVariation(v.Name); err != nil {
			return fmt.Errorf("variation #%d (%s) is not a valid function", i, v.Name)
		} else {
			v.Apply = variation
		}
	}

	return nil
}

func ValidateOutput(output string) error {
	if output == "" {
		return fmt.Errorf("output file is empty")
	}

	dir := filepath.Dir(output)

	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("output directory does not exist: %s", dir)
		}
		return fmt.Errorf("cannot access output directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("output path directory is a file, not a directory: %s", dir)
	}

	testFile := filepath.Join(dir, ".fractal_flame_write_test")
	file, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("cannot create file in output directory: %w", err)
	}
	file.Close()
	os.Remove(testFile)

	return nil

}
