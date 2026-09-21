# Kerangka Acuan Proyek — Project Reference Framework

Status: TEMPLATE / ACUAN
Tanggal: 21 Sep 2026
Tujuan: Agar proyek tidak melenceng dari desain yang sudah disepakati.

---

## 1. PRINSIP: "Satu Acuan, Semua Cek"

Setiap perubahan (fitur baru, perbaikan kode, dokumen) **harus dicek** terhadap acuan ini sebelum di-commit. Ini bukan "aturan kaku", tapi **checklist wajib** agar rekan (`MaruIsHere`, tim, atau kontributor) tidak membuat asumsi berbeda.

---

## 2. STRUKTUR ACUAN (Hierarki)

```
ACUAN PROYEK (File ini)
├── DESIGN.md          → Visi, fitur utama, batasan (scope)
├── ARCHITECTURE.md     → Komponen, aliran data, teknologi
├── MARKETPLACE.md      → Model bisnis, tier, plugin
├── COMPONENTS.md       → Daftar komponen + tanggung jawab
├── API_REFERENCE.md    → Endpoint + kontrak data
├── .gitignore          → File yang TIDAK boleh masuk repo
├── TODO.md             → Daftar tugas (prioritas)
└── DECISIONS.md        → ADR (kenapa pilih A, bukan B)
```

---

## 3. TEMPLATE: DESIGN.md (Contoh — Sudah Ada)

```markdown
# [NAMA PROYEK] — Design Specification
Versi: 1.0 | Tanggal: YYYY-MM-DD
Status: DRAFT | LOCK | IMPLEMENTED

## 1. Visi (1 kalimat)
"[Proyek ini adalah ... untuk ... dengan ...]"

## 2. Fitur Utama (Maksimal 5 — Prioritas)
- [ ] F1: ... (Kritik — Tanpa ini, proyek tidak jalan)
- [ ] F2: ... (Penting — Nilai tambah utama)
- [ ] F3-F5: ... (Tambahan — Bisa dipotong jika waktu tidak cukup)

## 3. Batasan (Scope — Apa yang TIDAK dilakukan)
- TIDAK: ...
- TIDAK: ...

## 4. Referensi Teknologi
- Backend: [Go / Node / ...]
- Frontend: [Vue 3 / React / ...]
- DB: [...]
- Integrasi: [...]
```

---

## 4. TEMPLATE: ARCHITECTURE.md

```markdown
# Arsitektur — [Nama]
Status: LOCK (tidak boleh berubah tanpa ADR)

## A. Diagram (ASCII / Link ke SVG)
## B. Aliran Data (Data Flow)
## C. Komponen & Tanggung Jawab
## D. Multi-Tenant / Isolasi (Jika ada)
## E. Offline / Sync (Jika ada)
```

---

## 5. TEMPLATE: MARKETPLACE.md (Jika SaaS)

```markdown
# Marketplace / SaaS Design
## 1. Tier & Pricing
## 2. Onboarding Flow
## 3. Plugin / Feature System
## 4. White-Label
## 5. Billing & Subscription
```

---

## 6. CHECKLIST: Sebelum Commit (Wajib)

```markdown
[ ] 1. Apakah perubahan sesuai DESIGN.md?
[ ] 2. Apakah ada komponen baru yang belum di COMPONENTS.md?
[ ] 3. Apakah .gitignore sudah melindungi file sensitif?
[ ] 4. Apakah ada file > 1MB yang masuk tracking?
[ ] 5. Apakah fitur ini sesuai tier/marketplace yang disepakati?
[ ] 6. Apakah dokumentasi (READ / API) sudah diperbarui?
```

---

## 7. ACUAN UNTUK REKAN (`MaruIsHere` / Kontributor)

Jika kamu (`MaruIsHere`) atau kontributor lain ingin menambah fitur:

1. **Baca file ini dulu** (`REFERENCE.md` / `01-marketplace-architecture.md`)
2. **Cek checklist** (bagian 6) — semua harus ✅
3. **Jika fitur baru besar** → buat `DECISIONS.md` (ADR) dulu
4. **Jika fitur tidak di DESIGN.md** → diskusi dulu (tidak langsung kode)

---

## 8. REFERENSI FILE DI PROYEK INI

| File Acuan | Lokasi | Status |
|------------|--------|--------|
| Design (Game) | `pulau-segara-docs/00_SYSTEM_CONTEXT.md` | ✅ LOCK |
| Komponen POS | `pos-golang-vue/docs/marketplace-design/02-components.md` | ✅ DRAFT |
| Arsitektur POS | `pos-golang-vue/docs/marketplace-design/01-marketplace-architecture.md` | ✅ DRAFT |
| ADR (Keputusan) | `pos-golang-vue/docs/decisions/adr-index.md` | ✅ 10 ADR |
| API Reference | `pos-golang-vue/docs/api/reference.md` | ✅ Lengkap |

---

*File ini (`REFERENCE.md`) dibuat sebagai **kerangka acuan** agar setiap perubahan di proyek ini (`pos-golang-vue`) atau proyek lain (`pulau-segara`) bisa dicek konsistensinya. Update file ini jika desain berubah — tapi perubahan besar butuh diskusi dulu.*
