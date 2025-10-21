package dbauthservice

import "context"

type dbAuthServiceBypassKey struct{}

// NewContextWithDbAuthServiceBypass creates a new Context object that instructs
// the authorization policies to bypass authorization checks.
func NewContextWithDbAuthServiceBypass(ctx context.Context) context.Context {
	return context.WithValue(ctx, dbAuthServiceBypassKey{}, dbAuthServiceBypassKey{})
}

// BypassDbAuthServiceFromContext checks whether the authorization checks should
// be bypassed for the given Context.
func BypassDbAuthServiceFromContext(ctx context.Context) bool {
	if value := ctx.Value(dbAuthServiceBypassKey{}); value != nil {
		return true
	}
	return false
}
