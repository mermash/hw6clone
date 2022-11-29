package main

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type AccessLogger struct {
	Logger *zap.Logger
}

func NewAccessLoggerMiddleware(log *zap.Logger) *AccessLogger {
	return &AccessLogger{
		Logger: log,
	}
}

func (ac *AccessLogger) AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		ac.Logger.Info(r.URL.Path,
			zap.String("requestID", RequestIDFromContext(r.Context())),
			zap.String("method", r.Method),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("url", r.URL.Path),
			zap.Duration("work_time", time.Since(start)),
		)
	})
}
