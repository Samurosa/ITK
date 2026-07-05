package interceptors

import (
	"ITK_Code/m/v2/internal/core/auth"
	"context"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]struct{}{
	"/user.UserService/Login":        {},
	"/user.UserService/Registration": {},
	"/user.UserService/RefreshToken": {},
}

func AuthInterceptor(
	log *zap.Logger,
	tokenManager auth.TokenManager,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := log.Named("auth interceptor")

		if _, ok := publicMethods[info.FullMethod]; ok {
			log.Info("A verification token is not required for this RPS.")
			return handler(ctx, req)
		}

		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Error("missing metadata")
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		values := meta.Get("authorization")
		if len(values) == 0 {
			log.Error("missing token")
			return nil, status.Errorf(codes.Unauthenticated, "missing token")
		}

		token := strings.TrimPrefix(values[0], "Bearer ")

		claims, err := tokenManager.ParseAccessToken(token)
		if err != nil {
			log.Error("invalid token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, auth.UserIDContextKey, claims.UserID)
		ctx = context.WithValue(ctx, auth.RoleContextKey, claims.Role)

		return handler(ctx, req)
	}
}
