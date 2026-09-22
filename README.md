# POS Golang + Vue — Sistem Kasir Modern (PWA)

Aplikasi Point of Sales (POS) kasir untuk UMKM / coffee shop / resto kecil. Backend REST API memakai **Go + Gin + GORM**, frontend kasir memakai **Vue 3 + Vite + PWA**. Database default **SQLite** (`pos.db`), bisa pindah live ke **MySQL** dari halaman Settings tanpa restart manual.

> Repo ini: `backend/` (Go module `pos-backend`) + `frontend/` (Vue 3). Backend sekaligus menyajikan file static `frontend/dist` saat production.

## Fitur yang sudah jalan

* **Register / Kasir** (`RegisterView.vue`): katalog produk + filter kategori + search/barcode, keranjang (qty, notes, validasi stok), diskon, pajak otomatis dari `store_setting.tax_percentage`, checkout multi-metode (`cash`, `qris`, `transfer`, `debit`), struk via `ReceiptModal.vue`.
* **Produk & Kategori** (`ProductsView.vue`): CRUD produk (`price`, `cost_price`, `stock`, `barcode`, `image_url`, `is_active`) + CRUD kategori. Search backend via `?search=` (nama/barcode), filter `?category_id=`.
* **Inventory** (`InventoryView.vue`): stok masuk (`in`) / keluar (`out`) dengan alasan (`pembelian_supplier`, `barang_rusak`, `barang_hilang`, `expired`, `promosi`, `retur_penjualan`), riwayat `stock-movements` + update stok atomik (transaction).
* **Orders** (`OrdersView.vue`): riwayat transaksi (`?limit=`, default `20`), detail per invoice, refund (status jadi `refunded` + stok dikembalikan + catat movement `in`).
* **Customers** (`CustomersView.vue` + `PaymentModal.vue`): CRUD customer (`name`, `phone`, `email`, `address`, `points`), pilih customer saat bayar.
* **Reports / Dashboard** (`ReportsView.vue`): `total_orders`, `total_revenue`, `total_items_sold`, top 5 produk terlaris, 5 transaksi terakhir.
* **Settings & Voucher** (`SettingsView.vue`): profil toko + QRIS (upload max 5MB jadi base64 `qris_image_url`), CRUD voucher (`percent`/`flat`), switch database `sqlite <-> mysql` live.
* **PWA**: installable (`vite-plugin-pwa`), `manifest` + `workbox`, banner `PwaInstallBanner.vue`.
* **Seed otomatis**: saat DB kosong, diisi 4 kategori, 9 produk, 3 customer, 5 voucher, 1 setting toko.

## Struktur repo

```text
pos-golang-vue/
├── backend/
│   ├── main.go              # entry point: load config, InitDB, router Gin, serve ../frontend/dist
│   ├── config/
│   │   └── config.go        # LoadConfig/SaveConfig (config.json)
│   ├── config.json          # port, db_engine, sqlite_path, mysql_dsn
│   ├── database/
│   │   └── database.go      # InitDB (sqlite/mysql), AutoMigrate 8 tabel, SeedInitialData
│   ├── models/
│   │   └── models.go        # Category, Product, Order, OrderItem, StoreSetting, Voucher, Customer, StockMovement
│   ├── handlers/
│   │   └── handlers.go      # 25 handlers REST (lihat docs/openapi.yaml)
│   ├── go.mod               # module pos-backend, go 1.27.1
│   └── pos.db               # SQLite default (di-ignore git, jangan commit)
├── frontend/
│   ├── src/
│   │   ├── App.vue          # tab navigation + fetch /api/settings
│   │   ├── main.js          # createApp
│   │   ├── views/           # Register, Products, Inventory, Customers, Orders, Reports, Settings (7 file)
│   │   ├── components/      # Navbar, ProductCard, CartDrawer, PaymentModal, ReceiptModal, PwaInstallBanner
│   │   └── assets/styles.css
│   ├── vite.config.js       # dev port 3000, proxy /api -> http://localhost:8080
│   └── package.json         # vue ^3.4, vite ^5.1, vite-plugin-pwa ^1.3
├── docs/
│   └── openapi.yaml         # spesifikasi API (sumber kebenaran endpoint)
├── package.json             # script root concurrently backend+frontend
└── README.md
```

## Tech stack (versi aktual)

