# Backend Makefile
# Supports Linux, macOS, WSL (with either native go or Windows go.exe), and Windows.

GO ?= $(shell if command -v go >/dev/null 2>&1; then echo go; elif command -v go.exe >/dev/null 2>&1; then echo go.exe; else echo go; fi)

.PHONY: run dev build test clean migrate-fresh migrate\:fresh db-seed db\:seed

run:
	$(GO) run ./cmd/api

dev: run

build:
	$(GO) build -o pos-backend ./cmd/api

test:
	$(GO) test ./...

migrate-fresh:
	$(GO) run ./cmd/migrate/main.go

migrate\:fresh: migrate-fresh

db-seed:
	$(GO) run ./cmd/seed/main.go

db\:seed: db-seed

clean:
	rm -f pos-backend pos-backend.exe pos-backend-linux pos-backend-mac pos-backend-win.exe
