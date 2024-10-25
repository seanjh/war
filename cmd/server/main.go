package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/seanjh/war/internal/appcontext"
	"github.com/seanjh/war/internal/game"
	"github.com/seanjh/war/internal/httputil"
)

var portFlag = flag.Int("port", 3000, "Listen port number")
var hostFlag = flag.String("host", "localhost", "Listen hostname")

func ping(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("pong!\n"))
	w.Header().Set("Conent-Type", "text/plain")
	if err != nil {
		ctx := appcontext.GetAppContext(r)
		ctx.Logger.Error("Failed to pong",
			"err", err,
		)
	}
}

func main() {
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	ctx := &appcontext.AppContext{
		Logger: logger,
		DB:     db,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("GET /ping", ping)
	game.SetupRoutes(mux)
	wrappedMux := ctx.Middleware(httputil.LogRequestMiddleware(mux, ctx.Logger))

	log.Printf("Starting server at %s:%d\n", *hostFlag, *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), wrappedMux))
}
