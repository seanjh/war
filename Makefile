.ONESHELL:

check:
	go fmt
	go mod tidy
	go mod verify
.PHONY: check

test:
	go test ./...
.PHONY: test

start:
	make -j 2 start-server start-tailwinds
.PHONY: start

build:
	go build -o ./bin/server ./cmd/server/main.go
.PHONY: build

sql-generate:
	sqlc generate
.PHONY: generate-sql

sql-migrate:
	migrate -verbose -database=sqlite3://./tmp/war.db -source=file://./internal/db/migrations up
.PHONY: sql-migrate

sql-migrate-drop:
	migrate -verbose -database=sqlite3://./tmp/war.db -source=file://./internal/db/migrations drop
.PHONY: sql-migrate

start-server:
	air -c .air.toml
.PHONY: start-server

start-tailwinds:
	pnpm exec tailwindcss -i ./styles/main.css -o ./assets/main.css --watch
