package httputil

import (
	"log/slog"
	"net/http"
	"time"
)

// LogRequest logs basic details about each handled HTTP request.
func LogRequestMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		slog.Info("Handled request",
			"method", r.Method,
			"path", r.URL.Path,
			"elapsed", elapsed,
		)
	})
}
