package domain

import "context"

type auditUserIDKey struct{}

func WithAuditUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, auditUserIDKey{}, userID)
}

func AuditUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(auditUserIDKey{}).(string)
	return userID, ok && userID != ""
}
