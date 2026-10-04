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

func DeviceIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		var device string

		log := logging.FromContext(ctx).Named("device-id-interceptor")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Debug("metadata not found")
			return handler(ctx, req)
		}

		deviceID := md.Get("device-id")
		if len(deviceID) > 1 {
			log.Debug("multiple device-id metadata values")
			return nil, status.Error(codes.InvalidArgument, "device-id must have one value")
		}

		if len(deviceID) > 0 {
			device = deviceID[0]
		}

		if device == "" {
			log.Debug("device not found in metadata")
			return handler(ctx, req)
		}

		ctx = sharedContext.UpdateRequestContext(ctx, func(rc *sharedContext.RequestContext) {
			rc.Metadata.DeviceID = device
		})

		return handler(ctx, req)
	}
}
