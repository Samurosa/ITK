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
		_ = logger.Sync()
	}()

	if err := app.Run(*cfgPath, logger); err != nil {
		logger.Error(
			"order server failed",
			zap.Error(err),
		)

		return 1
	}

	return 0
}
