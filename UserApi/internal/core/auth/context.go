package auth

import "context"

func GetJTIFromContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(JTIContextKey).(string)
	if !ok || value == "" {
		return "", ErrInvalidContext
	}

	return value, nil
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(UserIDContextKey).(string)
	if !ok || value == "" {
		return "", ErrInvalidContext
	}

	return value, nil
}
