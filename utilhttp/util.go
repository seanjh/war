package utilhttp

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// LogRequest logs basic details about each handled HTTP request.
func LogRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

// RequireMethods restricts HTTP handlers to the specified HTTP methods, and returns
// 405 Not Supported when the handler is invoked with an unsupported method.
func RequireMethods(h http.Handler, methods ...string) http.Handler {
	defaultMethods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodTrace,
		http.MethodPatch,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
	}
	m := map[string]bool{}
	supported := []string{}
	for _, method := range defaultMethods {
		m[method] = false
	}
	for _, method := range methods {
		_, ok := m[method]
		if !ok {
			log.Printf("Ignoring unrecognized method: %s", method)
			continue
		}
		supported = append(supported, method)
		m[method] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isSupported := m[r.Method]
		if isSupported {
			h.ServeHTTP(w, r)
			return
		}
		allowed := strings.Join(supported, ", ")
		w.Header().Add("Allow", allowed)
		http.Error(w, fmt.Sprintf("Unsupported HTTP method: %s", r.Method), http.StatusMethodNotAllowed)
	})
}

func RequireReadOnlyMethods(h http.Handler) http.Handler {
	return RequireMethods(h, http.MethodGet, http.MethodOptions, http.MethodHead)
}
