package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"net"

	usergrps "ITK_Code/m/v2/internal/adapters/inbound/grpc"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/interceptors"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/authentication"
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

var requiredMethods = map[string]struct{}{
	"/user.UserService/Login": {},
}

type GRPCApp struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	port       int
	serveDone  chan struct{}
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
			recovery.RecoveryInterceptor(),
			interceptors.ClientIPInterceptor(),
			interceptors.DeviceIDInterceptor(requiredMethods),
			authentication.AuthInterceptor(
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

func (a *GRPCApp) Start(
	ctx context.Context,
	onFailure func(error),
) error {
	var listenConfig net.ListenConfig

	listener, err := listenConfig.Listen(
		ctx,
		"tcp",
		fmt.Sprintf(":%d", a.port),
	)
	if err != nil {
		return fmt.Errorf("start grpc listener: %w", err)
	}

	a.serveDone = make(chan struct{})

	go func() {
		defer close(a.serveDone)

		serveErr := a.grpcServer.Serve(listener)
		if serveErr != nil &&
			!errors.Is(serveErr, grpc.ErrServerStopped) {
			if onFailure != nil {
				onFailure(fmt.Errorf("grpc serve: %w", serveErr))
			}
		}
	}()

	a.log.Info(
		"grpc UserServer started",
		zap.Int("port", a.port),
	)

	return nil
}

func (a *GRPCApp) Stop(ctx context.Context) error {
	if a.serveDone == nil {
		return nil
	}

	stopped := make(chan struct{})

	go func() {
		defer close(stopped)
		a.grpcServer.GracefulStop()
	}()

	select {
	case <-stopped:
		a.log.Info("grpc server gracefully stopped")

	case <-ctx.Done():
		a.log.Warn(
			"grpc shutdown deadline reached; forcing stop",
			zap.Error(ctx.Err()),
		)

		a.grpcServer.Stop()

		<-stopped

		a.log.Info("grpc server forcibly stopped")
	}

	select {
	case <-a.serveDone:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
