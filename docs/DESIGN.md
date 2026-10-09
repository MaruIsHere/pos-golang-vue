# System Design Document
## POS Kasir Pro

### 1. Architecture
- **Frontend:** Vue 3 (Composition API), Vite, Tailwind CSS v4, Vue Router, Pinia.
- **Backend:** Go 1.22+, Gin framework, GORM.
- **Database:** SQLite (dev/standalone via glebarez/sqlite) and MySQL (production).
- **AI/OCR:** `ocrs-cli` (Rust-based Tensor Engine for local ONNX model execution).

### 2. Project Structure
```text
pos-golang-vue/
├── backend/          # Go source code
│   ├── cmd/api/      # Main entry point
│   ├── internal/     # Handlers, models, routes, config
│   └── docs/         # OpenAPI specifications
├── frontend/         # Vue 3 source code
│   ├── src/          # Components, views, assets
│   └── dist/         # Compiled static files
├── bin/              # Compiled external tools (ocrs, ocrs.exe)
└── docs/             # Product and Design documentation
```

### 3. Deployment Strategy
- **Universal Binary:** The backend serves the frontend static files (`../frontend/dist`) directly via Gin `router.Static()`.
- **Cross-Compilation:** The Go backend and Rust OCR engine can be cross-compiled to run seamlessly on Linux or Windows.
- **Local Dev:** `npm run dev` proxies `/api` to the Go backend on port 8080.
