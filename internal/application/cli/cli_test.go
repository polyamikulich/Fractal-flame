package cli

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// ------ Тесты для cli.go ----------

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()

	assert.Contains(t, cmd.Use, "fractal-flame")
	assert.Equal(t, "Generate a flame fractal image", cmd.Short)
	assert.Equal(t, "Generate a flame fractal image", cmd.Long)
	assert.True(t, cmd.SilenceUsage)

	// Проверяем, что флаги установлены
	require.NotNil(t, cmd.Flags().Lookup("width"))
	require.NotNil(t, cmd.Flags().Lookup("height"))
	require.NotNil(t, cmd.Flags().Lookup("output_path"))
	require.NotNil(t, cmd.Flags().Lookup("seed"))
	require.NotNil(t, cmd.Flags().Lookup("iteration-count"))
	require.NotNil(t, cmd.Flags().Lookup("output_path"))
	require.NotNil(t, cmd.Flags().Lookup("threads"))
	require.NotNil(t, cmd.Flags().Lookup("gamma"))
	require.NotNil(t, cmd.Flags().Lookup("gamma-correction"))
	require.NotNil(t, cmd.Flags().Lookup("config"))
	require.NotNil(t, cmd.Flags().Lookup("symmetry-level"))
	require.NotNil(t, cmd.Flags().Lookup("affine-params"))
	require.NotNil(t, cmd.Flags().Lookup("functions"))
}

func TestRootCommand_HelpFlag(t *testing.T) {
	cmd := NewRootCommand()

	// Проверяем, что команда создана
	require.NotNil(t, cmd)

	// Перенаправление вывода в pipe
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	//defer w.Close()
	os.Stdout = w

	// Выполнение: запускаем команду с флагом --help
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()

	// Восстанавливаем stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	// Проверки
	assert.NoError(t, err) // Execute должен завершиться успешно при --help

	// Проверим, что в выводе есть все необходимые флаги
	assert.Contains(t, string(out), "Generate a flame fractal image")
	assert.Contains(t, string(out), "-h, --help")
	assert.Contains(t, string(out), "--config")
	assert.Contains(t, string(out), "-W, --width")
	assert.Contains(t, string(out), "Width of the output image")
	assert.Contains(t, string(out), "-H, --height")
	assert.Contains(t, string(out), "--threads")
	assert.Contains(t, string(out), "Number of threads to use")
	assert.Contains(t, string(out), "--gamma-correction")
	assert.Contains(t, string(out), "--gamma")
	assert.Contains(t, string(out), "Affine parameters: <a1>,<b1>,<c1>,<d1>,<e1>,<f1>/...")
	assert.Contains(t, string(out), "Symmetry level (N-way symmetry)")
	assert.Contains(t, string(out), "-o, --output_path")
	assert.Contains(t, string(out), "Path to JSON configuration file")
}

func TestGetConfig(t *testing.T) {
	// Случай, когда cliConfig == nil
	_, err := GetConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration not loaded")

	// Случай, когда cliConfig != nil
	expectedConfig := &model.FlameConfig{
		Size:           model.SizeConfig{Width: 800, Height: 600},
		Seed:           456,
		IterationCount: 500,
		Output:         "temp.png",
		Threads:        4,
		EnableGamma:    true,
		Gamma:          1.8,
		SymmetryLevel:  2,
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
		},
		Variations: []model.Variation{
			{Name: "sinusoidal", Weight: 0.5, Apply: nil},
		},
	}
	cliConfig = expectedConfig

	cfg, err := GetConfig()
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, cfg)

	// Сбросим для других тестов
	cliConfig = nil
}

// ------ Тесты для parse_cli.go ----------

