package main

import (
	"ITK_Code/m/v2/config"
	"ITK_Code/m/v2/internal/app"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	cfg := config.Load("config/local.yaml")

	err := godotenv.Load(cfg.Env)
	if err != nil {
		panic("error loading .env file path: " + cfg.Env + "error: " + err.Error())
	}

	secret := os.Getenv("JWT_SECRET")

	application, err := app.New(
		cfg,
		secret,
	)

	if err != nil {
		panic("create application failed error: " + err.Error())
	}

	application.Start()

	application.WaitSignal()

	application.Stop()
}
