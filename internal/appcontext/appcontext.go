package appcontext

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
)

type AppContext struct {
	Logger  *slog.Logger
	ReadDB  *sql.DB
	WriteDB *sql.DB
}

type key int

const appContextKey key = 0

// Middleware adds the application context to the request context.
func (c *AppContext) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), appContextKey, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAppContext returns the application context for the request.
func GetAppContext(r *http.Request) *AppContext {
	ctx, ok := r.Context().Value(appContextKey).(*AppContext)
	if !ok {
		panic("Failed to load app context from request")
	}
	return ctx
}
