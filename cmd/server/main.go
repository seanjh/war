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

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"

	"github.com/seanjh/war/internal/appcontext"
	"github.com/seanjh/war/internal/db"
	"github.com/seanjh/war/internal/game"
	"github.com/seanjh/war/internal/httputil"
)

var portFlag = flag.Int("port", 3000, "Listen port number")
var hostFlag = flag.String("host", "localhost", "Listen hostname")
var dsnFlag = flag.String("dsn", "file::memory:", "SQLite data source name")
var migrateFlag = flag.Bool("migrate", false, "Run the database migrations")

const connParams = "_fk=true&_busy_timeout=5000&_sync=1&_cache_size=1000000000&_journal=WAL&_txlock=immediate"

func main() {
	flag.Parse()

	connString := fmt.Sprintf("%s?%s", *dsnFlag, connParams)
	writeDB, err := sql.Open("sqlite3", fmt.Sprintf("%s&mode=rw", connString))
	if err != nil {
		log.Fatal(err)
	}
	writeDB.SetMaxOpenConns(1)
	readDB, err := sql.Open("sqlite3", fmt.Sprintf("%s&mode=ro", connString))
	if err != nil {
		log.Fatal(err)
	}
	readDB.SetMaxOpenConns(max(4, runtime.NumCPU()))

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	ctx := &appcontext.AppContext{
		Logger: logger,
		DBReader: &appcontext.AppContextDB{
			DB:    readDB,
			Query: db.New(readDB),
		},
		DBWriter: &appcontext.AppContextDB{
			DB:    writeDB,
			Query: db.New(writeDB),
		},
	}
	mux := game.SetupRoutes(httputil.SetupRoutes(http.NewServeMux()))
	wrappedMux := ctx.Middleware(httputil.LogRequestMiddleware(mux, ctx.Logger))

	if *migrateFlag {
		migrateConnString := fmt.Sprintf("sqlite3://%s?%s", *dsnFlag, connParams)
		log.Printf("Connecting to db: %s", migrateConnString)
		m, err := migrate.New("file://./internal/db/migrations", migrateConnString)
		if err != nil {
			log.Fatalf("Failed to connect to database to migrate: %v", err)
		}
		if err = m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to perform migration: %v", err)
		}
	}

	log.Printf("Starting server at %s:%d\n", *hostFlag, *portFlag)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *hostFlag, *portFlag), wrappedMux))
}
