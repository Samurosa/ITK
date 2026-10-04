package interceptors

import (
	"context"
	"errors"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RequestLogInterceptor records the outcome, including early interceptor failures.
// It runs after LoggerInterceptor and before recovery and validation.
func RequestLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		started := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err)
		switch {
		case errors.Is(err, context.Canceled):
			code = codes.Canceled
		case errors.Is(err, context.DeadlineExceeded):
			code = codes.DeadlineExceeded
		}

		level := zapcore.WarnLevel
		switch code {
		case codes.OK, codes.Canceled, codes.NotFound, codes.AlreadyExists, codes.Aborted:
			level = zapcore.DebugLevel
		case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
			level = zapcore.ErrorLevel
		}
		// Do not serialize requests, responses or status messages: they may contain credentials.
		logging.FromContext(ctx).Log(level, "grpc request completed",
			zap.String("grpc_code", code.String()),
			zap.Duration("duration", time.Since(started)),
		)
		return resp, err
	}
}
