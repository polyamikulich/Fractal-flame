package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/domain/model"
)

// cliConfig — глобальная переменная для хранения сконфигурированных значений.
var cliConfig *model.FlameConfig

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "fractal-flame [flags]",
		Short:        "Generate a flame fractal image",
		Long:         "Generate a flame fractal image",
		Args:         cobra.ArbitraryArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Начинаем с дефолтной конфигурации (3-ий приоритет)
			config := &model.FlameConfig{}
			config.SetDefaults()
			// printConfig(config)

			// Проверка наличия конфигурационного файла, чтение оттуда (2-ой приоритет)
			pathConfig, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			if pathConfig != "" {
				if err := parseJSON(pathConfig, config); err != nil {
					return err
				}
			}

			// printConfig(config)

			// Проверка аргументов командной строки (1-ый приоритет)
			if err := parseConfig(cmd, config); err != nil {
				return err
			}

			if err := validateConfig(config); err != nil {
				return err
			}

			cliConfig = config
			// printConfig(config)
			return nil
		},
	}

	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		printMainHelp(cmd)
	})

	// Выставляем флаги
	// Важно: дефолтно ставим 0 для int и "" для string - нужно для учёта приоритетов
	// Функция высставления правильных дефолтов лежит в domain\flame.go
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
	cmd.Flags().StringP("affine-params", "a", "", "Affine parameters: <a1>,<b1>,<c1>,<d1>,<e1>,<f1>/.../<an>,<bn>,<cn>,<dn>,<en>,<fn>; default: params for Sierpinski's Carpet")
	cmd.Flags().StringP("functions", "f", "", "Variations and weights: v1:w1,...,vk:wk; default: linear:1.0")

	return cmd
}

// GetConfig возвращает сконфигурированную структуру Config после выполнения команды.
// Вызывать только после cmd.Execute().
func GetConfig() (*model.FlameConfig, error) {
	if cliConfig == nil {
		return nil, fmt.Errorf("configuration not loaded, CLI execution may have failed")
	}
	return cliConfig, nil
}

// printConfig выводит конфиг для отладки
// func printConfig(cfg *model.FlameConfig) {
// 	fmt.Fprintf(os.Stderr, "✅ Config loaded:\n")
// 	fmt.Fprintf(os.Stderr, "   width: %d\n", cfg.Size.Width)
// 	fmt.Fprintf(os.Stderr, "   height: %d\n", cfg.Size.Height)
// 	fmt.Fprintf(os.Stderr, "   Seed: %d\n", cfg.Seed)
// 	fmt.Fprintf(os.Stderr, "   Output: %s\n", cfg.Output)
// 	fmt.Fprintf(os.Stderr, "   Threads: %d\n", cfg.Threads)
// 	fmt.Fprintf(os.Stderr, "   Iterations: %d\n", cfg.IterationCount)
// 	fmt.Fprintf(os.Stderr, "   Gamma: %.2f\n", cfg.Gamma)
// 	fmt.Fprintf(os.Stderr, "   Enable gamma: %v\n", cfg.EnableGamma)
// 	fmt.Fprintf(os.Stderr, "   Symmetry level: %d\n", cfg.SymmetryLevel)
// 	fmt.Fprintf(os.Stderr, "   Affine params: %v\n", cfg.AffineParams)
// 	fmt.Fprintf(os.Stderr, "   Variations (%d):\n", len(cfg.Variations))
// 	for i, v := range cfg.Variations {
// 		fmt.Fprintf(os.Stderr, "     [%d] %s (weight: %.2f), func: %v\n", i+1, v.Name, v.Weight, v.Apply != nil)
// 	}
// }

func printMainHelp(cmd *cobra.Command) {
	fmt.Print("Usage: fractal-flame [OPTIONS]\r\n")
	fmt.Printf("%s", cmd.Long)
	fmt.Print("\r\n")
	fmt.Print("Options:\r\n")

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		var flagLine string
		if flag.Shorthand != "" {
			flagLine = fmt.Sprintf("  -%s, --%-20s", flag.Shorthand, flag.Name)
		} else {
			flagLine = fmt.Sprintf("      --%-20s", flag.Name)
		}

		flagLine += fmt.Sprintf("   \t %-s", flag.Usage)

		if flag.DefValue != "" && flag.DefValue != "false" {
			flagLine += fmt.Sprintf(" (default \"%s\")", flag.DefValue)
		}

		fmt.Print(flagLine + "\r\n")
	})
}
