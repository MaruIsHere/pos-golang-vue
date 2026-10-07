# Blueprint Arsitektur SaaS (Multi-Tenant) POS Golang Vue

Dokumen ini menjadi **jangkar** dan panduan utama (roadmap) dalam merombak aplikasi dari *Single-Store* menjadi *True SaaS (Multi-Tenant)*, berdasarkan masukan UAT. Target market aplikasi ini adalah UMKM kecil, sehingga arsitektur harus tetap ringan (agile), mudah dikelola, dan tidak *over-engineered*.

## 🎯 Daftar Kebutuhan UAT
1. **Laporan Detail:** Mencakup laporan bahan baku, laba rugi, dan lainnya.
2. **Export CSV Rapi:** Tampilan export rapi dan mudah dibaca (formatting data).
3. **Pemisahan Inventory:** Membuat *Inventory Utama (Gudang)* dan *Inventory Lapak* agar stok barang tidak bercampur saat ada beberapa cabang/lapak.
4. **Penamaan per Kasir / SaaS:** Sistem harus bisa membedakan owner, lapak, dan kasir (Rule akun / Role-Based Access) yang mengarah ke sistem berlangganan (SaaS).

---

## 🏗️ 1. Skema Database (True SaaS Multi-Tenant)
Kita menggunakan pendekatan **Logical Isolation Multi-Tenancy**. Artinya, seluruh data klien berada dalam satu database, namun dipisahkan secara ketat menggunakan ID (Foreign Key).

### **Struktur Hirarki Utama:**
- **`Merchant` (Pemilik Usaha / Tenant):** Entitas tertinggi (misal: "Adi"). Semua data milik Adi akan diikat dengan `merchant_id`.
- **`Outlet` (Lapak / Cabang):** Cabang fisik yang dimiliki oleh Merchant. Diikat dengan `merchant_id` dan memiliki `outlet_id` sendiri.
- **`User` (Karyawan / Kasir):** Pekerja yang bertugas di lapak tertentu. Diikat dengan `merchant_id` dan `outlet_id`.

### **Aturan Isolasi Data (Wajib):**
Setiap tabel transaksi dan master data inti **WAJIB** memiliki `merchant_id` (dan/atau `outlet_id`) agar data satu Merchant (Adi) tidak akan pernah bocor ke Merchant lain.
Tabel yang wajib disesuaikan:
- `users`
- `products`
- `inventory`
- `transactions`
- `ingredients` / bahan baku

---

## 📦 2. Solusi Pemisahan Inventory (Poin 3)
Untuk memisahkan Gudang Utama dan Lapak, tabel `inventory` akan dirancang sebagai berikut:
1. Setiap row di `inventory` harus berelasi ke `outlet_id`.
2. Gudang utama dapat direpresentasikan sebagai **salah satu Outlet khusus** (misal dengan penanda `is_warehouse = true` di tabel `outlets`).
3. **Alur Mutasi Stok:** Barang diinput/masuk ke Gudang Utama -> Ditransfer (Mutasi) ke Lapak -> Terjual di Lapak (mengurangi stok lapak, bukan stok gudang).

---

## 💰 3. Strategi Laba Rugi & Tipe Produk (Poin 1)
Karena UMKM *event* memiliki jenis jualan yang beragam (FnB dan Merchandise), produk akan dibagi menjadi 2 tipe di dalam database (`product_type` di tabel `products`):
1. **Barang Jadi (Merchandise / Makanan Matang):**
   - Tidak menggunakan resep.
   - Penjualan memotong stok barang itu sendiri.
   - HPP (Modal) diambil langsung dari harga beli barang tersebut.
2. **Barang Olahan (FnB / Membutuhkan Resep):**
   - Memiliki resep (*Bill of Materials*).
   - Penjualan memotong stok *bahan baku* di gudang.
   - HPP diakumulasi dari total modal bahan baku yang terpakai.

**Metode Costing:** 
Sistem akan menggunakan **Harga Beli Terakhir (*Last Purchase Price*)** untuk menghitung HPP. Pendekatan ini dipilih karena komputasinya ringan untuk *database* dan jauh lebih mudah dipahami oleh UMKM.

---

## 📊 4. Strategi Export CSV (Poin 2)
Sistem akan menyediakan dua jenis laporan CSV untuk mencakup durasi *event* (1-3 hari) maupun toko permanen:
1. **Laporan Transaksi Detail:** 1 baris mewakili 1 struk/item terjual.
2. **Laporan Ringkasan Harian:** 1 baris mewakili akumulasi omset dan laba per hari.

**Format Angka & Mata Uang (Jalan Tengah):**
Format murni CSV tidak memiliki fitur *styling cell* seperti Excel (`.xlsx`). Jika nilai uang ditempel dengan tulisan "Rp" (misal: `Rp 15000`), program *spreadsheet* akan membacanya sebagai **Teks (Huruf)**, sehingga rumus `SUM()` dan kalkulasi keuangan akan rusak/gagal.
**Solusinya:** Simbol mata uang disematkan secara eksplisit pada **Header Kolom**. 
- *Contoh Header:* `Harga Modal (Rp) | Total Omset (Rp) | Laba Bersih (Rp)`
- *Contoh Row (Data):* `50000 | 150000 | 100000`

Dengan jalan tengah ini, angka yang di-*export* tetap murni (*Formula-Ready*), namun user awam tetap mendapat konteks yang jelas bahwa nominal tersebut adalah Rupiah.

---

## 🏢 5. Dukungan Akuntansi Badan Usaha (Ekspansi Poin 1)
Karena sistem menargetkan UMKM yang berpotensi memiliki legalitas Badan Usaha (CV/PT), struktur laporan akan diekspansi agar memenuhi standar akuntansi dasar:
1. **Pemecahan Komponen Transaksi:** Kolom laporan akan memisahkan Pendapatan Kotor (*Gross Sales*), Diskon, *Service Charge*, dan Pajak (PPN/PB1) agar pelaporan omset ke DJP valid.
2. **Rekapitulasi Metode Pembayaran:** Laporan wajib mencantumkan jenis pembayaran (Tunai, QRIS, Transfer Bank, EDC) untuk mempermudah UMKM melakukan rekonsiliasi (*cross-check*) dengan mutasi rekening bank mereka.
3. **Pencatatan Pengeluaran (*Expenses*):** Laba Kotor dari hasil jualan otomatis akan dikurangi dengan Biaya Operasional (Sewa lapak, Listrik, Gaji, dll) melalui modul pengeluaran untuk mendapatkan **Laba Bersih Operasional (*Net Profit*)**.

---

## 🚀 Langkah Selanjutnya (Next Action Plan)
- [ ] **Tahap 1: Database Migration (SaaS)** -> Mendesain rancangan field `merchant_id` dan `outlet_id` di schema database (PostgreSQL/MySQL).
- [ ] **Tahap 2: Logika Laba Rugi & HPP (Poin 1)** -> Merumuskan rumus dan pencatatan bahan baku (Cost of Goods Sold/HPP) agar laporan laba rugi akurat.
- [ ] **Tahap 3: Export CSV (Poin 2)** -> Merancang formatting output agar ramah dibaca (human-readable) ketika di-download oleh pemilik UMKM.
