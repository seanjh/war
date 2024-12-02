package httputil

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/seanjh/war/internal/appcontext"
)

// LogRequest logs basic details about each handled HTTP request.
func LogRequestMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		log.Info("Handled request",
			"method", r.Method,
			"path", r.URL.Path,
			"elapsed", fmt.Sprintf("%s", elapsed),
		)
	})
}

func Ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong!\n"))
	w.Header().Set("Conent-Type", "text/plain")
	if err != nil {
		ctx := appcontext.GetAppContext(r)
		ctx.Logger.Error("Failed to pong",
			"err", err,
		)
	}
}

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("GET /ping", Ping)
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	return mux
}
