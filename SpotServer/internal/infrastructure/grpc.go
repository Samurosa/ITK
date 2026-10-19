package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"net"

	"ITK_Code/m/v2/internal/adapters/inbound/grpc/interceptors"
	spGRPC "ITK_Code/m/v2/internal/adapters/inbound/grpc/server"
	"ITK_Code/m/v2/internal/application"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/authentication"
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
	serveDone  chan struct{}
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

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.LoggerInterceptor(log),
			recovery.RecoveryInterceptor(),
			authentication.AuthInterceptor(
				parser,
				publicMethods,
				validator,
			),
			interceptors.RequireAdmin(publicMethods),
		),
	)

	spGRPC.RegisterSpotService(grpcServer, spotService)

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
		return fmt.Errorf(
			"listen on grpc port %d: %w",
			a.port,
			err,
		)
	}

	a.serveDone = make(chan struct{})

	go func() {
		defer close(a.serveDone)

		serveErr := a.grpcServer.Serve(listener)

		if serveErr != nil &&
			!errors.Is(serveErr, grpc.ErrServerStopped) {
			if onFailure != nil {
				onFailure(
					fmt.Errorf("serve grpc: %w", serveErr),
				)
			}
		}
	}()

	a.log.Info(
		"grpc spot server listening",
		zap.Int("port", a.port),
	)

	return nil
}

func (a *GRPCApp) Stop(ctx context.Context) error {
	if a.serveDone == nil {
		return nil
	}

	a.log.Info("grpc spot server stopping")

	stopped := make(chan struct{})

	go func() {
		defer close(stopped)
		a.grpcServer.GracefulStop()
	}()

	select {
	case <-stopped:
		a.log.Info(
			"grpc spot server gracefully stopped",
		)

	case <-ctx.Done():
		a.log.Warn(
			"grpc graceful shutdown timed out; forcing stop",
			zap.Error(ctx.Err()),
		)

		a.grpcServer.Stop()

		a.log.Info(
			"grpc spot server force stopped",
		)

		select {
		case <-stopped:
		default:
			// GracefulStop is still completing.
		}
	}

	select {
	case <-a.serveDone:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}
