# Dokumentasi POS Kasir Pro

Dokumentasi ini dikelompokkan berdasarkan tujuan agar pembaca dapat memilih jalur belajar yang sesuai.

## Mulai dari sini

| Kebutuhan | Dokumen |
|---|---|
| Memahami tujuan produk dan ruang lingkup fitur | [PRD](product/PRD.md) |
| Memahami rancangan pengalaman dan alur antarmuka | [Desain produk dan UI](design/DESIGN.md) |
| Memahami komponen dan aliran data sistem | [System overview](system/SYSTEM_OVERVIEW.md) |
| Memahami target arsitektur multi-tenant | [Roadmap SaaS](system/SAAS_ROADMAP.md) |
| Menjalankan dan menggunakan aplikasi | [Panduan pengguna](guides/USER_GUIDE.md) |
| Meninjau hasil audit frontend | [Audit frontend, 11 Oktober 2026](audits/FRONTEND_AUDIT.md) |

## Dokumentasi berdasarkan fungsi

- `product/` — tujuan, pengguna, kebutuhan, dan batas ruang lingkup.
- `design/` — pola UI/UX dan referensi desain.
- `system/` — arsitektur aplikasi, model data, dan keputusan teknis.
- `guides/` — panduan operasional untuk pengguna.
- `api/` — spesifikasi API OpenAPI.
- `development/` — pemeliharaan, standar kontribusi, dan materi belajar developer.
- `audits/` — laporan pemeriksaan yang memiliki tanggal dan cakupan.
- `archive/` — catatan historis yang tidak boleh dianggap sebagai gambaran implementasi terkini.

Materi developer: [jalur belajar kode](development/LEARNING_GUIDE.md), [pedoman kontribusi](development/GUIDELINE.md), [maintenance](development/MAINTENANCE.md), [standar deployment aspiratif](development/deployment/production-standard.md).

## Status dokumentasi

Dokumen kebutuhan dan rancangan menjelaskan implementasi saat ini serta batas fiturnya. Item yang belum tersedia disebut sebagai rencana, bukan kemampuan aktif. Untuk detail endpoint, cocokkan [OpenAPI](api/openapi.yaml) dengan route aktual di `backend/internal/routes/routes.go`.
