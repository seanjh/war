package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime"

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

const connectionUrl = "file::memory:?cache=shared&_fk=true&_busy_timeout=5000&_sync=1&_cache_size=1000000000&_journal=WAL&_txlock=immediate"

func main() {
	flag.Parse()

	readDB, err := sql.Open("sqlite3", connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	readDB.SetMaxOpenConns(max(4, runtime.NumCPU()))

	writeDB, err := sql.Open("sqlite3", connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	writeDB.SetMaxOpenConns(1)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	ctx := &appcontext.AppContext{
		Logger:  logger,
		ReadDB:  readDB,
		WriteDB: writeDB,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("GET /ping", ping)
	game.SetupRoutes(mux)
	wrappedMux := ctx.Middleware(httputil.LogRequestMiddleware(mux, ctx.Logger))

	log.Printf("Starting server at %s:%d\n", *hostFlag, *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), wrappedMux))
}
