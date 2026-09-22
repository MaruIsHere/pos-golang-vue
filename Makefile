# Universal build/run helper (Windows / macOS / Linux).
# Requires: Go >= 1.27, Node >= 18.
#
#   make run          - backend (go run) + frontend (vite) via npm
#   make build        - frontend dist + backend binary for current OS
#   make build-linux  - binary for Linux server (from any OS)
#   make build-win    - binary for Windows (.exe, from any OS)
#   make build-mac    - binary for macOS (from any OS)
#   make start        - run local backend binary (builds if missing)
#   make clean        - remove local build artifacts

BACKEND_DIR := backend
BIN := $(BACKEND_DIR)/pos-backend

.PHONY: run build build-linux build-win build-mac start clean

run:
	npm run dev

build:
	npm --prefix frontend run build
	cd $(BACKEND_DIR) && go build -o pos-backend .

build-linux:
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o pos-backend-linux .

build-win:
	cd $(BACKEND_DIR) && GOOS=windows GOARCH=amd64 go build -o pos-backend-win.exe .

build-mac:
	cd $(BACKEND_DIR) && GOOS=darwin GOARCH=arm64 go build -o pos-backend-mac .

start:
	npm start

clean:
	rm -f $(BIN) $(BACKEND_DIR)/pos-backend-linux $(BACKEND_DIR)/pos-backend-win.exe $(BACKEND_DIR)/pos-backend-mac
