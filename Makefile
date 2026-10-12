# Universal build/run helper (Windows / macOS / Linux).
# Requires: Go >= 1.22, Node >= 18.
#
#   make run          - backend (go run) + frontend (vite) via npm
#   make build        - frontend dist + backend binary for current OS
#   make build-linux  - binary for Linux server (from any OS)
#   make build-win    - binary for Windows (.exe, from any OS)
#   make build-mac    - binary for macOS (from any OS)
#   make start        - run local backend binary (builds if missing)
#   make clean        - remove local build artifacts
#   make migrate:fresh / make migrate-fresh - reset database & run seed
#   make db:seed / make db-seed            - run initial seeders

BACKEND_DIR := backend
BIN := $(BACKEND_DIR)/pos-backend
GO ?= $(shell if command -v go >/dev/null 2>&1; then echo go; elif command -v go.exe >/dev/null 2>&1; then echo go.exe; else echo go; fi)

.PHONY: run build build-linux build-win build-mac start clean migrate-fresh migrate\:fresh db-seed db\:seed dev-backend dev-backend-windows dev-backend-wsl

run:
	npm run dev

build:
	npm --prefix frontend run build
	cd $(BACKEND_DIR) && $(GO) build -o pos-backend ./cmd/api

build-linux:
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 $(GO) build -o pos-backend-linux ./cmd/api

build-win:
	cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 $(GO) build -o pos-backend-win.exe ./cmd/api

build-mac:
	cd $(BACKEND_DIR) && GOOS=darwin GOARCH=arm64 $(GO) build -o pos-backend-mac ./cmd/api

start:
	npm start

dev-backend:
	node scripts/dev-backend.mjs

dev-backend-windows:
	node scripts/dev-backend.mjs --target=windows

dev-backend-wsl:
	node scripts/dev-backend.mjs --target=wsl

clean:
	rm -f $(BIN) $(BACKEND_DIR)/pos-backend-linux $(BACKEND_DIR)/pos-backend-win.exe $(BACKEND_DIR)/pos-backend-mac

# Database Migration
migrate-fresh:
	cd $(BACKEND_DIR) && $(GO) run ./cmd/migrate/main.go

migrate\:fresh: migrate-fresh

db-seed:
	cd $(BACKEND_DIR) && $(GO) run ./cmd/seed/main.go

db\:seed: db-seed
