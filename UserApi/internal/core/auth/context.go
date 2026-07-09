package auth

import "context"

func GetUserIDByContext(ctx context.Context) (string, error) {
	value, ok := ctx.Value(UserIDContextKey).(string)

	if !ok || value == "" {
		return "", ErrInvalidContext
	}
	return value, nil
}