func TestParseConfig_Valid(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().Int("width", 0, "")
	cmd.Flags().Int("height", 0, "")
	cmd.Flags().Int64("seed", 0, "")
	cmd.Flags().Int("iteration-count", 0, "")
	cmd.Flags().String("output_path", "", "")
	cmd.Flags().Int("threads", 0, "")
	cmd.Flags().Bool("gamma-correction", false, "")
	cmd.Flags().Float64("gamma", 0, "")
	cmd.Flags().Int("symmetry-level", 0, "")
	cmd.Flags().String("affine-params", "", "")
	cmd.Flags().String("functions", "", "")

	cmd.Flags().Set("width", "800")
	cmd.Flags().Set("height", "600")
	cmd.Flags().Set("seed", "123")
	cmd.Flags().Set("iteration-count", "1000")
	cmd.Flags().Set("output_path", "test.png")
	cmd.Flags().Set("threads", "2")
	cmd.Flags().Set("gamma-correction", "true")
	cmd.Flags().Set("gamma", "2.0")
	cmd.Flags().Set("symmetry-level", "3")
	cmd.Flags().Set("affine-params", "1.0,0.0,0.0,0.0,1.0,0.0")
	cmd.Flags().Set("functions", "linear:1.0")

	cfg := &model.FlameConfig{}
	err := parseConfig(cmd, cfg)
	require.NoError(t, err)

	assert.Equal(t, 800, cfg.Size.Width)
	assert.Equal(t, 600, cfg.Size.Height)
	assert.Equal(t, int64(123), cfg.Seed)
	assert.Equal(t, 1000, cfg.IterationCount)
	assert.Equal(t, "test.png", cfg.Output)
	assert.Equal(t, 2, cfg.Threads)
	assert.True(t, cfg.EnableGamma)
	assert.Equal(t, 2.0, cfg.Gamma)
	assert.Equal(t, 3, cfg.SymmetryLevel)
	assert.Len(t, cfg.AffineParams, 1)
	assert.Len(t, cfg.Variations, 1)
	assert.Equal(t, "linear", cfg.Variations[0].Name)
	assert.Equal(t, 1.0, cfg.Variations[0].Weight)
}

