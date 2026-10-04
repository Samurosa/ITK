package interceptors

import (
	"context"
	"net"
	"strings"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func ClientIPInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		var ip string

		log := logging.FromContext(ctx).Named("client-ip-interceptor")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Debug("metadata not found")
		}

		clientIP := md.Get("x-forwarded-for")

		if len(clientIP) > 0 {
			ip = strings.Split(clientIP[0], ",")[0]
		}

		if ip == "" {
			p, ok := peer.FromContext(ctx)
			if !ok || p == nil || p.Addr == nil {
				log.Error("peer not found")
				return nil, status.Error(codes.Internal, "client ip not found")
			}

			host, _, err := net.SplitHostPort(p.Addr.String())
			if err != nil {
				log.Error("invalid peer address", zap.Error(err))
				return nil, status.Error(codes.Internal, "invalid client ip")
			}
			log.Debug("peer address received")

			ip = host
		}

		ctx = sharedContext.UpdateRequestContext(ctx, func(rc *sharedContext.RequestContext) {
			rc.Metadata.ClientIP = ip
		})

		return handler(ctx, req)
	}
}
