# Laporan Audit Frontend

**Tanggal:** 11 Oktober 2026  
**Cakupan:** Pemeriksaan statis seluruh kode di `frontend/src` (halaman, komponen, store, composable, utilitas, dan integrasi API), ditambah typecheck dan build produksi.  
**Mode:** Read-only; audit tidak mengubah kode aplikasi.

## Ringkasan

Typecheck dan build frontend berhasil. Ditemukan dua masalah pada alur transaksi yang dapat memengaruhi penerapan diskon dan voucher. Build juga memperingatkan ukuran bundle halaman laporan yang besar.

## Temuan

### P2 — Diskon member tidak diteruskan sebagai diskon pesanan

**Lokasi:** [`PaymentModal.vue` (baris 527)](../../frontend/src/components/PaymentModal.vue#L527), [`PaymentModal.vue` (baris 612)](../../frontend/src/components/PaymentModal.vue#L612), [`PosView.vue` (baris 475)](../../frontend/src/views/PosView.vue#L475)

Modal pembayaran menghitung dan menampilkan potongan member, lalu mengirim `paid_amount` berdasarkan total setelah potongan. Namun payload pesanan dari POS tetap mengirim nilai diskon yang sebelumnya dihitung di keranjang dan tidak memasukkan diskon member. Backend menerima nilai pembayaran lebih rendah daripada total pesanan yang tercatat; akibatnya nilai total dan kembalian tidak sesuai dengan harga yang ditampilkan kepada kasir.

### P2 — Kode voucher yang diketik tidak dapat diterapkan

**Lokasi:** [`CartDrawer.vue` (baris 91)](../../frontend/src/components/CartDrawer.vue#L91), [`CartDrawer.vue` (baris 231)](../../frontend/src/components/CartDrawer.vue#L231)

Input voucher terikat ke `voucherCode`, sedangkan fungsi penerapan membaca `voucherInput`. Selain itu, tidak ada tombol terapkan atau handler Enter pada input tersebut. Karena itu, kode yang diketik manual tidak pernah diproses; voucher hanya bisa dipilih melalui tombol saran.

## Pemeriksaan

- `npm --prefix frontend run typecheck` — berhasil.
- `npm --prefix frontend run build` — berhasil, termasuk pembuatan service worker PWA.
- Build memberi peringatan bahwa chunk laporan berukuran sekitar **1,5 MB** (sekitar **526 KB gzip**). Ini berpotensi memperlambat pemuatan fitur laporan, terutama pada koneksi lambat.

## Batasan

Pemeriksaan alur dilakukan dengan membaca kode dan memeriksa kompilasi/build. Aplikasi tidak dijalankan melalui skenario transaksi end-to-end di browser, sehingga perilaku integrasi runtime belum diverifikasi.
