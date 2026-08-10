package infrastructure

import (
	usergrps "ITK_Code/m/v2/internal/adapters/inbound/grps"
	"ITK_Code/m/v2/internal/adapters/inbound/grps/interceptors"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Services interface {
	TokenManager() auth.TokenManager
	UserService() user.Service
	AuthService() auth.Service
	WalletService() wallet.Service
	SessionStorage() auth.SessionRepository
}

type GRPCApp struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	port       int
}

func NewGRPC(
	log *zap.Logger,
	services Services,
	port int,
) *GRPCApp {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestContextInterceptor(log),
			interceptors.ClientIPInterceptor(log),
			interceptors.AuthInterceptor(log,
				services.TokenManager(),
				services.SessionStorage(),
			),
		),
	)

	usergrps.RegisterUserService(grpcServer,
		services.UserService(),
		services.AuthService(),
		services.WalletService(),
		log)

	return &GRPCApp{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *GRPCApp) Run() error {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return err
	}

	a.log.Info(
		"grpcs UserServer started",
		zap.Any("port", a.port),
	)

	if err := a.grpcServer.Serve(l); err != nil {
		return err
	}
	return nil
}

func (a *GRPCApp) Stop() {

	done := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()
	select {

	case <-done:
		a.log.Info("GRPC UserServer gracefully stopped")

	case <-time.After(10 * time.Second):
		a.log.Info("GRPC UserServer timeout")
		a.grpcServer.Stop()
	}
}
