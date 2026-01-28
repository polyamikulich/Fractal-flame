package main

import (
	"log/slog"
	"os"

	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/application"
	"gitlab.education.tbank.ru/backend-academy-go-2025/homeworks/hw4-fractal-flame/internal/application/cli"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	rootCmd := cli.NewRootCommand()

	if err := rootCmd.Execute(); err != nil {
		slog.Error("CLI error", "error", err)
		os.Exit(2)
	}

	cfg, err := cli.GetConfig()
	if err != nil {
		slog.Error("Failed to load config after CLI execution", "error", err)
		os.Exit(1)
	}

	err = application.Run(cfg)
	if err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}
