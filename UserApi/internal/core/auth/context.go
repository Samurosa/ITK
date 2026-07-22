package auth

import "context"

func GetJTIFromContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(jtiContextKey).(string)
	if !ok || value == "" {
		return "", ErrInvalidContext
	}

	return value, nil
}

func WithJTI(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, jtiContextKey, ip)
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(userIDContextKey).(string)
	if !ok || value == "" {
		return "", ErrInvalidContext
	}

	return value, nil
}

func WithUserID(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, userIDContextKey, ip)
}

func GetClientIPFromContext(ctx context.Context) (string, error) {
	ip, ok := ctx.Value(clientIPKey).(string)

	if !ok || ip == "" {
		return "", ErrInvalidContext
	}

	return ip, nil
}

func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

func WithRole(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, roleContextKey, ip)
}

func WithDeviceID(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, deviceContextKey, ip)
}
