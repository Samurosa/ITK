package infrastructure

import (
	spGRPC "ITK_Code/m/v2/internal/adapters/inbound/grpc"
	"ITK_Code/m/v2/internal/core/spot"

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
	spotService spot.Service,
	port int,
) *GRPCApp {
	grpcServer := grpc.NewServer()

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
