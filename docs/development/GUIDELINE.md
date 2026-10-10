# Pedoman Kontribusi

Gunakan pedoman ini saat mengubah POS Kasir Pro. Ini adalah standar kerja proyek; klaim tentang perilaku yang sudah berjalan harus dicek pada kode dan test.

## 1. Sebelum mengubah kode

1. Periksa `git status` agar perubahan lokal yang sudah ada tidak tertimpa.
2. Baca [System Overview](../system/SYSTEM_OVERVIEW.md) dan [PRD](../product/PRD.md) yang relevan.
3. Ikuti alur fitur dari view frontend ke route, middleware, handler, dan model terkait.
4. Pastikan perubahan tidak mencampur data merchant atau outlet.

## 2. Backend Go

- Letakkan entrypoint di `backend/cmd/`; logika aplikasi di `backend/internal/`.
- Gunakan UUID sesuai model yang ada; representasi ID di frontend berupa string.
- Ambil merchant/outlet dari konteks middleware, bukan dari nilai request yang dipercaya begitu saja.
- Pastikan query baca/tulis membatasi data ke merchant dan outlet yang berhak diakses. Jangan menganggap seluruh handler sudah memenuhi isolasi SaaS; cek satu per satu.
- Validasi input dan otorisasi di backend. Menyembunyikan tombol/menu frontend tidak cukup.
- Gunakan transaksi database untuk operasi yang mengubah beberapa entitas atau stok secara bersamaan.
- Tangani error dengan status HTTP yang sesuai; jangan membocorkan secret atau DSN.

## 3. Frontend Vue

- Gunakan Vue 3 Composition API dan `<script setup lang="ts">` mengikuti pola file terkait.
- Gunakan komponen UI bersama dan composable/store yang sudah tersedia jika sesuai.
- Jaga state loading, kosong, error, dan sukses.
- Pertahankan pengalaman keyboard, label form, fokus, dan layout responsif.
- Hindari `any` tanpa alasan kuat; perbarui tipe API jika bentuk payload berubah.
- Jangan mempercayai perhitungan harga, diskon, pajak, stok, atau role dari browser sebagai sumber final.

## 4. Perubahan transaksi dan stok

Untuk perubahan checkout, periksa konsistensi antara UI, payload, handler, database, laporan, dan struk. Verifikasi subtotal, diskon, pajak, total akhir, jumlah dibayar, kembalian, outlet, dan efek stok.

## 5. Verifikasi

```bash
npm --prefix frontend run typecheck
npm --prefix frontend run build
cd backend && go test ./...
cd backend && go vet ./...
```

Tambahkan atau perbarui test untuk bug yang diperbaiki. Untuk perubahan visual/perilaku UI yang tidak tercakup test, dokumentasikan skenario smoke test.

## 6. Dokumentasi

- Perbarui [PRD](../product/PRD.md) bila ruang lingkup produk berubah.
- Perbarui [desain](../design/DESIGN.md) bila pola UI/alur berubah.
- Perbarui [System Overview](../system/SYSTEM_OVERVIEW.md) bila arsitektur atau konfigurasi berubah.
- Perbarui [Panduan Pengguna](../guides/USER_GUIDE.md) bila langkah operasional berubah.
- Tandai hal yang masih rencana; jangan menggambarkan roadmap sebagai fitur yang sudah tersedia.
