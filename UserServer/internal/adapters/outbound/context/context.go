package context

import (
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"context"
)

func GetJTIFromContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(dto.JtiContextKey).(string)
	if !ok || value == "" {
		return "", errors.ErrInvalidContext
	}

	return value, nil
}

func WithJTI(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, dto.JtiContextKey, ip)
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(dto.UserIDContextKey).(string)
	if !ok || value == "" {
		return "", errors.ErrInvalidContext
	}

	return value, nil
}

func WithUserID(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, dto.UserIDContextKey, ip)
}

func GetClientIPFromContext(ctx context.Context) (string, error) {
	ip, ok := ctx.Value(dto.ClientIPKey).(string)

	if !ok || ip == "" {
		return "", errors.ErrInvalidContext
	}

	return ip, nil
}

func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, dto.ClientIPKey, ip)
}

func WithRole(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, dto.RoleContextKey, ip)
}

func WithDeviceID(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, dto.DeviceContextKey, ip)
}
