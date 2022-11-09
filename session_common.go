package main

import (
	"context"
	"net/http"
)

type Session struct {
	ID     string
	UserID string
}

type SessionManager interface {
	Check(*http.Request) (*Session, error)
	Create(http.ResponseWriter, *User) (*Session, error)
	DestroyCurrent(http.ResponseWriter, *http.Request) error
	DestroyAll(http.ResponseWriter, *User) error
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
