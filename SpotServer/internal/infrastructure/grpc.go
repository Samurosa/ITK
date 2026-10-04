package infrastructure

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/interceptors"
	spGRPC "ITK_Code/m/v2/internal/adapters/inbound/grpc/server"
	"ITK_Code/m/v2/internal/application"

	"fmt"
	"net"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/interceptors/recovery"
	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	spotpb "github.com/Samurosa/exchange-contract/protobuf/gen/go/spot"

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
	spotService *application.Spot,
	port int,
	parser *jwt.Parser,
	validator session.Validator,
) *GRPCApp {
	publicMethods := map[string]struct{}{
		spotpb.SpotInstrumentService_GetSpot_FullMethodName:   {},
		spotpb.SpotInstrumentService_ListSpots_FullMethodName: {},
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.LoggerInterceptor(log),
		interceptors.RequestLogging(),
		recovery.RecoveryInterceptor(),
		interceptors.Authenticate(parser, publicMethods, validator),
		interceptors.RequireAdmin(publicMethods),
	))

	spGRPC.RegisterSpotService(grpcServer, spotService)

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
		"grpc spot server listening",
		zap.Int("port", a.port),
	)

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("serve grpc: %w", err)
	}
	return nil
}

func (a *GRPCApp) Stop() {
	a.log.Info("grpc spot server stopping")

	done := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()
	select {

	case <-done:
		a.log.Info("grpc spot server gracefully stopped")

	case <-time.After(10 * time.Second):
		a.log.Warn("grpc graceful shutdown timed out; forcing stop", zap.Duration("timeout", 10*time.Second))
		a.grpcServer.Stop()
		a.log.Info("grpc spot server force stopped")
	}
}
