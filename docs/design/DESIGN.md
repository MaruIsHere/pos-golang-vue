# Desain Produk dan UI — POS Kasir Pro

**Status:** Panduan desain berdasarkan UI yang ada. Dokumen ini bukan mockup pixel-perfect dan tidak mengesahkan fitur yang belum tersedia.

## 1. Prinsip desain

- **Utamakan alur kasir:** pencarian, pemilihan barang, ringkasan keranjang, dan pembayaran harus mudah ditemukan.
- **Tampilkan status dengan jelas:** loading, data kosong, error, sukses, dan tindakan berbahaya dibedakan.
- **Ramah sentuhan:** kontrol utama punya area klik yang memadai dan layout menyesuaikan ukuran layar.
- **Konsisten:** gunakan komponen UI bersama untuk tombol, input, badge, kartu, tabel, dan dialog.
- **Aksesibel:** label input, fokus keyboard, teks alternatif gambar, dan pesan status perlu dipertahankan saat UI dikembangkan.

## 2. Peta layar

| Layar | Tujuan |
|---|---|
| Login | Autentikasi staf |
| Kasir | Cari produk, menyusun keranjang, checkout, struk |
| Riwayat | Meninjau pesanan, bukti bayar, struk, dan refund sesuai izin |
| Produk | Mengelola katalog toko/master, kategori, merk, tipe, gambar |
| Inventori | Penerimaan, pengeluaran, dan riwayat mutasi stok |
| Pelanggan | Mengelola data pelanggan |
| Laporan | Melihat agregat penjualan dan ekspor |
| Toko | Mengelola cabang dan penempatan staf |
| Pengaturan | Profil/toko, voucher, staf, tema, dan database |

Visibilitas menu berbeda menurut role. Keputusan otorisasi final tetap berada pada backend.

## 3. Alur kasir

```text
Login → Pilih/konteks toko → Cari atau scan produk → Tambah ke keranjang
     → Atur jumlah/diskon → Pilih metode pembayaran → Simpan pesanan
     → Tampilkan/cetak struk → Perbarui katalog dan stok
```

### Interaksi penting

- Desktop menampilkan daftar produk dan keranjang berdampingan.
- Layar kecil memakai pemicu keranjang mengambang dan drawer.
- Pesanan tidak boleh dianggap selesai sebelum API mengembalikan sukses.
- Dialog konfirmasi dipakai untuk tindakan yang berpotensi menghapus atau mengubah data.
- Pembayaran non-tunai dapat menyertakan bukti unggahan. OCR hanya membantu membaca nominal, bukan memvalidasi dana masuk.

## 4. Pola layout dan komponen

- `Navbar.vue` mengatur navigasi utama dan menu tambahan untuk layar kecil.
- `AppButton`, `AppInput`, `AppBadge`, `AppCard`, `AppTable`, dan `AppDialogHost` menjadi elemen bersama.
- View berada di `frontend/src/views/`; komponen lintas layar berada di `frontend/src/components/`.
- Tema terang/gelap memakai kelas Tailwind dan preferensi lokal.
- PWA memakai manifest dan precache aset build; data API tidak otomatis tersedia offline.

## 5. Form dan validasi

- Validasi sederhana dapat dilakukan di browser untuk memberi umpan balik cepat.
- Backend tetap harus memvalidasi nilai uang, kuantitas, role, merchant, dan outlet.
- Pesan error harus memberi langkah pemulihan (misalnya coba lagi atau perbaiki isian).
- Input kuantitas harus mengikuti satuan produk: pcs bilangan bulat; gram/liter hingga tiga desimal.

## 6. Ekspor dan laporan

Ekspor disediakan dalam Excel, CSV, dan PDF. Tabel laporan perlu menjaga header dan format angka yang jelas. Hindari memuat pustaka ekspor besar pada jalur awal aplikasi jika tidak digunakan; build Oktober 2026 mencatat chunk laporan sekitar 1,5 MB sebelum gzip.

## 7. Arah perbaikan desain

1. Perbaiki konsistensi status interaksi dan validasi di seluruh form.
2. Audit aksesibilitas keyboard, screen reader, kontras, dan zoom.
3. Pisahkan modul ekspor/laporan agar pemuatan layar kasir tetap ringan.
4. Buat desain empty/error/loading state bersama dan gunakan konsisten.
5. Validasi desain dengan skenario penggunaan nyata pada desktop, tablet, dan ponsel.

## 8. Referensi desain tambahan

Konsep marketplace yang belum tentu menjadi bagian UI produksi tersimpan di [`marketplace-design/`](marketplace-design/REFERENCE.md). Perlakukan sebagai eksplorasi desain, bukan spesifikasi implementasi aktif.
