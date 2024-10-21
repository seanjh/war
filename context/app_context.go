package context

import (
	"log/slog"
	"net/http"

	"github.com/seanjh/war/storage"
)

type AppContext struct {
	Logger  *slog.Logger
	Storage *storage.Storage
}

type HandlerFuncWithContext = func(http.ResponseWriter, *http.Request, *AppContext)

// WithContext wraps the provided HTTP handler function to provide AppContext.
func (c *AppContext) WithContext(handler HandlerFuncWithContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}
