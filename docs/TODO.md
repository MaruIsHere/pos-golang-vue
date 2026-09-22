# TODO — Login Page & Autentikasi (`pos-golang-vue`)

Status: OPEN (untuk dikerjakan rekan — MaruIsHere)
Tanggal: 22 Sep 2026
Konteks: backend saat ini **tanpa autentikasi** (tidak ada tabel `User`, tidak ada middleware, tidak ada halaman login). Semua endpoint `/api/*` terbuka, termasuk `POST /orders/:id/refund`, `PUT /settings`, dan `POST /settings/switch-db`.

Acuan kode: `backend/models/models.go`, `backend/handlers/*.go`, `backend/routes/routes.go`, `backend/database/database.go`, `frontend/src/App.vue`, `docs/openapi.yaml`.

---

## Keputusan yang harus dikunci dulu (tanya Adi sebelum coding)

1. Satu device dipakai bergantian (butuh PIN/ganti kasir cepat) atau tiap kasir punya device sendiri?
2. Perlu role berbeda (`admin` vs `kasir`, mis. hanya admin boleh refund + ubah settings) atau semua boleh semua?
3. Kredensial berupa PIN saja atau username + password?

---

## 1. Tabel `users` + hash kredensial [backend]

- Tambah struct `User` di `backend/models/models.go`: `id`, `name`, `username` (unik, opsional jika PIN saja), `pin_hash`/`password_hash` (bcrypt, **jangan simpan plaintext**), `role` (`admin`/`kasir`), `is_active`, `created_at`.
- Daftarkan ke `AutoMigrate` di `backend/database/database.go` (sejajar 8 tabel lain).
- Kriteria selesai: `go run .` membuat tabel `users`; `go vet ./...` bersih.

## 2. Endpoint auth + middleware JWT [backend]

- File baru (saran): `backend/handlers/auth.go` + pasang middleware di `backend/routes/routes.go`.
- Endpoint: `POST /api/auth/login` (terima kredensial → kembalikan token + info user), `POST /api/auth/logout` (invalidasi sisi client / blacklist sederhana bila perlu), `GET /api/auth/me` (validasi token → profil).
- Middleware Gin: tolak tanpa `Authorization: Bearer <token>` valid (401), kecuali `GET /api/health`.
- Kriteria selesai: login salah → 401; akses `/api/products` tanpa token → 401; dengan token → 200.

## 3. Proteksi endpoint + kasir dari token [backend]

- `CashierName` pada `CreateOrder` (`backend/handlers/order.go`) diambil dari klaim token, hapus hardcode `"Kasir 1"`.
- Jika ada role: batasi `POST /orders/:id/refund`, `PUT /settings`, `POST /settings/switch-db`, `DELETE /products/:id` hanya untuk `admin` (403 untuk kasir).
- Kriteria selesai: kasir tidak bisa refund (403); admin bisa; nama kasir di struk sesuai user login.

## 4. Halaman login + guard [frontend]

- File baru (saran): `frontend/src/views/LoginView.vue`; daftarkan di `frontend/src/App.vue` (tab `login` bila belum ada token).
- Guard: belum ada token di `localStorage` → tampilkan login, sembunyikan 7 tab lain.
- Semua `fetch('/api/...')` di `views/` + `components/` menyertakan header `Authorization: Bearer <token>`; respons 401 → kembali ke login.
- Kriteria selesai: refresh browser tanpa token → login; login sukses → masuk kasir; token kedaluwarsa → tendangan otomatis ke login.

## 5. Seed admin + dokumentasi [backend + docs]

- Seed 1 user admin awal di `SeedInitialData` (`backend/database/database.go`) bila tabel kosong; tulis kredensial awal di chat internal (jangan commit password produksi).
- Tambah skema `User`, `LoginInput`, dan path `/auth/*` ke `docs/openapi.yaml`; validasi dengan `npx @redocly/cli lint docs/openapi.yaml` (harus 0 error).
- Kriteria selesai: fresh `pos.db` langsung bisa login admin; OpenAPI valid.

---

## Perintah verifikasi (jalankan dari root repo)

```bash
cd backend && gofmt -l handlers/ routes/   # harus kosong
cd backend && go vet ./...                 # harus bersih
cd backend && go run .                     # + curl login, akses tanpa token, refund sebagai kasir
npx @redocly/cli lint docs/openapi.yaml    # harus valid
```
