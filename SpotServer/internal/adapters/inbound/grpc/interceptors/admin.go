package interceptors

import (
	"ITK_Code/m/v2/internal/core/dto"
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RequireAdmin(publicMethods map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, public := publicMethods[info.FullMethod]; public {
			return handler(ctx, req)
		}

		role, authenticated := sharedContext.Role(ctx)

		if !authenticated {
			return nil, status.Error(codes.Unauthenticated, "authentication required")
		}

		if role != string(dto.AdminRole) {
			return nil, status.Error(codes.PermissionDenied, "administrator role required")
		}

		return handler(ctx, req)
	}
}
