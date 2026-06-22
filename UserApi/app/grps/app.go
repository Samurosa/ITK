package grpsApp

import (
	usergrps "ITK_Code/m/v2/internal/adapters/grps/user"
	"fmt"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	port       int
}

func New(
	log *zap.Logger,
	userService usergrps.User,
	port int,
) *App {
	grpcServer := grpc.NewServer()

	usergrps.RegisterUserService(grpcServer, userService)

	return &App{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *App) Run() {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		a.log.Fatal("error accessing port:", zap.Error(err))
	}

	a.log.Info(
		"grpcs server started",
		zap.Any("port", a.port),
	)

	if err := a.grpcServer.Serve(l); err != nil {
		a.log.Fatal("grpc server failed", zap.Error(err))
	}
}

func (a *App) Stop() {
	a.log.Info("grpc server stopped")
	a.grpcServer.GracefulStop()
}

// проверять токены
