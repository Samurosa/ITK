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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secret := os.Getenv("JWT_SECRET")
	configPath := os.Getenv("CONFIG_PATH")

	cfg := config.Load(configPath)

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

	application := grpsApp.New(logger, cfg.GRPC.Port, cfg.Token_ttl, secret)

	go application.GrpcApp.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GrpcApp.Stop()

	logger.Debug("application stop")

}
