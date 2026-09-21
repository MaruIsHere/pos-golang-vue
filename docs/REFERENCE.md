# Checklist Kebutuhan Dev Profesional — Belajar Go + Industri

Status: ACUAN BELAJAR / MENTORING
Tanggal: 21 Sep 2026
Target: MaruIsHere (baru lulus — belajar Go + standar industri)
Proyek Acuan: `pos-golang-vue`
Mentor: satriyadi (pengguna utama proyek ini)

---

## 0. FILOSOFI: Dari "Vibe" ke "Profesional"

Tidak cukup kode jalan. Kode profesional harus:
- **Bisa dibaca** oleh tim lain (bukan hanya penulis)
- **Bisa diuji** (test otomatis, bukan "di laptop jalan")
- **Bisa dideploy** ulang (bukan manual copy-paste)
- **Bisa dipahami** alasan desainnya (bukan "saya coba-coba sampai jalan")

> **Pelajaran utama:** "Kode yang bagus bukan kode yang cepat ditulis — tapi kode yang cepat dipahami 6 bulan kemudian."

---

## 1. GO LANG — Dasar yang Wajib (Referensi Proyek Ini)

### 1.1 Bahasa (Syntax & Konsep)
| Topik | Status (Centang) | Bukti di Proyek |
|-------|------------------|-----------------|
| [ ] Package (`package main`) | ⬜ | `backend/main.go` |
| [ ] Import (`github.com/gin-gonic/gin`) | ⬜ | `backend/main.go` baris 11 |
| [ ] Function (`func main()`) | ⬜ | `backend/main.go` baris 15 |
| [ ] Variable (`cfg := config.LoadConfig()`) | ⬜ | `backend/main.go` baris 16 |
| [ ] Struct (melihat `models/`) | ⬜ | `models/models.go` |
| [ ] Pointer vs Value (`*database.InitDB`) | ⬜ | `database/` |
| [ ] Interface (`database.InitDB`) | ⬜ | `database/database.go` |
| [ ] Go Modules (`go.mod`, `go.sum`) | ⬜ | `backend/` |
| [ ] Error Handling (`if err != nil`) | ⬜ | `main.go` baris 22 |

**Tugas Praktis (1-2 jam):**
1. Buka `backend/main.go` — jelaskan apa setiap baris lakukan (tulis komentar dalam bahasa Indonesia di file `BELAJAR.md`)
2. Ubah `port` di `main.go` — coba jalankan (`go run main.go`)
3. Hapus satu import (`cors`) — lihat error apa yang muncul → perbaiki

---

### 1.2 Framework Go (Proyek Ini: Gin)
| Topik | Bukti di Proyek | Tugas |
|-------|-----------------|-------|
| Router (`router := gin.Default()`) | `main.go` baris 26 | Tambah 1 endpoint baru (`GET /ping`) |
| Middleware (`router.Use(...)`) | `main.go` baris 29 | Tambah middleware log waktu request |
| Handler (`handlers/handlers.go`) | 592 baris — besar! | Baca 50 baris pertama, jelaskan apa yang dilakukan |
| Config (`config.LoadConfig()`) | `config/` | Ubah port di `.env` / config, jalankan ulang |

**Tugas Praktis:**
- Tambahkan endpoint `GET /health` di `handlers/` yang return `{"status":"ok"}`
- Tambahkan middleware yang print waktu setiap request

---

### 1.3 Database & ORM (Proyek Ini: GORM / SQLite / MySQL)
| Topik | Bukti | Tugas |
|-------|-------|-------|
| Init DB (`database.InitDB`) | `database/` | Lihat `database.go` — jelaskan `InitDB` |
| Query (melihat `models/`) | `models/models.go` | Buat query baru: `SELECT * FROM products WHERE stock > 10` |
| Migration | Tidak ada versi! | Buat `migration/` folder dengan versi `v1` |
| Backup / Restore | Tidak terdokumentasi | Buat script `scripts/backup-db.sh` |

**Tugas Praktis:**
- Buat file `backend/scripts/backup-db.sh` (bash) yang `cp backend/*.db backup/`
- Buat `backend/scripts/test-query.go` yang query semua produk

---

## 2. PROFESIONAL DEV — Standar Industri (Referensi `docs/deployment/production-standard.md`)

