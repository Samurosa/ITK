package interceptors

import (
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DeviceIDInterceptor(requiredMethods map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := requiredMethods[info.FullMethod]; !ok {
			return handler(ctx, req)
		}

		log := logging.FromContext(ctx).Named("device-id-interceptor")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Debug("metadata not found")
			return nil, status.Error(codes.InvalidArgument, "metadata not found")
		}

		deviceID := md.Get("device-id")
		if len(deviceID) > 1 {
			log.Debug("multiple device-id metadata values")
			return nil, status.Error(codes.InvalidArgument, "device-id must have one value")
		}

		if len(deviceID) > 0 {
			return nil, status.Error(codes.InvalidArgument, "device-id metadata must have empty value")
		}

		if deviceID[0] == "" {
			log.Debug("device not found in metadata")
			return nil, status.Error(codes.InvalidArgument, "device-id must have empty value")
		}

		ctx = sharedContext.UpdateRequestContext(ctx, func(rc *sharedContext.RequestContext) {
			rc.Metadata.DeviceID = deviceID[0]
		})

		return handler(ctx, req)
	}
}
