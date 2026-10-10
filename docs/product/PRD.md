# Product Requirements Document — POS Kasir Pro

**Status:** Gambaran produk yang berjalan (MVP), 11 Oktober 2026  
**Bahasa:** Indonesia

## 1. Ringkasan

POS Kasir Pro adalah aplikasi web kasir untuk usaha kecil dengan satu atau beberapa toko. Aplikasi menyediakan pencatatan penjualan, katalog produk, mutasi stok, pelanggan, voucher, laporan, dan pengelolaan pengguna. Frontend dibuat dengan Vue 3 dan berkomunikasi dengan REST API Go/Gin. Data disimpan di SQLite secara default atau MySQL jika dikonfigurasi.

## 2. Masalah yang ingin diselesaikan

- Kasir perlu mencatat transaksi dengan cepat dan melihat ketersediaan stok.
- Pemilik perlu mengelola katalog, pengguna, toko, dan ringkasan penjualan dari satu aplikasi.
- Perubahan stok di luar penjualan perlu tercatat sebagai mutasi dengan alasan.
- Usaha multi-toko perlu memilih konteks toko dan memisahkan produk toko dari katalog master.

## 3. Pengguna dan kebutuhan

| Peran | Kebutuhan utama | Akses antarmuka utama |
|---|---|---|
| Kasir (`KASIR`) | Menjual produk, menerima pembayaran, melihat riwayat | Kasir, Riwayat |
| Kepala kasir (`KEPALA_KASIR`) | Mengawasi kasir, mengelola produk/stok/pelanggan/laporan dan refund | Kasir, Produk, Inventori, Pelanggan, Riwayat, Laporan |
| Owner (`OWNER`) | Mengelola usaha, toko, staf, konfigurasi, dan laporan | Semua menu |
| Admin (`ADMIN`) | Mengelola operasional dan konfigurasi sesuai kebijakan aplikasi | Semua menu |
| Administrator (`ADMINISTRATOR`) | Role legacy yang dikenali sebagian middleware API | Akses harus diverifikasi; belum ditampilkan sebagai role bawaan di menu frontend |

Hak akses route frontend dan middleware API tidak sepenuhnya memakai daftar role yang sama. Penegakan keamanan harus selalu dilakukan di backend; tabel di atas merangkum menu yang ditampilkan, bukan pengganti pemeriksaan otorisasi API.

## 4. Ruang lingkup fitur yang tersedia

### 4.1 Penjualan

- Cari produk dan tambahkan ke keranjang.
- Kelola kuantitas, hapus item, dan kosongkan keranjang.
- Hitung subtotal, diskon keranjang, dan pajak berdasarkan pengaturan.
- Proses tunai, QRIS, atau transfer; unggah bukti pembayaran non-tunai.
- Cetak/tampilkan struk setelah pesanan berhasil.
- Scan barcode melalui input pencarian atau kamera pada perangkat yang mendukung.

### 4.2 Katalog dan inventori

- Kelola produk, kategori bertingkat, merk/artist, dan tipe produk.
- Pisahkan produk toko dari katalog produk master; impor produk master ke toko.
- Kelola stok masuk/keluar dan lihat histori mutasi.
- Satuan stok yang didukung: pcs, gram, dan liter.

### 4.3 Operasional usaha

- Kelola toko dan penempatan pengguna ke toko.
- Kelola pelanggan, voucher, profil, pengaturan toko, dan tema.
- Lihat ringkasan dan grafik penjualan serta ekspor Excel, CSV, dan PDF.
- Ubah engine database melalui halaman pengaturan untuk deployment yang sudah menyiapkan koneksi MySQL.

## 5. Kebutuhan non-fungsional

- Antarmuka responsif untuk desktop dan perangkat sentuh.
- Navigasi route memerlukan autentikasi dan memeriksa role; API terlindungi menggunakan bearer token.
- Aplikasi dapat dipasang sebagai PWA dan menyimpan aset statis untuk pemuatan shell. **Ini bukan mode transaksi offline**: sinkronisasi transaksi saat offline belum termasuk.
- Build frontend harus lolos `vue-tsc` dan Vite; backend dibangun dengan toolchain Go yang tercantum di `backend/go.mod`.
- Gambar produk dibatasi dan dikompresi sebelum diunggah.

## 6. Batasan dan hal yang belum dijanjikan

- Tidak ada sinkronisasi transaksi offline atau resolusi konflik multi-perangkat.
- Belum ada modul resep/BOM bahan baku, akuntansi laba-rugi penuh, atau pencatatan biaya operasional yang terhubung ke laporan.
- Tidak ada sistem langganan/billing tenant, pemulihan kata sandi, atau rotasi refresh token yang terdokumentasi.
- OCR dijalankan lokal melalui executable `ocrs` jika tersedia; kualitas hasil pembacaan bergantung pada gambar dan model OCR. OCR bukan verifikasi transfer otomatis.
- Beberapa teks pada antarmuka mungkin menyebut fitur lebih luas daripada implementasi backend; cek perilaku API sebelum mengandalkan fitur untuk operasi keuangan.

## 7. Kriteria penerimaan utama

1. Pengguna dapat login, dan akses halaman/menu dibatasi menurut role.
2. Kasir dapat membuat pesanan berisi satu atau lebih produk dan melihat struk hasil transaksi.
3. Transaksi mengurangi stok produk toko terkait dan muncul pada riwayat/laporan.
4. Pengguna berwenang dapat membuat produk dan mencatat mutasi stok yang valid.
5. Owner/admin dapat mengelola toko, staf, dan pengaturan yang tersedia.
6. Ekspor laporan menghasilkan berkas yang dapat dibuka pada format yang dipilih.

## 8. Metrik keberhasilan yang disarankan

- Waktu dari membuka POS sampai transaksi tersimpan.
- Persentase transaksi yang gagal dan alasan kegagalan.
- Selisih stok fisik dengan stok sistem.
- Keberhasilan ekspor laporan dan waktu pemuatan dashboard.
- Jumlah insiden akses lintas toko/merchant.

Metrik ini merupakan usulan pengukuran; aplikasi belum mendokumentasikan telemetri untuk semuanya.