### 2.1 Kode Bersih (Clean Code)
| Prinsip | Contoh Buruk (Vibe) | Contoh Baik (Profesional) |
|---------|---------------------|---------------------------|
| Nama variabel | `x`, `a`, `tmp` | `customerName`, `orderTotal` |
| Fungsi panjang | 500 baris | Maksimal 30-50 baris, 1 tugas |
| Komentar | Tidak ada / "ini buat ini" | Jelaskan **kenapa**, bukan **apa** |
| Magic number | `if count > 7` | `if count > MAX_ITEMS_PER_ORDER` |

**Tugas:** Buka `backend/handlers/handlers.go`. Cari 3 nama variabel/fungsi yang bisa diperbaiki. Ganti dengan nama lebih deskriptif.

---

### 2.2 Testing (Wajib — Tidak Boleh Kosong)
| Tipe | Status Proyek | Tugas |
|------|---------------|-------|
| Unit Test (Go) | ❌ Tidak ada | Buat `backend/handlers/handlers_test.go` — test 1 endpoint |
| Integration Test | ❌ Tidak ada | Buat `tests/integration/` — test API dengan `curl` |
| E2E (Frontend) | ❌ Tidak ada | Buat `frontend/src/components/CartDrawer.spec.vue` (Vue Test Utils) |

**Referensi:** `docs/deployment/production-standard.md` — bagian CI/CD (baris 49-87) menjelaskan pipeline test yang wajib.

---

### 2.3 Version Control (Git)
| Prinsip | Status Proyek | Perbaikan |
|---------|---------------|-----------|
| Commit message konvensional | ✅ Ada (`feat:`, `fix:`, `docs:`) | Pertahankan |
| `.gitignore` lengkap | ✅ Baru (`9990ee05`) | Pertahankan |
| Branch protection (`main`) | ❌ Belum terverifikasi | Buat `branch-protection.md` |
| PR Review (Pull Request) | ❌ Belum terdokumentasi | Buat `CONTRIBUTING.md` |

**Tugas:** Buat `CONTRIBUTING.md` (di root repo) yang berisi:
- Cara buat PR
- Aturan commit message
- Checklist sebelum PR (`README.md` bagian 6)

---

### 2.4 Dokumentasi (Tidak Boleh Kosong)
| Dokumen | Status | Tindakan |
|---------|--------|----------|
| `README.md` (proyek) | ✅ Lengkap | Sudah baik |
| `docs/deployment/production-standard.md` | ✅ Lengkap (365 baris) | Sudah baik — ini acuan |
| `docs/marketplace-design/REFERENCE.md` | ✅ Lengkap (123 baris) | Sudah baik |
| `API.md` / `docs/api/reference.md` | ✅ Lengkap (12973 chars) | Sudah baik |
| `CHANGELOG.md` | ❌ Belum ada | Buat — catat setiap versi release |
| `LICENSE` | ❌ Belum ada | Tambahkan MIT License |

**Tugas:** Buat `CHANGELOG.md` (format: `Keep a Changelog`) dan `LICENSE`.

---

### 2.5 Security (Wajib — Referensi `production-standard.md` Baris 270)
| Prinsip | Status Kode | Perlu Diperiksa |
|---------|-------------|----------------|
| Helmet (`security headers`) | Tidak terlihat di `main.go` | ✅ Perlu ditambah |
| Rate Limit | Tidak terlihat | ✅ Perlu ditambah (`docs/deployment/production-standard.md` baris 91) |
| Input Validation (`zod` / `joi`) | Belum diverifikasi | Periksa `handlers/` |
| Audit Log (operasi sensitif) | Belum diverifikasi | Periksa kode |
| JWT Rotation (refresh token reuse) | Belum diverifikasi | Periksa `auth/` (belum ada) |

---

### 2.6 Database & Migration
| Prinsip | Status | Perbaikan |
|---------|--------|-----------|
| Schema Versioned (`prisma/migrations/`) | ❌ Tidak ada Prisma | Buat folder `migrations/` atau pakai `sql-migrate` |
| Backup Otomatis (`scripts/backup-db.sh`) | ❌ Ada dokumen tapi belum diuji | Uji script — jalankan satu kali |
| Restore Uji (`scripts/restore-db.sh`) | ❌ Belum ada | Buat script |

---

### 2.7 Monitoring & Observability
| Prinsip | Status | Tindakan |
|---------|--------|----------|
| Health Check (`/health`) | ✅ Ada (`docs/deployment/`) | Implementasi kode belum diverifikasi |
| Metrics (`/metrics`) | ❌ Tidak terlihat | Tambah endpoint Prometheus |
| Alert (`Slack` / `Email`) | ❌ Hanya dokumen | Buat `.env` template + script alert |

---

## 3. RENCANA BELAJAR (`MaruIsHere` — 8 Minggu)

