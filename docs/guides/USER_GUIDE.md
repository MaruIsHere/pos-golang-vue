# Panduan Penggunaan — POS Kasir Pro

Panduan ini mencakup persiapan aplikasi, penggunaan menu, dan pemecahan masalah untuk versi yang ada di repositori. Label menu dan hak akses mengikuti implementasi frontend/backend saat dokumen dibuat.

## 1. Persiapan dan menjalankan aplikasi

### Kebutuhan

- Go sesuai versi pada `backend/go.mod` (saat ini `go 1.27.1`).
- Node.js dan npm yang mendukung Vite pada `frontend/package.json`.
- Browser modern; izin kamera diperlukan jika memakai scanner kamera.
- SQLite berjalan tanpa server database tambahan. MySQL bersifat opsional.

### Development

Dari root proyek:

```bash
npm install
npm --prefix frontend install
npm run dev
```

Frontend tersedia di `http://localhost:3000`, API default di `http://localhost:8080`. Vite meneruskan `/api` dan `/uploads` ke backend. Jika belum ada file config lokal, backend membuat `backend/config.local.json` berdasarkan konfigurasi default.

Alternatif: `make run` menjalankan `npm run dev`.

### Build dan jalankan mode server

```bash
npm run build
npm start
```

Build menghasilkan aset `frontend/dist` dan binary backend di `backend/pos-backend`. Backend menyajikan UI di port yang dikonfigurasi (default 8080). `npm start` membangun binary backend jika belum ada; build frontend perlu dilakukan setelah perubahan UI.

### Akun awal untuk development

Seed lokal membuat akun berikut jika belum tersedia:

| Username | Password awal | Role |
|---|---|---|
| `owner` | `owner123` | OWNER |
| `admin` | `admin123` | ADMIN |
| `kepalakasir` | `kepala123` | KEPALA_KASIR |
| `kasir` | `kasir123` | KASIR |

Kredensial ini hanya untuk database development/demo. Ganti atau hapus sebelum server dibuka ke jaringan. Jangan gunakan akun seed untuk data usaha sungguhan.

## 2. Hak akses ringkas

Menu yang tampil:

| Menu | Kasir | Kepala kasir | Owner/Admin |
|---|:---:|:---:|:---:|
| Kasir | ✓ | ✓ | ✓ |
| Riwayat pesanan | ✓ | ✓ | ✓ |
| Produk | — | ✓ | ✓ |
| Inventori | — | ✓ | ✓ |
| Pelanggan | — | ✓ | ✓ |
| Laporan | — | ✓ | ✓ |
| Toko | — | — | ✓ |
| Pengaturan | — | — | ✓ |

Daftar menu bukan pengganti matriks izin API. Beberapa operasi (misalnya refund) memiliki aturan tambahan di backend.

## 3. Login dan navigasi

1. Buka alamat aplikasi dan masukkan username serta password.
2. Owner/admin diarahkan ke laporan; role operasional diarahkan ke halaman Kasir.
3. Gunakan navbar untuk membuka menu yang tersedia. Pada layar kecil, menu tambahan dapat berada dalam menu navigasi sekunder.
4. Gunakan Logout sebelum meninggalkan perangkat bersama.

## 4. Proses penjualan

1. Pastikan nama toko aktif yang tampil sesuai lokasi kerja. Owner/admin dapat memilih konteks toko dari kontrol toko pada navbar.
2. Cari produk berdasarkan nama, merk, tipe, barcode, atau scan barcode.
3. Tekan kartu produk untuk menambahkannya ke keranjang.
4. Atur kuantitas dengan tombol plus/minus untuk pcs, atau input kuantitas hingga tiga angka desimal untuk gram/liter. Jumlah tidak dapat melebihi stok yang dimuat di katalog.
5. Periksa subtotal, pajak, dan diskon yang muncul pada keranjang.
6. Tekan **Proses Bayar**, pilih Tunai, QRIS, atau Transfer Bank. Untuk non-tunai, unggah bukti jika diperlukan.
7. Isi atau pilih pelanggan bila dibutuhkan, lalu selesaikan transaksi.
8. Setelah API menyatakan transaksi sukses, tinjau struk. Stok diperbarui dari data server.

### Catatan pembayaran

- OCR pada unggahan bukti hanya membantu mendeteksi angka. Kasir tetap harus memverifikasi pembayaran melalui aplikasi/bank yang benar.
- Jangan menganggap gambar bukti sebagai konfirmasi dana telah diterima.
- Audit frontend tertanggal 11 Oktober 2026 mencatat potongan member belum tercermin pada diskon pesanan yang dikirim dan input voucher manual belum tersambung. Hindari mengandalkan dua perilaku tersebut sampai diperbaiki dan diverifikasi.

## 5. Katalog produk

1. Buka **Produk**.
2. Pilih **Produk Kasir Toko** untuk katalog yang dijual di toko aktif atau **Katalog Produk Utama** untuk katalog master.
3. Gunakan pencarian dan filter kategori, merk, serta tipe.
4. Tekan **Produk Baru** untuk menambah data. Pilih kategori, nama, harga, satuan, stok, dan data lain yang relevan.
5. Kelola kategori, merk, dan tipe lewat tombol **Master Data**.
6. Di katalog master, impor satu produk atau beberapa produk ke toko yang dipilih.
7. Produk dengan gambar harus memenuhi batas unggah yang ditampilkan aplikasi.

