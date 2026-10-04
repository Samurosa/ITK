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

// RequestLogging runs after LoggerInterceptor and before recovery and auth so
// rejected requests and recovered panics also have a completion record.
func RequestLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		started := time.Now()
		response, err := handler(ctx, req)
		code := status.Code(err)
		switch {
		case errors.Is(err, context.Canceled):
			code = codes.Canceled
		case errors.Is(err, context.DeadlineExceeded):
			code = codes.DeadlineExceeded
		}

		fields := []zap.Field{
			zap.String("grpc_code", code.String()),
			zap.Duration("duration", time.Since(started)),
		}
		if err != nil {
			fields = append(fields, zap.Error(err))
		}
		logging.FromContext(ctx).Log(requestLogLevel(code), "grpc request completed", fields...)
		return response, err
	}
}

func requestLogLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK, codes.NotFound, codes.Canceled:
		return zap.DebugLevel
	case codes.InvalidArgument, codes.OutOfRange, codes.AlreadyExists,
		codes.FailedPrecondition, codes.Unauthenticated, codes.PermissionDenied,
		codes.ResourceExhausted, codes.DeadlineExceeded, codes.Aborted:
		return zap.WarnLevel
	default:
		return zap.ErrorLevel
	}
}
