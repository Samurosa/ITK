package infrastructure

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/interceptors"
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/server"
	"ITK_Code/m/v2/internal/application"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/interceptors/recovery"
	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"

	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	log        *zap.Logger
	grpcServer *grpc.Server
	port       int
}

func NewGRPC(
	log *zap.Logger,
	orderService *application.OrderService,
	jwtParser *jwt.Parser,
	validator session.Validator,
	port int,
) *GRPCApp {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.LoggerInterceptor(log),
		interceptors.RequestLog(),
		recovery.RecoveryInterceptor(),
		interceptors.Authentication(jwtParser, validator),
	), grpc.ChainStreamInterceptor(
		logging.LoggerStreamInterceptor(log),
		interceptors.StreamLog(),
		recovery.RecoveryStreamInterceptor(),
		interceptors.StreamAuthentication(jwtParser, validator),
	))

	server.NewOrderServer(grpcServer, orderService)

	return &GRPCApp{
		log:        log,
		grpcServer: grpcServer,
		port:       port,
	}
}

func (a *GRPCApp) Run() error {

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("listen on grpc port %d: %w", a.port, err)
	}

	a.log.Info(
		"grpc order server listening",
		zap.String("address", l.Addr().String()),
	)

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("serve grpc: %w", err)
	}
	return nil
}

func (a *GRPCApp) Stop() {
	a.log.Info("grpc order server stopping")

	done := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()
	select {

	case <-done:
		a.log.Info("grpc order server gracefully stopped")

	case <-time.After(10 * time.Second):
		a.log.Warn("grpc graceful shutdown timed out; forcing stop", zap.Duration("timeout", 10*time.Second))
		a.grpcServer.Stop()
		a.log.Info("grpc order server stopped after forced shutdown")
	}
}
