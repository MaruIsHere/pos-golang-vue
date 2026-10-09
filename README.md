# 🚀 POS Kasir Pro (Golang + Vue 3 + Tailwind v4 + AI OCR)

POS Kasir Pro adalah aplikasi Point of Sale (POS) modern yang dirancang untuk kecepatan, kemandirian (offline-first), dan efisiensi resource. Dibangun dengan arsitektur **Micro-binary**, seluruh aplikasi (Backend & Frontend) dapat dikemas ke dalam satu *executable file* yang ringan, didukung oleh mesin AI (OCR) bawaan untuk membaca struk tanpa butuh koneksi API eksternal.

---

## 🏗️ Arsitektur & Teknologi

Sistem dirancang agar kompatibel dengan lingkungan dengan *resource* rendah (seperti CPU lama) namun tetap mempertahankan standar modern.

*   **Backend:** Go (Golang) 1.22+ dengan Gin Framework dan GORM.
*   **Frontend:** Vue 3 (Composition API), Vite, Tailwind CSS v4, Pinia.
*   **Database:** SQLite (Default, *Zero-config* via `glebarez/sqlite`) dengan opsi beralih ke MySQL.
*   **AI Engine (OCR):** `ocrs-cli` (Pure Rust Tensor Engine). Ukuran binary hanya ~19MB, 100% lokal, didesain untuk mendeteksi nominal uang.
*   **UI/UX:** Skeleton Loading modern (menghindari spinner *laggy*), PWA-ready, Responsif (Desktop & Mobile), dan Dark Mode.

---

## ✨ Fitur Utama

1.  **🤖 Smart AI OCR Struk & QRIS**
    Unggah bukti transfer atau struk, dan sistem akan membaca nominal angka secara otomatis dan mengisinya di form pembayaran berkat integrasi Neural Network lokal (`ocrs-cli`).
2.  **⚡ Micro-binary & Multi-platform**
    Frontend Vue di-build menjadi file statis (`dist`), kemudian di-*embed* ke dalam binary Go. Anda dapat mendistribusikan aplikasi sebagai satu file tunggal untuk Windows (`.exe`), Linux, atau macOS.
3.  **📦 Manajemen Inventaris Kasir (POS)**
    *   Sistem Keranjang, Kalkulasi Pajak, dan Diskon.
    *   Multi-Metode Pembayaran (Tunai, QRIS, Transfer Bank).
    *   Mutasi stok keluar-masuk secara *realtime*.
4.  **🔐 Role-Based Access Control (RBAC)**
    Akses bertingkat untuk keamanan toko: `OWNER`, `ADMINISTRATOR`, `ADMIN`, `KEPALA_KASIR`, `KASIR`.
5.  **🗄️ Dynamic Database Switcher**
    Berjalan langsung dengan SQLite saat diunduh, dan bisa dialihkan ke MySQL/MariaDB dari menu *Settings* untuk kebutuhan multi-kasir (*production* skala menengah).

---

## 🛠️ Panduan Instalasi & Development

Proyek ini menggunakan _monorepo style_ dengan folder `backend/` dan `frontend/`.

### Persyaratan
*   [Go 1.22+](https://go.dev/)
*   [Node.js 18+](https://nodejs.org/) & `npm`
*   *(Opsional)* Rust/Cargo jika ingin mengompilasi ulang mesin OCR.

### Menjalankan Mode Development
Kami telah menyiapkan *script* otomatis (`npm run dev` di *root* folder) untuk menyalakan frontend (Vite) dan backend (Go) secara bersamaan:

```bash
# Clone repositori
git clone <url-repo>
cd pos-golang-vue

# Install dependency utama
npm install

# Jalankan Backend dan Frontend bersamaan
npm run dev
```

*   **Frontend Vite:** Berjalan di `http://localhost:3000` (dengan *proxy* internal ke backend).
*   **Backend Go API:** Berjalan di `http://localhost:8080`.

**Pengujian di Mobile via Cloudflare Tunnel:**
Jika Anda ingin menguji responsivitas UI dari HP secara langsung, jalankan:
```bash
cloudflared tunnel --url http://localhost:3000
```
Lalu buka *link* `.trycloudflare.com` yang dihasilkan di browser HP Anda.

---

## ⚙️ Kompilasi Mesin OCR (Opsional)

Kami telah menyertakan *binary* OCR untuk Linux (`bin/ocrs`) dan Windows (`bin/ocrs.exe`). Go akan otomatis mendeteksi OS yang berjalan (melalui `runtime.GOOS`).
Namun, jika Anda perlu menargetkan arsitektur lain (misalnya Apple Silicon/ARM), Anda bisa melakukan kompilasi ulang dari *source* Rust:

```bash
# Pastikan Rust terinstal (https://rustup.rs/)
cargo install ocrs-cli
# Salin binary ke folder aplikasi
cp ~/.cargo/bin/ocrs ./bin/
```

---

## 📖 Dokumentasi Lanjutan

Untuk memahami lebih dalam mengenai aturan *commit*, rancangan sistem, dan spesifikasi API, silakan lihat file-file berikut di folder `docs/`:

*   [Spesifikasi Kebutuhan Produk (PRD)](docs/PRD.md)
*   [Desain & Arsitektur Sistem](docs/DESIGN.md)
*   [Panduan Penggunaan](PANDUAN_PENGGUNAAN.md)
*   [OpenAPI/Swagger Spec](backend/docs/openapi.yaml)

---
*Dibuat dengan ❤️ untuk performa maksimal pada resource minimal.*