Periksa toko aktif sebelum membuat atau mengimpor produk agar stok/katalog masuk ke konteks yang benar.

## 6. Inventori dan mutasi stok

- **Penerimaan (Stock In):** pilih produk, jumlah masuk, dan alasan, lalu simpan.
- **Pengeluaran (Stock Out):** pilih produk, jumlah keluar, alasan (rusak, hilang, expired, promosi, pemakaian internal), dan catatan opsional.
- Kalkulator konversi membantu menghitung kg ke gram, ml ke liter, atau isi kemasan ke pcs. Tinjau satuan dan hasil sebelum menerapkan ke kuantitas.
- **Riwayat Mutasi:** menampilkan mutasi manual. Penjualan kasir merupakan pengurangan stok otomatis dan tidak ditampilkan sebagai mutasi manual pada layar ini.

## 7. Riwayat dan refund

- Buka **Riwayat** untuk melihat pesanan terbaru.
- Buka bukti pembayaran bila tersedia dan gunakan tindakan struk untuk melihat/cetak rincian.
- Refund tersedia untuk role yang diizinkan backend. Konfirmasi nomor/invoice dan dampak stok sebelum menyetujui refund.
- Jangan mengulang refund jika hasil request belum jelas; refresh riwayat terlebih dahulu.

## 8. Pelanggan dan manfaat member

- Buka **Pelanggan** untuk menambah, mengubah, mencari, atau menghapus data pelanggan sesuai hak akses.
- Saat checkout, pilih pelanggan yang sudah terdaftar atau gunakan pilihan pelanggan umum.
- Data pelanggan berisi informasi pribadi; masukkan hanya informasi yang diperlukan dan batasi akses akun staf.

## 9. Laporan dan ekspor

- Buka **Laporan** untuk ringkasan transaksi, grafik, produk terlaris/kurang laku, dan tabel penjualan.
- Tekan **Refresh** untuk mengambil data terbaru.
- Gunakan **Export Excel**, **Export CSV**, atau **Export PDF** untuk mengunduh laporan yang sedang ditampilkan.
- Perlakukan hasil ekspor sebagai data bisnis sensitif dan simpan di lokasi yang aksesnya dibatasi.

## 10. Toko dan staf (Owner/Admin)

- Di **Toko**, buat atau perbarui toko dan atur staf yang ditempatkan di toko.
- Di **Pengaturan → Tim**, buat akun dan kelola profil/password staf.
- Gunakan role minimum yang dibutuhkan untuk tugas staf.
- Setelah mengubah konteks toko, pastikan nama toko aktif benar sebelum memproses penjualan.

## 11. Pengaturan

- **Toko/Profil:** kelola identitas usaha dan informasi yang digunakan pada struk.
- **Voucher:** buat voucher bertipe persentase atau nominal sesuai opsi form. Periksa tanggal/ketentuan pada versi aplikasi; UI saat ini tidak menjanjikan mekanisme masa berlaku atau batas penggunaan.
- **Tim:** kelola akun staf yang tersedia.
- **Sistem:** pilih tema dan konfigurasi engine database.

Pergantian database bukan migrasi data otomatis yang boleh dilakukan tanpa persiapan. Buat backup, siapkan database/DSN MySQL, dan rencanakan perpindahan data sebelum mengubah engine.

## 12. Instalasi PWA dan dukungan offline

Aplikasi menawarkan banner instalasi PWA dan menyimpan aset statis tertentu untuk membuka shell aplikasi. Pencatatan penjualan dan operasi utama membutuhkan API aktif; transaksi offline tidak disimpan untuk sinkronisasi kemudian.

## 13. Pemecahan masalah

| Gejala | Tindakan |
|---|---|
| Tidak bisa login | Periksa username/password, alamat server, lalu coba lagi. Minta owner reset password jika akun tidak diketahui. |
| Pesan koneksi/server | Pastikan backend aktif pada port yang dikonfigurasi dan Vite proxy mengarah ke alamat yang benar. |
| Menu tidak tersedia | Role akun mungkin tidak memiliki izin. Hubungi owner/admin. |
| Produk tidak terlihat | Periksa toko aktif, mode katalog, filter pencarian, kategori, serta status stok. |
| Kamera barcode gagal | Beri izin kamera dan gunakan HTTPS atau localhost; gunakan input barcode sebagai alternatif. |
| Bukti/unggahan gagal | Periksa koneksi, ukuran/format file, dan kapasitas penyimpanan server. |
| UI server tidak berubah setelah edit | Jalankan ulang `npm run build`; server menyajikan file di `frontend/dist`. |
| Data toko tampak salah setelah switch | Hentikan perubahan, pastikan config/engine dan database yang aktif, lalu minta administrator memeriksa backup. |

## 14. Backup dan perlindungan data

- Untuk SQLite, hentikan aplikasi atau gunakan mekanisme backup konsisten sebelum menyalin database.
- Sertakan file database dan direktori `backend/uploads/` jika backup harus mempertahankan gambar/bukti.
- Untuk MySQL, gunakan prosedur backup database yang dikelola administrator.
- Uji pemulihan backup secara berkala; file backup yang belum pernah diuji belum tentu dapat dipulihkan.
