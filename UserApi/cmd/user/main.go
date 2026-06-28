package main

import (
	grpsApp "ITK_Code/m/v2/app"
	"ITK_Code/m/v2/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load("config/local.yaml")
	err := godotenv.Load(cfg.Env)
	if err != nil {
		log.Fatalf("Error loading .env file path:%s", cfg.Env)
	}

	secret := os.Getenv("JWT_SECRET")

	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("error create logger: ", zap.Error(err))
	}

	defer func(log *zap.Logger) {
		err := log.Sync()
		if err != nil {
			log.Fatal("error sync logger: ", zap.Error(err))
		}
	}(logger)

	application := grpsApp.New(logger, cfg.GRPC.Port, cfg.TokenTTl, secret)

	go application.GrpcApp.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GrpcApp.Stop()

	logger.Debug("application stop")

}
