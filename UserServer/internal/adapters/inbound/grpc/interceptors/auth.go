package interceptors

import (
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/authentication"
	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	"google.golang.org/grpc"
)

// AuthInterceptor supplies the request logger to the shared authentication API.
func AuthInterceptor(parser *jwt.Parser, publicMethods map[string]struct{}, validator session.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return authentication.AuthInterceptor(parser, publicMethods, validator)(ctx, req, info, handler)
	}
}
