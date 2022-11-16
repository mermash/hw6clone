package main

import (
	"context"
)

type Session struct {
	ID     string
	UserID string
}

type ctxKey int

const sessionKey ctxKey = 1

func SessionFromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(sessionKey).(*Session)
	if !ok {
		return nil, ErrNoAuth
	}
	return sess, nil
}
