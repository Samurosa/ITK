package interceptors

import (
	reqCtx "ITK_Code/m/v2/internal/core/context"
	"ITK_Code/m/v2/internal/core/dto"
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func ClientIPInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		log := log.Named("client-ip-interceptor")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Info("ip from metadata not found")
		}

		ip := md.Get("x-forwarded-for")[0]

		if ip == "" {
			p, ok := peer.FromContext(ctx)
			if !ok {
				log.Error("peer not found")
				return nil, status.Error(codes.Internal, "client ip not found")
			}

			host, _, err := net.SplitHostPort(p.Addr.String())
			if err != nil {
				log.Error("invalid peer address", zap.Error(err))
				return nil, status.Error(codes.Internal, "invalid client ip")
			}
			log.Info("peer address received")

			ip = host
		}

		ctx, err := reqCtx.UpdateRequestContext(ctx, func(rc *dto.RequestContext) {
			rc.Metadata.ClientIP = ip
		})
		if err != nil {
			log.Error("update request context", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to update request context")
		}

		return handler(ctx, req)
	}
}