### Minggu 1-2: Go Dasar + Kode Proyek Ini
| Hari | Tugas | Bukti |
|------|-------|-------|
| 1-2 | Baca `main.go`, `models/`, `database/` — tulis komentar | File `BELAJAR.md` |
| 3-4 | Buat endpoint baru (`/ping`) + test (`go test`) | Commit `test: add ping endpoint` |
| 5-7 | Perbaiki 3 nama variabel/fungsi | Commit `refactor: improve naming` |

### Minggu 3-4: Testing + Security
| Hari | Tugas | Bukti |
|------|-------|-------|
| 1-3 | Buat `handlers/handlers_test.go` — 1 test | Commit `test: add handler test` |
| 4-5 | Tambah `helmet` middleware | Commit `security: add helmet` |
| 6-7 | Buat `scripts/backup-db.sh` + uji | Commit `chore: add backup script` |

### Minggu 5-6: Dokumentasi + Deployment
| Hari | Tugas | Bukti |
|------|-------|-------|
| 1-2 | Buat `CHANGELOG.md` + `LICENSE` | Commit `docs: add changelog and license` |
| 3-4 | Buat `CONTRIBUTING.md` | Commit `docs: add contributing` |
| 5-7 | Setup CI pipeline (`.github/workflows/deploy.yml`) | PR `feat: add CI pipeline` |

### Minggu 7-8: Marketplace + Plugin
| Hari | Tugas | Bukti |
|------|-------|-------|
| 1-3 | Buat `Tenant` model (`prisma` atau Go struct) | Commit `feat: add tenant model` |
| 4-5 | Buat feature toggle (`feature-flags` sederhana) | Commit `feat: add feature flags` |
| 6-7 | Buat `PluginInstance` + 1 plugin (`inventory`) | Commit `feat: add plugin system` |

---

## 4. METRIK KESUKSESAN (Bagaimana Tahu "Profesional"?)

```markdown
[ ] Kode bisa dibaca tanpa penjelasan oral dari penulis
[ ] Semua endpoint punya test (unit atau integration)
[ ] Semua operasi sensitif punya audit log
[ ] Semua secret tidak di-commit (`.env` tidak di repo)
[ ] CI pipeline jalan otomatis (build, test, security scan)
[ ] Deploy bisa rollback dalam < 3 menit
[ ] Dokumentasi (README, API, ChangeLog) lengkap dan up-to-date
[ ] .gitignore melindungi semua file sensitif (verified dengan `git ls-files`)
```

---

## 5. REFERENSI UNTUK BELAJAR LEBIH LANJUT

| Topik | Sumber | Keterangan |
|-------|--------|------------|
| Go Fundamentals | `https://go.dev/doc/effective_go` | Dokumen resmi — wajib baca |
| Clean Code (Go) | `https://github.com/golang-standards/project-layout` | Struktur proyek standar |
| Testing Go | `https://pkg.go.dev/testing` | Package testing |
| Security (OWASP) | `https://owasp.org/www-project-top-ten/` | Top 10 risiko web |
| CI/CD (GitHub) | `docs/deployment/production-standard.md` baris 49-87 | Contoh pipeline lengkap |
| Multi-Tenant | `docs/marketplace-design/01-marketplace-architecture.md` | Arsitektur proyek |
| Database Migration | `https://www.postgresql.org/docs/current/` | Backup/Restore |

---

## 6. CATATAN MENTOR (`satriyadi` → `MaruIsHere`)

> "Saya (`satriyadi`) membuat dokumen ini bukan untuk menghakimi — tapi sebagai **peta jalan**. Kode yang ada (`pos-golang-vue`) sudah nyata (`main.go`, `handlers/`, `database/`), tapi belum melalui semua standar produksi. Tugasmu (`MaruIsHere`) bukan "menyempurnakan semua sekaligus" — tapi **menambah satu layer profesional setiap minggu**.
>
> **Minggu ini:** Pilih 1 tugas dari Minggu 1 (Go Dasar) → kerjakan → commit → push. Itu sudah langkah besar dari "vibe coding" ke "dev profesional"."

---

*File ini (`REFERENCE.md`) dibuat sebagai kerangka acuan proyek POS (`pos-golang-vue`) — bukan untuk proyek `pulau-segara` (game). Dokumen game (`pulau-segara-docs/`) hidup terpisah di `~/Projects/gamedev/pulau-segara/`. Jangan mencampur standar/aturan antar proyek. Update file ini jika standar POS berubah — perubahan besar butuh diskusi dulu (lihat `docs/deployment/production-standard.md`).