| Layer | Stack |
|---|---|
| Backend | Go `1.27.1`, `gin-gonic/gin v1.12.0`, `gorm v1.31.2`, `glebarez/sqlite v1.11.0` (pure-Go), `gorm/driver/mysql v1.6.0`, `go-sql-driver/mysql v1.8.1`, `gin-contrib/cors v1.7.8` |
| Frontend | `vue ^3.4.0`, `vite ^5.1.0`, `vite-plugin-pwa ^1.3.0`, `lucide-vue-next ^0.344.0` |
| DB | SQLite default, MySQL 8 opsional |

## Prasyarat

* Go `>= 1.27` (`go version`)
* Node.js `>= 18` + npm
* MySQL opsional (hanya jika `db_engine=mysql`)

## Quickstart (development)

```bash
# 1. Backend (dari folder backend, port 8080)
cd backend
go mod tidy
go run .
# cek: curl http://localhost:8080/api/health

# 2. Frontend (terminal baru, port 3000, proxy /api -> 8080)
cd frontend
npm install
npm run dev
# buka http://localhost:3000
```

Atau dari root sekaligus (Windows / macOS / Linux — perlu koneksi internet pertama kali untuk `concurrently`):

```bash
npm run dev
```

> Semua script root universal: `dev:backend` memakai `go run .`, `start` memakai `node scripts/start-backend.mjs` (otomatis pilih `pos-backend` / `pos-backend.exe` + build jika belum ada). Jangan commit binary / `*.db` (sudah di `.gitignore`).

## Build production (universal Linux / Windows / macOS)

```bash
# Cara singkat (semua OS):
make build     # frontend dist + backend binary OS lokal
make start     # = npm start, jalankan binary lokal (build otomatis jika belum ada)

# Cross-compile tanpa pindah OS (butuh waktu, SQLite pure-Go dikompilasi per target):
make build-linux   # -> backend/pos-backend-linux (untuk server)
make build-win     # -> backend/pos-backend-win.exe (untuk Windows)
make build-mac     # -> backend/pos-backend-mac (untuk macOS arm64)
```

Manual (setara `make`, dijalankan dari folder backend agar `config.json` + serve dist benar):

```bash
cd backend
go build -o pos-backend . && ./pos-backend
# buka http://localhost:8080 (API + frontend static jadi satu)
```

Alasan cross-compile works: SQLite driver yang dipakai (`glebarez/sqlite` / `modernc.org/sqlite`) pure-Go, tidak butuh CGO/`gcc`.

## Konfigurasi (`backend/config.json`)

```json
{
  "port": "8080",
  "db_engine": "sqlite",
  "sqlite_path": "pos.db",
  "mysql_dsn": "root:@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local"
}
```

* `POST /api/settings/switch-db` akan update file ini + row `store_settings` otomatis.
* Backend harus dijalankan dari folder `backend/` karena path `config.json` dan `../frontend/dist` relatif terhadap working directory (`main.go:92-94`, `config/config.go:18`).

## Database

* `InitDB(engine, mysqlDsn, sqlitePath)` di `database/database.go:41`.
* `AutoMigrate`: `Category`, `Product`, `Order`, `OrderItem`, `StoreSetting`, `Voucher`, `Customer`, `StockMovement`.
* MySQL: database dibuat otomatis (`CREATE DATABASE IF NOT EXISTS`, `database.go:20`).
* Seed hanya jika tabel kosong: 4 kategori (`Makanan Utama`, `Minuman`, `Cemilan & Snack`, `Dessert`), 9 produk (contoh `Nasi Goreng Special Egg 28000`, `Kopi Susu Gula Aren 18000`), 3 customer (Budi/Siti/Dewi), 5 voucher (`DISKON10`, `DISKON20`, `HEMAT10K`, `HEMAT50K`, `POSHEMAT`), 1 setting (`KASIR COFFEE & BISTRO`).

## API reference (ringkas)

Base URL dev: `http://localhost:8080/api`. Spec lengkap dan contoh payload: [`docs/openapi.yaml`](docs/openapi.yaml).

