package auth

import "context"

type contextKey string

const userIDKey contextKey = "userID"

func WithUserId(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func UserIdFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)

	return id, ok
}
