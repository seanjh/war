package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/seanjh/war/utilhttp"
)

var portFlag = flag.Int("port", 3000, "Listen port number")
var hostFlag = flag.String("host", "localhost", "Listen hostname")

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong!\n"))
	w.Header().Set("Conent-Type", "text/plain")
	if err != nil {
		log.Printf("Failed to pong: %v", err)
	}
}

func renderGame() http.Handler {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "game.html"),
	))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Rendering game template")
		tmpl.ExecuteTemplate(w, "layout", nil)
	})
}

func main() {
	flag.Parse()

	http.Handle("/", utilhttp.RequireReadOnlyMethods(utilhttp.LogRequest(renderGame())))

	fs := utilhttp.RequireReadOnlyMethods(utilhttp.LogRequest(http.FileServer(http.Dir("./assets"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.Handle("/ping", utilhttp.LogRequest(utilhttp.RequireReadOnlyMethods(http.HandlerFunc(ping))))

	log.Printf("Starting server at %s:%d\n", *hostFlag, *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), nil))
}