| Method & Path | Handler | Keterangan |
|---|---|---|
| `GET /health` | inline `main.go:41` | `{status, system, db_engine}` |
| `GET /categories` | `GetCategories` | list kategori |
| `POST /categories` | `CreateCategory` | body `{name, icon}` |
| `DELETE /categories/:id` | `DeleteCategory` | hapus kategori |
| `GET /products?category_id=&search=` | `GetProducts` | search nama/barcode, preload category |
| `POST /products` | `CreateProduct` | buat produk, return + category |
| `PUT /products/:id` | `UpdateProduct` | 404 jika tidak ada |
| `DELETE /products/:id` | `DeleteProduct` | hapus produk |
| `POST /orders` | `CreateOrder` | checkout, cek stok, transaction, hitung `grandTotal=(total-discount)+tax` |
| `GET /orders?limit=20` | `GetOrders` | order desc, preload items |
| `GET /orders/:id` | `GetOrderById` | detail + items |
| `POST /orders/:id/refund` | `RefundOrder` | 400 jika sudah `refunded`, kembalikan stok |
| `GET /customers` | `GetCustomers` | order desc |
| `POST /customers` | `CreateCustomer` | `name` wajib |
| `PUT /customers/:id` | `UpdateCustomer` | 404 jika tidak ada |
| `DELETE /customers/:id` | `DeleteCustomer` | hapus customer |
| `GET /stock-movements` | `GetStockMovements` | preload product, order desc |
| `POST /stock-movements` | `CreateStockMovement` | `{product_id, type:in/out, quantity>0, reason, notes}`, update stok atomik |
| `GET /reports/dashboard` | `GetDashboardStats` | omset, total order, item terjual, top 5, recent 5 |
| `GET /settings` | `GetSettings` | 1 row pertama |
| `PUT /settings` | `UpdateSettings` | update profil toko + QRIS |
| `POST /settings/switch-db` | `SwitchDatabase` | `{engine: sqlite/mysql, mysql_dsn?}`, re-`InitDB` |
| `GET /vouchers` | `GetVouchers` | order desc |
| `POST /vouchers` | `CreateVoucher` | `{code wajib unik, type, value, description}`, auto `is_active=true` |
| `DELETE /vouchers/:id` | `DeleteVoucher` | hapus voucher |

Contoh cepat (teruji terhadap handler):

```bash
curl http://localhost:8080/api/health
curl "http://localhost:8080/api/products?search=kopi"
curl -X POST http://localhost:8080/api/orders \
  -H 'Content-Type: application/json' \
  -d '{"customer_name":"Umum","payment_method":"cash","paid_amount":50000,"discount":0,"tax":2800,"items":[{"product_id":4,"quantity":2}]}'
```

## Frontend pages -> API yang dipakai (terverifikasi)

* `RegisterView`: `GET /categories`, `GET /products`, `POST /orders`
* `ProductsView`: `GET /products`, `GET /categories`, `POST /categories`, `DELETE /products/:id`
* `InventoryView`: `GET /products`, `GET /stock-movements`, `POST /stock-movements`
* `CustomersView`: `GET /customers`, `DELETE /customers/:id`
* `OrdersView`: `GET /orders?limit=50`, `POST /orders/:id/refund`
* `ReportsView`: `GET /reports/dashboard`
* `SettingsView` + `PaymentModal` + `CartDrawer`: `GET/PUT /settings`, `GET/POST/DELETE /vouchers`, `POST /settings/switch-db`, `GET /customers`

## Troubleshooting

* `8080 already in use`: ganti `port` di `backend/config.json`.
* Frontend `fetch /api/...` 404 saat dev: pastikan Vite proxy jalan (`vite.config.js:47-50` proxy ke `8080`) dan backend `go run .` aktif.
* Production blank page: pastikan `frontend/dist` sudah di-build dan backend dijalankan dari folder `backend/` (karena `router.Static("/assets", "../frontend/dist/assets")`).
* `pos.db` terkunci: jangan jalankan 2 backend bersamaan di file DB yang sama.
* MySQL `connection refused`: cek DSN, user, dan MySQL jalan di `127.0.0.1:3306`.

## Coverage dokumentasi

* 25/25 endpoint REST terdokumentasi di `docs/openapi.yaml` (health 1, categories 3, products 4, orders 4, customers 4, stock 2, dashboard 1, settings 3, vouchers 3).
* 8/8 model GORM terdokumentasi sebagai schema.
* Contoh `curl` + payload checkout / stock-movement / switch-db diambil dari `handlers.go` dan pemakaian nyata di `RegisterView.vue` / `InventoryView.vue` / `SettingsView.vue`.
