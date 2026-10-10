# System Overview — POS Kasir Pro

Dokumen ini menjelaskan bagaimana bagian utama sistem terhubung dan menjadi bahan belajar untuk membaca repositori. Implementasi aktual di kode menjadi sumber kebenaran jika ada perbedaan.

## 1. Gambaran arsitektur

```text
Browser / PWA
  └─ Vue 3 + TypeScript + Vite + Pinia
       └─ Axios (/api, bearer JWT)
            └─ Go + Gin API (:8080)
                 └─ GORM
                      ├─ SQLite (default)
                      └─ MySQL (opsi konfigurasi)

Backend juga menyajikan frontend/dist dan file unggahan saat aplikasi dijalankan sebagai server.
OCR bukti bayar dipanggil sebagai executable lokal `bin/ocrs` atau `bin/ocrs.exe`.
```

Frontend tidak di-embed ke binary Go pada konfigurasi saat ini. Build membuat `frontend/dist`; server Go menyajikan direktori tersebut dari path relatif.

## 2. Struktur repositori

| Path | Isi |
|---|---|
| `frontend/src/views/` | Halaman yang dipetakan router |
| `frontend/src/components/` | Komponen UI dan alur POS |
| `frontend/src/stores/` | State auth, pengaturan, dan konteks toko |
| `frontend/src/utils/` | Klien API dan utilitas ekspor |
| `backend/cmd/api/` | Entrypoint HTTP API |
| `backend/cmd/migrate/`, `backend/cmd/seed/` | Utilitas migrasi dan seed |
| `backend/internal/routes/` | Deklarasi route dan middleware per route |
| `backend/internal/handlers/` | Validasi request dan logika endpoint |
| `backend/internal/models/` | Model GORM dan bentuk data API |
| `backend/internal/middleware/` | Autentikasi JWT dan otorisasi role |
| `backend/internal/config/` | Konfigurasi JSON dan environment override |
| `backend/internal/database/` | Koneksi, AutoMigrate, dan seed awal |
| `scripts/` | Runner development dan start lintas platform |
| `bin/` | Executable OCR untuk platform yang disediakan |

## 3. Aliran request

1. Vue Router memuat view secara lazy dan menjaga halaman yang memerlukan autentikasi.
2. Store auth membaca token dari `localStorage`.
3. Klien Axios menambahkan `Authorization: Bearer <token>` pada request.
4. Gin route mengarahkan request ke handler; route terlindungi memakai middleware autentikasi dan role.
5. Handler mengakses model melalui GORM dan mengembalikan JSON.
6. Store/view memperbarui state tampilan berdasarkan respons.

Respons `401` akan mengakhiri sesi frontend. Validasi dan izin di UI bukan batas keamanan; API harus selalu dianggap sebagai otoritas.

## 4. Model domain utama

- **Merchant:** ruang usaha/tenant.
- **Outlet:** lokasi toko yang terkait merchant.
- **User:** staf dengan role dan kemungkinan outlet.
- **Product:** produk katalog master atau produk toko; produk toko dapat terhubung ke master melalui `master_product_id`.
- **Category, Artist, ProductType:** atribut katalog.
- **Order dan OrderItem:** transaksi serta snapshot item saat transaksi.
- **StockMovement:** catatan perubahan stok masuk/keluar.
- **Customer, Voucher, StoreSetting:** data pendukung penjualan dan konfigurasi.
- **Expense:** model tersedia, tetapi modul pencatatan pengeluaran operasional belum tampak pada route/API aktif.

Relasi merchant/outlet ada pada model, tetapi cakupan isolasi perlu ditinjau per handler sebelum sistem digunakan sebagai SaaS multi-tenant untuk banyak pelanggan.

## 5. Route frontend dan API

Route frontend utama didefinisikan di `frontend/src/router/index.ts`: `/`, `/orders`, `/products`, `/inventory`, `/customers`, `/reports`, `/settings`, dan `/stores`; autentikasi menggunakan `/login`.

Route API didefinisikan di `backend/internal/routes/routes.go`. Kelompok fungsinya:

- Auth: login.
- Profil dan pengguna.
- Katalog: kategori, produk, merk, tipe produk, gambar produk.
- Penjualan: order, detail/riwayat, upload bukti, refund.
- Inventori: stock movements.
- Pelanggan dan voucher.
- Laporan dashboard.
- Toko, pengaturan, dan switch database.

Spesifikasi OpenAPI berada di [`../api/openapi.yaml`](../api/openapi.yaml); cocokkan status aktualnya dengan route sebelum menjadikannya kontrak final.

## 6. Konfigurasi dan penyimpanan

- Default disimpan pada `backend/config.json`.
- `backend/config.local.json` dibuat untuk konfigurasi lokal dan diabaikan Git.
- Environment override: `POS_PORT`, `POS_DB_ENGINE`, `POS_SQLITE_PATH`, `POS_MYSQL_DSN`.
- SQLite memakai path relatif terhadap folder config, WAL, dan satu koneksi terbuka untuk menserialisasi penulisan.
- Upload disimpan di `backend/uploads/` saat proses dijalankan dari folder backend.

Jangan memasukkan kredensial MySQL produksi atau file konfigurasi lokal ke Git. Backup database sebelum menjalankan migrasi atau mengganti engine.

## 7. Menjalankan sistem

```bash
# Development: Vite :3000 dan API Go :8080
npm run dev

# Build frontend dan backend
npm run build

# Menjalankan backend melalui runner lintas platform
npm start
```

Vite mem-proxy `/api` dan `/uploads` ke `http://127.0.0.1:8080` secara default. Target proxy dapat diubah dengan `VITE_API_PROXY_TARGET`.

## 8. Keamanan dan batas operasional

- JWT disimpan di browser `localStorage`; lindungi perangkat kasir dan gunakan HTTPS di jaringan produksi.
- Seed development membuat akun default yang kata sandinya tercantum di panduan pengguna. Ganti/hapus akun tersebut sebelum deployment yang dapat diakses jaringan.
- CORS pada entrypoint saat ini mengizinkan semua origin tanpa credentials; evaluasi dan batasi origin pada deployment publik.
- SQLite cocok untuk satu server sederhana; sistem tidak menyediakan replikasi atau sinkronisasi offline.
- PWA cache aset statis, bukan database transaksi lokal.
- OCR berjalan lokal dan hasil pembacaan harus diperiksa oleh petugas.

## 9. Jalur belajar kode yang disarankan

1. `backend/cmd/api/main.go` — startup, konfigurasi, middleware, static serving.
2. `backend/internal/routes/routes.go` — peta API dan role per route.
3. `backend/internal/middleware/auth.go` — pembentukan konteks autentikasi.
4. `backend/internal/models/models.go` — skema domain.
5. Satu alur ujung-ke-ujung: `frontend/src/views/PosView.vue` → `POST /api/orders` → `backend/internal/handlers/order.go`.
6. `frontend/src/utils/api.ts` dan `frontend/src/stores/auth.ts` — komunikasi API dan sesi UI.

## 10. Batasan yang diketahui

- Unit/integration test backend masih terbatas; cakupan penuh belum dijamin.
- Hasil audit frontend mencatat dua masalah pada diskon member dan input voucher manual: [`../audits/FRONTEND_AUDIT.md`](../audits/FRONTEND_AUDIT.md).
- Build frontend menghasilkan chunk laporan besar; lihat laporan audit untuk ukuran yang terukur saat audit.
- Dokumen roadmap SaaS bukan bukti bahwa seluruh isolasi tenant sudah diterapkan.
