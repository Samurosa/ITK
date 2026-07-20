package main

import (
	"ITK_Code/m/v2/config"
	"ITK_Code/m/v2/internal/app"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	cfg := config.Load("config/local.yaml")

	err := godotenv.Load(cfg.Env)
	if err != nil {
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
