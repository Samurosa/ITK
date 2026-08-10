package main

import (
	"flag"
	"fmt"
	"os"

	"ITK_Code/m/v2/internal/app"
	"ITK_Code/m/v2/internal/config"
)

func main() {
	cfgPath := flag.String(
		"config",
		"./internal/config/local.yaml",
		"config file path",
	)

	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Printf(
			"error loading config file path: %s, error: %s\n",
			*cfgPath,
			err,
		)
		return
	}

	secret, err := os.ReadFile(cfg.JWTSecretPath)
	if err != nil {
		fmt.Printf(
			"error reading JWT secret from %s: %s\n",
			cfg.JWTSecretPath,
			err,
		)
		return
	}

	application, err := app.New(
		cfg,
		string(secret),
	)
	if err != nil {
		fmt.Printf(
			"create application failed: %s\n",
			err,
		)
		return
	}

	application.Start()

	application.WaitSignal()

	application.Stop()
}
