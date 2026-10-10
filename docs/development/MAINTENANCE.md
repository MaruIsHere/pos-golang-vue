# Maintenance dan Rilis

Panduan ringkas untuk menjaga perubahan kecil, terverifikasi, dan mudah dipulihkan. Perintah berikut mengikuti struktur repositori saat ini.

## Verifikasi perubahan

```bash
# Frontend
npm --prefix frontend run typecheck
npm --prefix frontend run build

# Backend
cd backend
gofmt -w ./cmd ./internal
go test ./...
go vet ./...
go build ./cmd/api
```

Jalankan `gofmt -w` hanya pada file Go yang memang berubah; perintah di atas ditampilkan sebagai pola, bukan instruksi untuk menulis ulang semua file saat review. Untuk pengecekan format tanpa mengubah file:

```bash
gofmt -l ./cmd ./internal
```

## Perubahan database

- Model GORM dimigrasikan melalui `AutoMigrate` saat startup database.
- Tinjau perubahan model pada SQLite dan MySQL jika keduanya didukung.
- Backup sebelum migrasi, penggantian database, atau operasi yang memodifikasi data.
- Uji pemulihan backup pada database terpisah sebelum mengandalkannya.

## Perubahan API

1. Ubah route, handler, validasi, dan model sesuai kebutuhan.
2. Perbarui `docs/api/openapi.yaml` jika spesifikasinya mencakup endpoint tersebut.
3. Bandingkan spesifikasi dengan `backend/internal/routes/routes.go`; OpenAPI saat ini belum dijamin lengkap.
4. Jalankan test backend dan build frontend bila kontrak UI berubah.

## Perubahan frontend

- Pastikan state loading, error, data kosong, dan sukses memiliki perilaku jelas.
- Uji role dan konteks toko yang relevan.
- Jalankan typecheck dan build produksi.
- Untuk perubahan checkout, verifikasi jumlah item, subtotal, diskon, pajak, total, pembayaran, stok, dan struk sebagai satu alur.

## Build dan menjalankan hasil build

```bash
npm run build
npm start
```

Backend melayani aset `frontend/dist`; karena itu UI deployment lokal berubah setelah frontend dibuild ulang. Perintah `make build-linux`, `make build-win`, dan `make build-mac` hanya membangun binary backend untuk target masing-masing; jalankan build frontend secara terpisah untuk menghasilkan aset web.

## Keamanan deployment

- Ganti akun seed dan simpan konfigurasi/DSN rahasia di luar Git.
- Sajikan aplikasi melalui HTTPS di jaringan yang tidak sepenuhnya dipercaya.
- Batasi CORS pada origin aplikasi yang sah.
- Lindungi file database, upload, backup, dan export laporan.
- Tinjau role dan isolasi merchant/outlet di setiap handler sebelum membuka layanan multi-tenant ke publik.

Dokumen standar produksi aspiratif: [`deployment/production-standard.md`](deployment/production-standard.md). Rekomendasi di sana harus diverifikasi terhadap implementasi.
