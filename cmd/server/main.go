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
	"github.com/seanjh/war/internal/db"
	"github.com/seanjh/war/internal/game"
	"github.com/seanjh/war/internal/httputil"
)

var portFlag = flag.Int("port", 3000, "Listen port number")
var hostFlag = flag.String("host", "localhost", "Listen hostname")
var sqliteConnectionUrl = flag.String("sqliteUrl", "file::memory:", "SQLite connection URL")

const connectionUrl = "file::memory:?cache=shared&_fk=true&_busy_timeout=5000&_sync=1&_cache_size=1000000000&_journal=WAL&_txlock=immediate"

func main() {
	flag.Parse()

	readDB, err := sql.Open("sqlite3", fmt.Sprintf("%s&mode=ro", connectionUrl))
	if err != nil {
		log.Fatal(err)
	}
	readDB.SetMaxOpenConns(max(4, runtime.NumCPU()))

	writeDB, err := sql.Open("sqlite3", fmt.Sprintf("%s&mode=rw", connectionUrl))
	if err != nil {
		log.Fatal(err)
	}
	writeDB.SetMaxOpenConns(1)

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	ctx := &appcontext.AppContext{
		Logger:     logger,
		ReadQuery:  db.New(readDB),
		WriteQuery: db.New(writeDB),
	}
	mux := game.SetupRoutes(httputil.SetupRoutes(http.NewServeMux()))
	wrappedMux := ctx.Middleware(httputil.LogRequestMiddleware(mux, ctx.Logger))

	log.Printf("Starting server at %s:%d\n", *hostFlag, *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), wrappedMux))
}
