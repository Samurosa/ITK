package infrastructure

import (
	usergrps "ITK_Code/m/v2/internal/adapters/inbound/grpc"

	"ITK_Code/m/v2/internal/adapters/inbound/grpc/interceptors"
	"fmt"
	"net"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/interceptors/recovery"
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
	user usergrps.UserService,
	auth usergrps.AuthService,
	port int,
	tokenParser *sharedjwt.Parser,
	sessionValidator sharedsession.Validator,
) *GRPCApp {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.LoggerInterceptor(log),
			interceptors.RequestLogInterceptor(),
			recovery.RecoveryInterceptor(),
			interceptors.ClientIPInterceptor(),
			interceptors.DeviceIDInterceptor(),
			interceptors.AuthInterceptor(
				tokenParser,
				publicMethods,
				sessionValidator,
			),
		),
	)

	usergrps.RegisterUserService(grpcServer, user, auth)

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
		a.log.Warn("grpc graceful shutdown timed out; forcing stop")
		a.grpcServer.Stop()
		a.log.Info("grpc UserServer stopped")
	}
}
