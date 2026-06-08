package main

import (
	config "ITK_Code/m/v2/config"
	"os"
	"os/signal"
	"syscall"

	grpsApp "ITK_Code/m/v2/internal/app"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load("./config/local.yaml")

	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("error create logger: ", zap.Error(err))
	}
	defer log.Sync()

	application := grpsApp.New(log, cfg.GRPC.Port, cfg.Token_ttl)

	go application.GrpcApp.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.GrpcApp.Stop()

	log.Debug("application stop")

}
