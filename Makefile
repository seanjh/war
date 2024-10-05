.ONESHELL:

check:
	go fmt
	go mod tidy
	go mod verify
.PHONY: check

start:
	make -j 2 start-server start-tailwinds
.PHONY: start

build:
	go build -o ./tmp/server ./cmd/server.go
.PHONY: build

start-server:
	air -c .air.toml
.PHONY: start-server

start-tailwinds:
	pnpm exec tailwindcss -i ./styles/main.css -o ./assets/main.css --watch
