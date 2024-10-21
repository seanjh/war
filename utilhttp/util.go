package utilhttp

import (
	"log"
	"net/http"
)

// LogRequest logs basic details about each handled HTTP request.
func LogRequest(handler http.HandlerFunc) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		handler(w, r)
	})
}
