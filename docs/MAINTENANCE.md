# MAINTENANCE — Ritme & Prosedur Update Dependensi (`pos-golang-vue`)

Status: SOP TIM (berlaku untuk Adi + rekan)
Tanggal: 22 Sep 2026
Tujuan: project long-running ini tetap up-to-date tanpa breaking surprise.

Prinsip: **update kecil sering, update besar terjadwal. Satu update = satu commit + satu verifikasi.**

---

## 1. Ritme (tempel di kalender)

| Kapan | Apa | Perintah |
|---|---|---|
| Tiap selesai fitur | `gofmt`, `vet`, smoke curl | kebiasaan berjalan, jangan dilepas |
| Bulanan (±30 mnt) | Patch + minor Go & npm, audit | lihat §2 dan §3 |
| Kuartalan (±2–4 jam) | Evaluasi 1 major tertunda | prosedur §4, branch khusus |
| Tiap rilis Go (±6 bulan) | Naik toolchain + baca release notes | go.dev/blog, ubah angka di `backend/go.mod` |

---

## 2. Backend Go — prosedur update

```bash
cd backend

# INTIP dulu (read-only, selalu aman):
go list -m -u all          # baris dengan [x.y.z] = update tersedia

# PATCH saja (risiko nol), mis. cors 1.7.8 -> 1.7.9:
go get -u=patch ./...
go mod tidy

# Verifikasi wajib tiap update:
go vet ./...
go build -o /tmp/cek-backend .
go run .                   # + curl http://localhost:8080/api/health
```

Major (mis. gin v1 → v2 suatu hari): baca release notes GitHub dulu, kerjakan di branch `chore/update-<nama>-v<major>`, jangan di `main`. Yang biasanya pecah = API library, bukan bahasa Go (jaminan kompatibilitas Go 1).

Naik toolchain Go: upgrade via package manager OS (`go version` harus ≥ angka di `go.mod`), lalu sesuaikan baris `go X.YY` di `backend/go.mod` + `go mod tidy`.

---

## 3. Frontend npm — prosedur update

```bash
cd frontend

# INTIP (read-only): kolom Current / Wanted / Latest
npm outdated
# Wanted  = aman auto-update (masih dalam ^ di package.json)
# Latest  = keputusan manual (beda major)

# AMAN saja:
npm update                 # hanya naik ke Wanted

# MAJOR satu per satu, JANGAN sekaligus:
npm install vite@8         # lalu verifikasi §5 sebelum lanjut ke paket berikut
```

`npm audit` tiap bulan; `npm audit fix` hanya untuk non-breaking (cek diff-nya).

---

## 4. Prosedur migrasi MAJOR (checklist)

1. `git status` bersih → buat branch `chore/update-<paket>-v<major>`.
2. Baca changelog major: breaking changes + codemod yang disediakan.
3. Install 1 paket → `npm run build` (atau `go build`) → catat error.
4. Perbaiki kode yang pecah (mis. icon rename, API berubah) → build hijau.
5. Smoke test: `npm run dev` / `go run .`, buka semua halaman, cek console browser + log backend.
6. Commit di branch → merge ke `main` hanya jika langkah 4–5 hijau.
7. Contoh riil di repo ini: migrasi Vite 5 → 8 + plugin-vue 5 → 6 + lucide 0.344 → 1.0 (Sep 2026), tanpa ubah kode Vue.

---

## 5. Verifikasi standar tiap update (definisi "selesai")

```bash
cd backend && gofmt -l handlers/ routes/   # harus kosong
cd backend && go vet ./...                 # harus bersih
cd backend && go run .                     # curl /api/health + 1 endpoint tulis
cd frontend && npm run build               # harus sukses
cd frontend && npm run dev                 # buka http://localhost:3000, cek console
npx @redocly/cli lint docs/openapi.yaml    # bila API berubah
```

---

## 6. Otomatisasi (nanti, saat tim tumbuh)

Dependabot (PR update otomatis mingguan) + CI (`go vet`, `go build`, `npm run build` tiap PR). Butuh file `.github/dependabot.yml` + workflow Actions — belum dipasang (lihat TODO).
