package main

import (
	"ITK_Code/m/v2/internal/app"
	"ITK_Code/m/v2/internal/config"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	cfgPath := "intenal/config/local.yaml"
	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Printf("error loading config file path: %s, error: %s", cfgPath, err.Error())
		return
	}
	if err := godotenv.Load(cfg.Env); err != nil {
		fmt.Printf("error loading .env file path: %s, error: %s", cfg.Env, err.Error())
		return
	}

	secret := os.Getenv("JWT_SECRET")

	application, err := app.New(
		cfg,
		secret,
	)
	if err != nil {
		fmt.Printf("create application failed error: %s", err.Error())
		return
	}

	application.Start()

	application.WaitSignal()

	application.Stop()
}
