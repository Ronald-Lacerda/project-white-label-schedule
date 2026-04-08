package shared

import "context"

type contextKey string

const (
	tenantKey contextKey = "establishment_id"
	userKey   contextKey = "user_id"
)

func WithEstablishmentID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tenantKey, id)
}

func EstablishmentIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(tenantKey).(string)
	return v
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userKey, id)
}

func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userKey).(string)
	return v
}
