package appcontext

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/seanjh/war/internal/db"
)

type AppContextDB struct {
	Query *db.Queries
	DB    *sql.DB
}

type AppContext struct {
	Logger   *slog.Logger
	DBReader *AppContextDB
	DBWriter *AppContextDB
}

type key string

const appContextKey key = "appcontext"

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
