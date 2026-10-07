# Panduan Pengembangan (Developer Guideline)
**Sistem POS Golang Vue - Arsitektur SaaS Multi-Tenant v1.0.0-MVP**

Dokumen ini ditulis untuk menyelaraskan *vibe coding* seluruh developer agar tetap berada dalam satu konteks arsitektur yang sama. Tolong baca dengan seksama sebelum melakukan perubahan pada *codebase*!

## 1. Arsitektur Multi-Tenant (Sangat Krusial)
Sistem ini telah bermigrasi dari aplikasi kasir toko tunggal menjadi **SaaS (Software as a Service)** untuk banyak pemilik bisnis sekaligus. 
Hierarki data kita adalah: `Merchant` (Pemilik Bisnis) -> `Outlet` (Cabang Toko) -> `User` / `Product` / `Order` dll.

**ATURAN EMAS DATABASE:**
- **JANGAN PERNAH** mengambil data dari database tanpa mem-filter berdasarkan `merchant_id` dan `outlet_id`. 
- Kegagalan menerapkan ini akan menyebabkan data Toko A bocor dan bisa diakses oleh Toko B.
- Saat membuat API baru, selalu ekstrak `merchant_id` dan `outlet_id` dari *JWT Token* via Context Gin:
  ```go
  merchantID, _ := c.Get("merchant_id")
  outletID, _ := c.Get("outlet_id")
  // Gunakan dalam query GORM:
  db.Where("merchant_id = ?", merchantID)
  ```

## 2. Penggunaan UUID
Seluruh ID di dalam sistem (Primary Key maupun Foreign Key) **bukan lagi integer/angka**.
- Tipe data ID di Database dan Struct Golang = `uuid.UUID` (Package: `github.com/google/uuid`).
- Tipe data ID di Typescript/Vue (Frontend) = `string`.
- Jangan pernah menggunakan konversi `Number(id)` atau nilai *fallback* `0` di Frontend. Gunakan *empty string* `""` atau `null`.

## 3. Struktur Direktori Backend
Gunakan pola desain *Clean-ish Architecture* yang sudah ada. Jangan membuang file sembarangan.
- `cmd/api/main.go` -> *Entry point* aplikasi utama.
- `cmd/migrate/main.go` -> *Entry point* untuk skrip reset & migrasi database.
- `internal/models/` -> Tempat deklarasi struktur tabel Database (GORM).
- `internal/database/` -> Koneksi DB dan AutoMigrate.
- `internal/handlers/` -> Controller logika bisnis API.
- `internal/routes/` -> Registrasi *endpoint* API (Gin Router).

## 4. Pembayaran & Pelaporan (QRIS / Tunai)
- Sistem sudah membedakan pelaporan pendapatan Tunai (*Cash*) dan rekening (*QRIS/Debit*).
- Jika ada penambahan metode pembayaran, pastikan field `PaymentReference` dan `PlatformFee` di tabel `Order` ikut diperhitungkan di dalam agregasi SQL pada `internal/handlers/report.go`.

## 5. UI Frontend & CSS
- Aplikasi menggunakan **Vue 3 (Composition API / `<script setup>`)** dan **Tailwind CSS**.
- Sistem didesain untuk **Responsive**: Layout kasir (POS) akan melipat (menggunakan *Drawer* atau susunan vertikal) jika diakses dari *Mobile* atau *Tablet*. 
- Jika mengubah komponen UI, tes selalu responsivitas layar (ukuran iPhone SE hingga layar Desktop 1080p).
- Dilarang membypass *type checking* (`as any`). Pastikan `npx vue-tsc --noEmit` menghasilkan `0 error`.

---
*Stay sharp, stick to the context, and happy coding! 🚀*
