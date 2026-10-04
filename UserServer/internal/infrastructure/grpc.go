package infrastructure

import (
	usergrps "ITK_Code/m/v2/internal/adapters/inbound/grpc"

	"ITK_Code/m/v2/internal/adapters/inbound/grpc/interceptors"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"fmt"
	"net"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/authentication"
	sharedjwt "github.com/Samurosa/exchange-common/shared/auth/jwt"
	sharedsession "github.com/Samurosa/exchange-common/shared/auth/session"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var publicMethods = map[string]struct{}{
	"/user.UserService/Login":        {},
	"/user.UserService/Registration": {},
	"/user.UserService/RefreshToken": {},
}

type GRPCApp struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	port       int
}

func NewGRPC(
	log *zap.Logger,
	user user.Service,
	auth auth.Service,
	port int,
	tokenParser *sharedjwt.Parser,
	sessionValidator sharedsession.Validator,
) *GRPCApp {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.ClientIPInterceptor(log),
			interceptors.DeviceIDInterceptor(log),
			authentication.AuthInterceptor(
				log,
				tokenParser,
				publicMethods,
				sessionValidator,
			),
		),
	)

	usergrps.RegisterUserService(grpcServer, user, auth, log)

	return &GRPCApp{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *GRPCApp) Run() error {
	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", a.port),
	)
	if err != nil {
		return err
	}

	a.log.Info(
		"grpc UserServer started",
		zap.Int("port", a.port),
	)

	return a.grpcServer.Serve(listener)
}

func (a *GRPCApp) Stop() {
	done := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		a.log.Info("grpc UserServer gracefully stopped")
	case <-time.After(10 * time.Second):
		a.log.Info("grpc UserServer stop timeout")
		a.grpcServer.Stop()
	}
}
