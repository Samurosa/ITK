package app

import (
	grpsApp "ITK_Code/m/v2/internal/app/grps"
	"time"

	"go.uber.org/zap"
)

type App struct {
	GrpcApp *grpsApp.App
}

func New(
	log *zap.Logger,
	port int,
	tokenTTL time.Duration,
) *App {
	app := grpsApp.New(log, port)

	return &App{
		GrpcApp: app,
	}
}