func TestParseConfig_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		args        []string // аргументы командной строки
		expectError bool
	}{
		{
			name:        "invalid width",
			args:        []string{"--width", "abc"},
			expectError: true,
		},
		{
			name:        "invalid height",
			args:        []string{"--height", "B"},
			expectError: true,
		},
		{
			name:        "invalid seed",
			args:        []string{"--seed", "not_a_number"},
			expectError: true,
		},
		{
			name:        "invalid iterations",
			args:        []string{"--iteration-count", "not_a_number"},
			expectError: true,
		},
		{
			name:        "invalid threads",
			args:        []string{"--threads", "d"},
			expectError: true,
		},
		{
			name:        "invalid gamma",
			args:        []string{"--gamma", "negative_gamma"},
			expectError: true,
		},
		{
			name:        "invalid symmetry level",
			args:        []string{"--symmetry-level", "not_count"},
			expectError: true,
		},
		{
			name:        "invalid affine params format",
			args:        []string{"--affine-params", "1.0,2.0"}, // не хватает коэффициентов
			expectError: true,
		},
		{
			name:        "invalid functions weight",
			args:        []string{"--functions", "linear:not_a_number"},
			expectError: true,
		},
		{
			name:        "valid config",
			args:        []string{"--width", "800", "--functions", "linear:1.0"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			cmd.Flags().IntP("width", "W", 0, "Width of the output image")
			cmd.Flags().IntP("height", "H", 0, "Height of the output image")
			cmd.Flags().Int64("seed", 0, "Seed for random number generator")
			cmd.Flags().IntP("iteration-count", "i", 0, "Number of iterations")
			cmd.Flags().StringP("output_path", "o", "", "Output image file path")
			cmd.Flags().IntP("threads", "t", 0, "Number of threads to use")
			cmd.Flags().BoolP("gamma-correction", "g", false, "Enable gamma correction")
			cmd.Flags().Float64("gamma", 0, "Gamma value for correction (if enabled)")
			cmd.Flags().IntP("symmetry-level", "s", 0, "Symmetry level (N-way symmetry)")
			cmd.Flags().String("config", "", "Path to JSON configuration file")
			cmd.Flags().StringP("affine-params", "a", "", "Affine parameters")
			cmd.Flags().StringP("functions", "f", "", "Variations and weights")

			for i := 0; i < len(tt.args); i += 2 {
				if i+1 >= len(tt.args) {
					t.Fatalf("Invalid args format: %v", tt.args)
				}
				flagName := strings.TrimPrefix(tt.args[i], "--")
				flagName = strings.TrimPrefix(flagName, "-")
				flagValue := tt.args[i+1]
				err := cmd.Flags().Set(flagName, flagValue)
				if err != nil {
					if tt.expectError {
						return
					} else {
						t.Fatalf("Failed to set flag %s to %s: %v", flagName, flagValue, err)
					}
				}
			}

			cfg := &model.FlameConfig{}
			err := parseConfig(cmd, cfg)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseFunctions(t *testing.T) {
	s := "linear:1.0,sinusoidal:0.5"
	vars, err := parseFunctions(s)
	require.NoError(t, err)
	require.Len(t, vars, 2)

	assert.Equal(t, "linear", vars[0].Name)
	assert.Equal(t, 1.0, vars[0].Weight)
	assert.Nil(t, vars[0].Apply)

	assert.Equal(t, "sinusoidal", vars[1].Name)
	assert.Equal(t, 0.5, vars[1].Weight)
	assert.Nil(t, vars[1].Apply)
}

func TestParseFunctions_InvalidFormat(t *testing.T) {
	s := "linear:1.0,invalid"
	_, err := parseFunctions(s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid function format")
}

func TestParseFunctions_InvalidWeight(t *testing.T) {
	s := "linear:abc"
	_, err := parseFunctions(s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid weight")
}

func TestParseAffine_Valid(t *testing.T) {
	s := "1.0,0.0,0.0,0.0,1.0,0.0/0.1,0.3,0,0.3,0.1,0.4"
	affines, err := parseAffine(s, 123)
	require.NoError(t, err)
	require.Len(t, affines, 2)

	assert.Equal(t, 1.0, affines[0].A)
	assert.Equal(t, 0.0, affines[0].B)
	assert.Equal(t, 0.0, affines[0].C)
	assert.Equal(t, 0.0, affines[0].D)
	assert.Equal(t, 1.0, affines[0].E)
	assert.Equal(t, 0.0, affines[0].F)

	assert.Equal(t, 0.1, affines[1].A)
	assert.Equal(t, 0.3, affines[1].B)
	assert.Equal(t, 0.0, affines[1].C)
	assert.Equal(t, 0.3, affines[1].D)
	assert.Equal(t, 0.1, affines[1].E)
	assert.Equal(t, 0.4, affines[1].F)

	// Color устанавливается в NewAffineTransform
	assert.NotEmpty(t, affines[0].ColorR)
	assert.NotEmpty(t, affines[0].ColorG)
	assert.NotEmpty(t, affines[0].ColorB)

	assert.NotEmpty(t, affines[1].ColorR)
	assert.NotEmpty(t, affines[1].ColorG)
	assert.NotEmpty(t, affines[1].ColorB)
}

func TestParseAffine_InvalidNumCoeffs(t *testing.T) {
	s := "1.0,0,0,0,1/0,1,1,1,1,1"
	_, err := parseAffine(s, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected 6 coefficients")
}

func TestParseAffine_InvalidValue(t *testing.T) {
	s := "1.0,0.0,0.0,0.0,abc,0.0"
	_, err := parseAffine(s, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid coefficient")
}

// ------ Тесты для parse_json.go ----------

// Тест для parseJSON
func TestParseJSON_AfterEmptyConfig(t *testing.T) {
	// Создаём временный JSON-файл
	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testConfig := model.FlameConfig{
		Size:           model.SizeConfig{Width: 100, Height: 200},
		Seed:           456,
		IterationCount: 500,
		Output:         "temp.png",
		Threads:        4,
		EnableGamma:    true,
		Gamma:          1.8,
		SymmetryLevel:  2,
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
		},
		Variations: []model.Variation{
			{Name: "sinusoidal", Weight: 0.5, Apply: nil},
		},
	}

	data, err := json.Marshal(testConfig)
	require.NoError(t, err)
	_, err = tempFile.Write(data)
	require.NoError(t, err)
	tempFile.Close()

	// Загружаем в пустую конфигурацию
	cfg := &model.FlameConfig{}
	err = parseJSON(tempFile.Name(), cfg)
	require.NoError(t, err)

	// Проверяем, что поля заполнены
	assert.Equal(t, 100, cfg.Size.Width)
	assert.Equal(t, 200, cfg.Size.Height)
	assert.Equal(t, int64(456), cfg.Seed)
	assert.Equal(t, 500, cfg.IterationCount)
	assert.Equal(t, "temp.png", cfg.Output)
	assert.Equal(t, 4, cfg.Threads)
	assert.True(t, cfg.EnableGamma)
	assert.Equal(t, 1.8, cfg.Gamma)
	assert.Equal(t, 2, cfg.SymmetryLevel)
	assert.Len(t, cfg.AffineParams, 1)
	assert.Equal(t, 0.5, cfg.AffineParams[0].A)
	assert.Len(t, cfg.Variations, 1)
	assert.Equal(t, "sinusoidal", cfg.Variations[0].Name)
	assert.Equal(t, 0.5, cfg.Variations[0].Weight)
	// Apply не устанавливается в parseJSON, только в validate
	assert.Nil(t, cfg.Variations[0].Apply)
}

func TestParseJSON_AfterNotEmptyConfig(t *testing.T) {
	// Создаём временный JSON-файл
	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	testConfig := model.FlameConfig{
		Size:           model.SizeConfig{Width: 100, Height: 200},
		Seed:           456,
		IterationCount: 500,
		Output:         "temp.png",
		Threads:        4,
		EnableGamma:    true,
		Gamma:          1.8,
		SymmetryLevel:  2,
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
		},
		Variations: []model.Variation{
			{Name: "sinusoidal", Weight: 0.5, Apply: nil},
		},
	}

	data, err := json.Marshal(testConfig)
	require.NoError(t, err)
	_, err = tempFile.Write(data)
	require.NoError(t, err)
	tempFile.Close()

	// Загружаем в непустую конфигурацию
	cfg := &model.FlameConfig{
		Size:        model.SizeConfig{Width: 1080},
		EnableGamma: false,
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
			{A: 0.5, B: 0.5, C: 0.0, D: 0.5, E: 0.5, F: 0.0},
		},
		Variations: []model.Variation{
			{Name: "linear", Weight: 0.5, Apply: nil},
			{Name: "swirl", Weight: 0.8, Apply: nil},
		},
	}
	err = parseJSON(tempFile.Name(), cfg)
	require.NoError(t, err)

	// Проверяем, что поля заполнены так, как в JSON
	assert.Equal(t, 100, cfg.Size.Width)
	assert.Equal(t, 200, cfg.Size.Height)
	assert.Equal(t, int64(456), cfg.Seed)
	assert.Equal(t, 500, cfg.IterationCount)
	assert.Equal(t, "temp.png", cfg.Output)
	assert.Equal(t, 4, cfg.Threads)
	assert.True(t, cfg.EnableGamma)
	assert.Equal(t, 1.8, cfg.Gamma)
	assert.Equal(t, 2, cfg.SymmetryLevel)
	assert.Len(t, cfg.AffineParams, 1)
	assert.Equal(t, 0.5, cfg.AffineParams[0].A)
	assert.Len(t, cfg.Variations, 1)
	assert.Equal(t, "sinusoidal", cfg.Variations[0].Name)
	assert.Equal(t, 0.5, cfg.Variations[0].Weight)
	// Apply не устанавливается в parseJSON, только в validate
	assert.Nil(t, cfg.Variations[0].Apply)
}

func TestParseJSON_NotExist(t *testing.T) {
	tempFile := "tempfile.json"
	cfg := &model.FlameConfig{}
	err := parseJSON(tempFile, cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestParseJSON_NotJSON(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.md")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	cfg := &model.FlameConfig{}
	err = parseJSON(tempFile.Name(), cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must have .json extension")
}

func TestParseJSON_Invalid(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	data := []byte(`{"invalid json"}`)
	_, err = tempFile.Write(data)
	require.NoError(t, err)
	tempFile.Close()

	cfg := &model.FlameConfig{}
	err = parseJSON(tempFile.Name(), cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

// ---------Тесты для validate_cli.go --------

func TestValidateConfig(t *testing.T) {
	// Подготовим валидный конфиг
	cfg := &model.FlameConfig{
		Size:           model.SizeConfig{Width: 800, Height: 600},
		Seed:           123,
		IterationCount: 1000,
		Output:         "test.png",
		Threads:        2,
		EnableGamma:    true,
		Gamma:          2.0,
		SymmetryLevel:  3,
		AffineParams: []model.AffineTransform{
			{A: 0.5, B: 0.0, C: 0.0, D: 0.0, E: 0.5, F: 0.0},
		},
		Variations: []model.Variation{
			{Name: "linear", Weight: 1.0, Apply: nil}, // Apply будет установлен
		},
	}

	err := validateConfig(cfg)
	require.NoError(t, err)

	// Проверим, что Apply был установлен
	assert.NotNil(t, cfg.Variations[0].Apply)
}

func TestValidateInt_Width(t *testing.T) {
	cfg := &model.FlameConfig{Size: model.SizeConfig{Width: -1}}
	err := ValidateInt(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "width must be positive")
}

func TestValidateInt_Height(t *testing.T) {
	cfg := &model.FlameConfig{Size: model.SizeConfig{Width: 1, Height: 0}}
	err := ValidateInt(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "height must be positive")
}

func TestValidateInt_Iterations(t *testing.T) {
	cfg := &model.FlameConfig{Size: model.SizeConfig{Width: 1, Height: 1}, IterationCount: 0}
	err := ValidateInt(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "iteration count must be positive")
}

func TestValidateInt_Threads(t *testing.T) {
	cfg := &model.FlameConfig{Size: model.SizeConfig{Width: 1, Height: 1}, IterationCount: 2000, Threads: 0}
	err := ValidateInt(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "threads must be positive")
}

func TestValidateInt_Gamma(t *testing.T) {
	cfg := &model.FlameConfig{Size: model.SizeConfig{Width: 1, Height: 1}, IterationCount: 2000, Threads: 3, Gamma: 0}
	err := ValidateInt(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "gamma must be positive")
}

func TestValidateInt_Symmetry(t *testing.T) {
	cfg := &model.FlameConfig{Size: model.SizeConfig{Width: 1, Height: 1}, IterationCount: 2000, Threads: 3, Gamma: 2, SymmetryLevel: 0}
	err := ValidateInt(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symmetry level must be >= 1")
}

func TestValidateAffineParams_Empty(t *testing.T) {
	err := ValidateAffineParams([]model.AffineTransform{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "affine params must be non-empty")
}

func TestValidateAffineParams_InvalidCoeffs(t *testing.T) {
	affines := []model.AffineTransform{
		{A: 0.0, B: 0.0, C: 0.0, D: 0.0, E: 0.0, F: 0.0},
	}
	// Присваиваем NaN
	affines[0].A = math.NaN()
	err := ValidateAffineParams(affines)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid coefficients (NaN or Inf)")
}

func TestValidateFunctions_Empty(t *testing.T) {
	err := ValidateFunctions([]model.Variation{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variations must be non-empty")
}

func TestValidateFunctions_InvalidWeight(t *testing.T) {
	funcs := []model.Variation{
		{Name: "linear", Weight: 0, Apply: nil},
	}
	err := ValidateFunctions(funcs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "has invalid weight (NaN or Inf or non-positive)")
}

func TestValidateFunctions_InvalidName(t *testing.T) {
	funcs := []model.Variation{
		{Name: "nonexistent", Weight: 1.0, Apply: nil},
	}
	err := ValidateFunctions(funcs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not a valid function")

	// Apply не должен быть установлен
	assert.Nil(t, funcs[0].Apply)
}

func TestValidateOutputPath(t *testing.T) {
	// Создаём временную директорию для тестов
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		outputPath  string
		expectError bool
	}{
		{
			name:        "valid path in existing dir",
			outputPath:  filepath.Join(tempDir, "valid_output.png"),
			expectError: false,
		},
		{
			name:        "empty path",
			outputPath:  "",
			expectError: true,
		},
		{
			name:        "nonexistent directory",
			outputPath:  filepath.Join("nonexistent_dir", "file.png"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		// Пропускаем кейс, который тестируется отдельно
		if tt.name == "path to existing file as dir (edge case)" {
			continue
		}

		t.Run(tt.name, func(t *testing.T) {
			cfg := &model.FlameConfig{Output: tt.outputPath}
			err := ValidateOutput(cfg.Output)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ------- Тесты для проверки корректной приоритетности ---------

// createTestJSONFile создаёт временный JSON-файл с заданными данными для тестов
func createTestJSONFile(t *testing.T, data string) string {
	file, err := os.CreateTemp("", "test_config_*.json")
	assert.NoError(t, err)

	_, err = file.WriteString(data)
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.Remove(file.Name()) // Удаляем файл после завершения теста
	})

	return file.Name()
}

// TestCLIPriority_DefaultsUsedIfNoJSONNoCLI проверяет, что используются дефолтные значения, если нет JSON и CLI
func TestCLIPriority_DefaultsUsedIfNoJSONNoCLI(t *testing.T) {
	cmd := NewRootCommand()
	cmd.SetArgs([]string{
		// Никаких флагов, нет --config
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	cfg, err := GetConfig()
	assert.NoError(t, err)

	// Проверяем, что установлены дефолтные значения из model.FlameConfig.SetDefaults()
	defaultCfg := &model.FlameConfig{}
	defaultCfg.SetDefaults()

	assert.Equal(t, defaultCfg.Size.Width, cfg.Size.Width)
	assert.Equal(t, defaultCfg.Size.Height, cfg.Size.Height)
	assert.Equal(t, defaultCfg.IterationCount, cfg.IterationCount)
	assert.Equal(t, defaultCfg.Output, cfg.Output)
	assert.Equal(t, defaultCfg.Seed, cfg.Seed)
	assert.Equal(t, defaultCfg.Threads, cfg.Threads)
	assert.Equal(t, defaultCfg.EnableGamma, cfg.EnableGamma)
	assert.Equal(t, defaultCfg.Gamma, cfg.Gamma)
	assert.Equal(t, defaultCfg.SymmetryLevel, cfg.SymmetryLevel)
	// AffineParams и Variations также должны быть заполнены дефолтами
	assert.Len(t, cfg.AffineParams, len(defaultCfg.AffineParams))
	assert.Len(t, cfg.Variations, len(defaultCfg.Variations))
}

// TestCLIPriority_JSONWinsOverDefaults проверяет, что JSON имеет приоритет над дефолтами
func TestCLIPriority_JSONWinsOverDefaults(t *testing.T) {
	jsonContent := `{
		"size": { 
			"width": 800, 
			"height": 600
		},
		"iteration_count": 1000,
		"output_path": "json_default.png",
		"seed": 456,
		"threads": 2,
		"gamma_correction": true,
		"gamma": 2.0,
		"symmetry_level": 3,
		"affine_params": [
			{ "a": 0.8, "b": 0.1, "c": 0.1, "d": 0.1, "e": 0.8, "f": 0.1 }
		],
		"functions": [
			{ "name": "sinusoidal", "weight": 0.5 }
		]
	}`
	jsonPath := createTestJSONFile(t, jsonContent)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{
		"--config", jsonPath,
		// Не передаём CLI-флаги, только JSON
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	cfg, err := GetConfig()
	assert.NoError(t, err)

	// Проверяем, что значения из JSON перезаписали дефолты
	assert.Equal(t, 800, cfg.Size.Width)
	assert.Equal(t, 600, cfg.Size.Height)
	assert.Equal(t, 1000, cfg.IterationCount)
	assert.Equal(t, "json_default.png", cfg.Output)
	assert.Equal(t, int64(456), cfg.Seed)
	assert.Equal(t, 2, cfg.Threads)
	assert.Equal(t, true, cfg.EnableGamma)
	assert.Equal(t, 2.0, cfg.Gamma)
	assert.Equal(t, 3, cfg.SymmetryLevel)
	// Проверяем affine-params из JSON
	assert.Len(t, cfg.AffineParams, 1)
	assert.Equal(t, 0.8, cfg.AffineParams[0].A)
	// Проверяем functions из JSON
	assert.Len(t, cfg.Variations, 1)
	assert.Equal(t, "sinusoidal", cfg.Variations[0].Name)
	assert.Equal(t, 0.5, cfg.Variations[0].Weight)
}

// TestCLIPriority_CLIWins проверяет, что CLI-флаги имеют наивысший приоритет
func TestCLIPriority_CLIWins(t *testing.T) {
	jsonContent := `{
		"size": { "width": 100, "height": 100 },
		"iteration_count": 100,
		"output_path": "json_output.png",
		"seed": 123,
		"threads": 1,
		"gamma_correction": true,
		"gamma": 1.8,
		"symmetry_level": 2,
		"affine_params": [
			{ "a": 1.0, "b": 0.0, "c": 0.0, "d": 0.0, "e": 1.0, "f": 0.0 }
		],
		"functions": [
			{ "name": "linear", "weight": 1.0 }
		]
	}`
	jsonPath := createTestJSONFile(t, jsonContent)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{
		"--config", jsonPath,
		"--width", "500",
		"--iteration-count", "5000",
		"-o", "cli_output.png",
		"--seed", "999",
		"--threads", "4",
		"--gamma-correction=false", // CLI отключает гамму, JSON включает
		"--gamma", "2.5",
		"--symmetry-level", "5",
		"--affine-params", "0.5,0,0,0,0.5,0.5", // CLI переопределяет
		"--functions", "sinusoidal:1.0", // CLI переопределяет
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	cfg, err := GetConfig()
	assert.NoError(t, err)

	// Проверяем, что CLI-значения перезаписали JSON
	assert.Equal(t, 500, cfg.Size.Width)
	assert.Equal(t, 5000, cfg.IterationCount)
	assert.Equal(t, "cli_output.png", cfg.Output)
	assert.Equal(t, int64(999), cfg.Seed)
	assert.Equal(t, 4, cfg.Threads)
	// CLI указал gamma-correction=false, хотя в JSON было true
	assert.Equal(t, false, cfg.EnableGamma)
	// gamma=2.5 из CLI, хотя в JSON было 1.8
	assert.Equal(t, 2.5, cfg.Gamma)
	assert.Equal(t, 5, cfg.SymmetryLevel)
	// CLI переопределило affine-params
	assert.Len(t, cfg.AffineParams, 1)
	assert.Equal(t, 0.5, cfg.AffineParams[0].A)
	// CLI переопределило functions
	assert.Len(t, cfg.Variations, 1)
	assert.Equal(t, "sinusoidal", cfg.Variations[0].Name)
}

// TestCLIPriority_MixedPriority проверяет смешанный сценарий: часть из JSON, часть из CLI, часть дефолт
func TestCLIPriority_MixedPriority(t *testing.T) {
	jsonContent := `{
		"size": { "width": 123, "height": 456 },
		"iteration_count": 999,
		"output_path": "json_default.png",
		"seed": 789,
		"threads": 3,
		"gamma_correction": true,
		"gamma": 1.9,
		"symmetry_level": 4,
		"affine_params": [
			{ "a": 0.7, "b": 0.2, "c": 0.1, "d": 0.2, "e": 0.7, "f": 0.1 }
		],
		"functions": [
			{ "name": "spherical", "weight": 0.8 }
		]
	}`

	jsonPath := createTestJSONFile(t, jsonContent)

	cmd := NewRootCommand()
	cmd.SetArgs([]string{
		"--config", jsonPath,
		"--width", "789", // CLI переопределяет JSON
		"--threads", "8", // CLI переопределяет JSON
		// iteration_count не переопределяется CLI -> остаётся из JSON
		// output_path не переопределяется CLI -> остаётся из JSON
		// остальные не переопределяются -> остаются из JSON или становятся дефолтными
	})

	err := cmd.Execute()
	assert.NoError(t, err)

	cfg, err := GetConfig()
	assert.NoError(t, err)

	// CLI: width=789, threads=8
	assert.Equal(t, 789, cfg.Size.Width)

	assert.Equal(t, 8, cfg.Threads)
	// JSON: остальные
	assert.Equal(t, 456, cfg.Size.Height)           // из JSON
	assert.Equal(t, 999, cfg.IterationCount)        // из JSON
	assert.Equal(t, "json_default.png", cfg.Output) // из JSON
	assert.Equal(t, int64(789), cfg.Seed)           // из JSON
	assert.Equal(t, true, cfg.EnableGamma)          // из JSON
	assert.Equal(t, 1.9, cfg.Gamma)                 // из JSON
	assert.Equal(t, 4, cfg.SymmetryLevel)           // из JSON
	// Affine и Functions из JSON
	assert.Len(t, cfg.AffineParams, 1)
	assert.Equal(t, 0.7, cfg.AffineParams[0].A)
	assert.Len(t, cfg.Variations, 1)
	assert.Equal(t, "spherical", cfg.Variations[0].Name)
	assert.Equal(t, 0.8, cfg.Variations[0].Weight)
}
