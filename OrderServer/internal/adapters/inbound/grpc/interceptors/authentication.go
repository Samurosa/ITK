package interceptors

import (
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/authentication"
	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	"google.golang.org/grpc"
)

func Authentication(parser *jwt.Parser, validator session.Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return authentication.AuthInterceptor(parser, nil, validator)(ctx, req, info, handler)
	}
}

func StreamAuthentication(parser *jwt.Parser, validator session.Validator) grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return authentication.AuthStreamInterceptor(parser, validator)(srv, stream, info, handler)
	}
}
