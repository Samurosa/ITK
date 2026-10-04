package interceptors

import (
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/authentication"
	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	"google.golang.org/grpc"
)

// Authenticate uses the request logger, including the gRPC method, for shared auth.
func Authenticate(parser *jwt.Parser, publicMethods map[string]struct{}, validator session.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, public := publicMethods[info.FullMethod]; public {
			return handler(ctx, req)
		}
		ctx, err := authentication.Authenticate(ctx, parser, validator)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}
