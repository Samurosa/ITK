package main

import (
	"flag"
	"os"

	"ITK_Code/m/v2/internal/app"

	"go.uber.org/zap"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	cfgPath := flag.String(
		"config",
		"./internal/config/local.yaml",
		"config file path",
	)

	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		return 1
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			// Log synchronization is best-effort.
			_ = err
		}
	}()

	if err := app.Run(*cfgPath, logger); err != nil {
		logger.Error(
			"spot server failed",
			zap.Error(err),
		)

		return 1
	}

	return 0
}
