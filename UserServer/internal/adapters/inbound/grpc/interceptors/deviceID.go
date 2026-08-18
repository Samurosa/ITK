package interceptors

import (
	reqCtx "ITK_Code/m/v2/internal/core/context"
	"ITK_Code/m/v2/internal/core/dto"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DeviceIDInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		var device string

		log := log.Named("device-id-interceptor")

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Debug("metadata not found")
		}

		deviceID := md.Get("device-id")

		if len(deviceID) > 0 {
			device = deviceID[0]
		}

		if device == "" {
			log.Error("device not found in metadata")
			return nil, nil
		}

		ctx, err := reqCtx.UpdateRequestContext(ctx, func(rc *dto.RequestContext) {
			rc.Metadata.DeviceID = device
		})
		if err != nil {
			log.Error("update request context", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to update request context")
		}

		return handler(ctx, req)
	}
}
