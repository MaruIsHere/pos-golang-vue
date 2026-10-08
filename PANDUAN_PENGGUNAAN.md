# 📚 Panduan Penggunaan Sistem (User Documentation)
**POS Kasir Pro - Multi-Tenant SaaS (v1.0.0-MVP)**

Dokumen ini berisi panduan operasional lengkap dari sisi pengguna dan administrator untuk menjalankan sistem kasir ini.

---

## 1. 🚀 Cara Menjalankan Aplikasi
Aplikasi ini sangat mudah dijalankan. Buka terminal di *root* folder proyek, lalu gunakan perintah berikut:

- **Menjalankan Aplikasi (Mode Development):**
  ```bash
  make run
  ```
  *Akses frontend di `http://localhost:3000` (atau `http://localhost:5173` jika menggunakan default Vite tanpa konfigurasi).*

- **Reset Database Total (Migrate Fresh):**
  ```bash
  make migrate-fresh
  ```
  *Gunakan ini jika ingin menghapus seluruh data dan mereset skema database ke awal. (Mirip `php artisan migrate:fresh`).*

- **Inject Data Dummy (Seeder):**
  ```bash
  make db-seed
  ```
  *Gunakan ini untuk memasukkan data awal seperti Akun Karyawan, Outlet, dan Produk contoh. (Mirip `php artisan db:seed`).*

---

## 2. 🔐 Akun & Hak Akses (Role-Based Access)
Sistem sudah dilengkapi dengan 4 akun default saat di-seed. **Password untuk semua akun di bawah ini adalah:** `[username]123` *(contoh: password owner adalah `owner123`)*.

| Username | Jabatan | Akses & Wewenang |
|---|---|---|
| **`owner`** | Pemilik Toko / Boss | Akses 100%. Bisa menambah karyawan baru, melihat semua laporan, dan mengatur toko. |
| **`admin`** | Back-Office | Bisa melihat laporan dan atur stok barang, tapi **TIDAK BISA** menambah karyawan baru. |
| **`kepalakasir`** | Supervisor Kasir | Bisa mengatur inventaris dasar dan memproses refund transaksi. |
| **`kasir`** | Staf Penjaga Lapak | Akses paling rendah. Hanya bisa masuk ke layar **KASIR (POS)** untuk memproses pesanan. |

---

## 3. 🛒 Alur Kerja Kasir (Lapak Event / Reguler)
1. **Buka Layar Kasir:** Kasir login menggunakan akun `kasir`, sistem akan langsung membuka halaman POS (Point of Sale).
2. **Pilih Barang:** Klik *card* barang yang dibeli (misal: Takoyaki).
3. **Checkout & Pembayaran:** 
   - Klik tombol **Bayar**.
   - Pilih metode pembayaran: **Cash** atau **QRIS / Transfer**.
   - *(Fitur Baru!)* Jika pelanggan membayar menggunakan QRIS, kasir bisa mencatat **Nomor Referensi** (contoh: 4 digit terakhir nomor HP pembeli) sebagai bukti jika terjadi selisih dengan pihak Event Organizer.
   - Jika Event Organizer memotong biaya admin (Fee), masukkan potongannya di kolom **Platform Fee**.
4. **Selesai:** Stok barang akan otomatis berkurang.

---

## 4. 📦 Manajemen Inventaris (Stock)
- Buka menu **Inventaris**.
- Semua *Stock Movement* (Barang Masuk & Keluar) dicatat dengan ketat.
- Terdapat kalkulator konversi otomatis (Misal: Kasir membeli bahan baku tepung 1 Kg, sistem bisa mengkonversinya menjadi satuan 1000 gram agar akurat saat digunakan sebagai *recipe* HPP).

---

## 5. 📊 Laporan Harian (Dashboard)
Sistem laporan sudah dirancang secara spesifik untuk rekonsiliasi uang di laci kasir vs uang di rekening bank.
- Buka menu **Laporan (Dashboard)**.
- Anda akan melihat pemisahan yang sangat detail:
  - **Total Revenue (Omset Kotor):** Keseluruhan uang dari transaksi hari ini.
  - **Total Cash (Uang Fisik):** Uang yang seharusnya **ada di laci kasir** malam ini. 
  - **Total QRIS / Debit:** Uang yang mengendap di bank / ditahan pihak Event Organizer (cair T+1).
  - **Total Platform Fee:** Uang yang hilang dipotong oleh biaya MDR Bank atau Fee Panitia Event (Penting agar perhitungan Laba Bersih akurat).

---

## 6. 📱 Dukungan Perangkat (Responsive)
Aplikasi ini dapat dibuka di berbagai perangkat tanpa mengorbankan kenyamanan:
- **Desktop/Laptop:** Tampilan layar penuh, keranjang belanja (Cart) terbuka di sebelah kanan secara permanen.
- **Tablet (iPad):** Resolusi sempurna untuk dijadikan *Kiosk* atau *Stand POS* di meja lapak.
- **Smartphone (HP):** Tampilan mengecil dengan mulus, keranjang belanja berubah menjadi tombol melayang (Drawer) di bawah layar agar kasir leluasa memilih menu walau dari layar kecil.
