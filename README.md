# POS Kasir Pro

POS Kasir Pro adalah aplikasi kasir berbasis web untuk mengelola penjualan dan operasional toko. Repositori ini berisi frontend Vue 3/TypeScript, backend Go/Gin, database SQLite atau MySQL, serta utilitas OCR lokal untuk membantu membaca nominal pada bukti pembayaran.

> PWA menyimpan aset aplikasi untuk membantu membuka shell UI. Proses transaksi tetap memerlukan backend aktif; aplikasi belum menyediakan transaksi offline dan sinkronisasi otomatis.

## Fitur

- POS dengan katalog, keranjang, pajak/diskon, pembayaran tunai/QRIS/transfer, dan struk.
- Pencarian barcode dan scanner kamera pada perangkat yang mendukung.
- Produk master dan katalog produk per toko, kategori, merk, dan tipe.
- Penerimaan/pengeluaran stok dan riwayat mutasi.
- Riwayat pesanan, refund sesuai hak akses, pelanggan, voucher, dan pengelolaan toko/staf.
- Dashboard penjualan dan ekspor Excel, CSV, serta PDF.
- Role-based access, tema terang/gelap, dan PWA.

## Teknologi

- **Frontend:** Vue 3, TypeScript, Vite, Vue Router, Pinia, Tailwind CSS.
- **Backend:** Go, Gin, GORM.
- **Database:** SQLite default; MySQL dapat dikonfigurasi.
- **OCR:** executable lokal `bin/ocrs` / `bin/ocrs.exe` jika tersedia.

## Menjalankan untuk development

Prasyarat: Go yang memenuhi `backend/go.mod`, Node.js/npm yang mendukung Vite, dan Git.

```bash
git clone <url-repository>
cd pos-golang-vue
npm install
npm --prefix frontend install
npm run dev
```

- Frontend Vite: `http://localhost:3000`
- API Go: `http://localhost:8080`
- Proxy `/api` dan `/uploads` diarahkan ke `127.0.0.1:8080` secara default.
- Untuk mengganti target proxy: set `VITE_API_PROXY_TARGET` sebelum menjalankan Vite.

Server pertama kali dapat membuat `backend/config.local.json` dan database SQLite sesuai konfigurasi. Jangan commit file config lokal atau gunakan kredensial seed di server publik.

## Build dan menjalankan server lokal

```bash
npm run build
npm start
```

Build menjalankan typecheck dan Vite, lalu membangun backend Go. Backend menyajikan aset dari `frontend/dist` dan default berjalan pada port 8080. Aset frontend **tidak di-embed** ke binary backend.

Perintah lintas platform tersedia melalui Makefile:

```bash
make run
make build
make build-linux
make build-win
make build-mac
```

## Akun demo lokal

Seed backend membuat akun `owner`, `admin`, `kepalakasir`, dan `kasir` dengan password awal masing-masing `owner123`, `admin123`, `kepala123`, dan `kasir123` jika akun tersebut belum ada. Akun ini untuk development/demo; ganti password atau hapus sebelum server dibuka ke jaringan.

## Konfigurasi backend

Default ada di `backend/config.json`. Konfigurasi lokal disimpan ke `backend/config.local.json` (diabaikan Git). Environment variable yang didukung:

| Variable | Fungsi |
|---|---|
| `POS_PORT` | Port HTTP backend |
| `POS_DB_ENGINE` | `sqlite` atau `mysql` |
| `POS_SQLITE_PATH` | Path database SQLite |
| `POS_MYSQL_DSN` | DSN koneksi MySQL |
| `VITE_API_PROXY_TARGET` | Target proxy API/upload Vite |

Backup database sebelum migrasi atau perubahan engine. Untuk deployment jaringan publik, gunakan HTTPS, ganti akun default, dan batasi konfigurasi CORS sesuai domain yang dipercaya.

## Struktur kode

```text
frontend/src/views/       # Halaman aplikasi
frontend/src/components/  # Komponen UI
frontend/src/stores/      # State autentikasi, toko, pengaturan
backend/cmd/api/          # Entrypoint backend
backend/internal/         # Route, handler, model, middleware, config, database
scripts/                  # Runner development/start lintas platform
docs/                     # PRD, desain, system, panduan, API, dan audit
```

## Pemeriksaan

```bash
npm --prefix frontend run typecheck
npm --prefix frontend run build
cd backend && go test ./...
```

## Dokumentasi

Mulai dari [indeks dokumentasi](docs/README.md).

- [PRD — kebutuhan produk](docs/product/PRD.md)
- [Desain produk dan UI](docs/design/DESIGN.md)
- [System overview untuk belajar arsitektur](docs/system/SYSTEM_OVERVIEW.md)
- [Roadmap arsitektur SaaS](docs/system/SAAS_ROADMAP.md)
- [Panduan penggunaan lengkap](docs/guides/USER_GUIDE.md)
- [Spesifikasi API OpenAPI](docs/api/openapi.yaml)
- [Audit frontend](docs/audits/FRONTEND_AUDIT.md)
- [Jalur belajar developer](docs/development/LEARNING_GUIDE.md)
- [Pedoman kontribusi](docs/development/GUIDELINE.md)
- [Panduan maintenance](docs/development/MAINTENANCE.md)

## Catatan status

Dokumen desain dan roadmap dapat memuat rencana yang belum diimplementasikan. Periksa [System Overview](docs/system/SYSTEM_OVERVIEW.md), route API aktif, dan UI aktual sebelum mengandalkan fitur untuk operasi produksi.
