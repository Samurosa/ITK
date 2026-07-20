package infrastructure

import (
	usergrps "ITK_Code/m/v2/internal/adapters/inbound/grps/user"
	"ITK_Code/m/v2/internal/adapters/inbound/interceptors"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"ITK_Code/m/v2/internal/core/wallet"
	"fmt"
	"net"

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
		grpc.UnaryInterceptor(interceptors.AuthInterceptor(log,
			services.TokenManager(),
			services.SessionStorage(),
		)),
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

func (a *GRPCApp) Run() {

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

func (a *GRPCApp) Stop() {
	a.log.Info("GRPC server stopped")
	a.grpcServer.GracefulStop()
}
