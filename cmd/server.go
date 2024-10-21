package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/seanjh/war/context"
	"github.com/seanjh/war/game"
	u "github.com/seanjh/war/utilhttp"
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

func main() {
	flag.Parse()

	ctx := &context.AppContext{}

	fs := http.FileServer(http.Dir("./assets"))
	// TODO(sean) wrap this handler
	http.Handle("GET /assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("GET /ping", u.LogRequest(ping))
	game.SetupHandlers(ctx)

	log.Printf("Starting server at %s:%d\n", *hostFlag, *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), nil))
}
