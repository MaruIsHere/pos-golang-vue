# Jalur Belajar Developer

Dokumen ini membantu mempelajari repositori lewat satu alur fitur. Ikuti kode aktual dan tes yang tersedia; beberapa dokumen di `docs/archive/` adalah catatan lama.

## 1. Mulai dari peta aplikasi

1. Baca [README](../../README.md) untuk perintah utama.
2. Baca [System Overview](../system/SYSTEM_OVERVIEW.md) untuk arsitektur dan batas fitur.
3. Baca [PRD](../product/PRD.md) untuk tujuan produk dan ruang lingkup.
4. Buka `frontend/src/router/index.ts` untuk peta halaman.
5. Buka `backend/internal/routes/routes.go` untuk peta API dan middleware.

## 2. Pelajari alur transaksi

Telusuri berurutan:

1. `frontend/src/views/PosView.vue` — pemuatan produk, keranjang, dan request checkout.
2. `frontend/src/components/CartDrawer.vue` — ringkasan keranjang dan diskon/voucher.
3. `frontend/src/components/PaymentModal.vue` — detail pembayaran.
4. `frontend/src/utils/api.ts` — bearer token dan penanganan respons.
5. `backend/internal/routes/routes.go` — route `POST /orders`.
6. `backend/internal/handlers/order.go` — validasi, perhitungan transaksi, penyimpanan order/item, stok.
7. `backend/internal/models/models.go` — struktur `Order` dan `OrderItem`.

Saat membaca, catat kontrak payload frontend dan bandingkan dengan field request backend. Nilai uang, diskon, pajak, outlet, dan stok harus konsisten dari UI sampai database.

## 3. Pelajari autentikasi

- Frontend: `frontend/src/stores/auth.ts`, `frontend/src/router/index.ts`, `frontend/src/utils/api.ts`.
- Backend: `backend/internal/handlers/auth.go`, `backend/internal/middleware/auth.go`, dan deklarasi middleware/role di `backend/internal/routes/routes.go`.
- Periksa khususnya bahwa setiap operasi terlindungi di backend, bukan hanya disembunyikan dari menu.

## 4. Perintah verifikasi yang tersedia

```bash
npm --prefix frontend run typecheck
npm --prefix frontend run build
cd backend && go test ./...
cd backend && go vet ./...
```

Jalankan test yang relevan setelah perubahan. Untuk alur UI yang tidak tercakup test otomatis, lakukan smoke test dengan akun dan database development.

## 5. Panduan pemeliharaan

Lihat [MAINTENANCE](MAINTENANCE.md). Panduan deployment aspiratif ada di [`deployment/production-standard.md`](deployment/production-standard.md); cocokkan semua rekomendasi dengan kontrol yang benar-benar ada di kode.

## 6. Aturan dokumentasi

- Tandai keputusan sebagai **implementasi**, **rencana**, atau **saran**.
- Jangan menyebut mode PWA sebagai transaksi offline kecuali sinkronisasi API sudah dibangun dan diuji.
- Perbarui route/API dan panduan pengguna saat perilaku aplikasi berubah.
- Jangan dokumentasikan secret atau DSN produksi.
