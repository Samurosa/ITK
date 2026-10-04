package infrastructure

import (
	"ITK_Code/m/v2/internal/adapters/inbound/grpc/interceptors"
	spGRPC "ITK_Code/m/v2/internal/adapters/inbound/grpc/server"
	"ITK_Code/m/v2/internal/core/spot"

	"fmt"
	"net"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/authentication"
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
	spotService spot.Service,
	port int,
	parser *jwt.Parser,
	validator session.Validator,
) *GRPCApp {
	publicMethods := map[string]struct{}{
		spotpb.SpotInstrumentService_GetSpot_FullMethodName:   {},
		spotpb.SpotInstrumentService_ListSpots_FullMethodName: {},
	}
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		authentication.AuthInterceptor(log, parser, publicMethods, validator),
		interceptors.RequireAdmin(publicMethods),
	))

	spGRPC.RegisterSpotService(grpcServer,
		spotService,
		log,
	)

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
		"grpcs Spot server started",
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
		a.log.Info("GRPC Spot server gracefully stopped")

	case <-time.After(10 * time.Second):
		a.log.Info("GRPC Spot server timeout")
		a.grpcServer.Stop()
	}
}
