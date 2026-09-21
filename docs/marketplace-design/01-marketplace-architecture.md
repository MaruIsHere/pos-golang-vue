# Desain Sistem POS Marketplace — Arsitektur & Komponen

Status: DRAFT (untuk review rekan — MaruIsHere)
Tanggal: 21 Sep 2026
Proyek: pos-golang-vue (Go + Vue 3 PWA)
Target: Dijual di marketplace (SaaS / License)

---

## 1. VISI: POS Yang Bisa Dijual

Bukan sekadar aplikasi kasir — ini **platform** yang dibeli oleh pemilik usaha kecil/menengah (UMKM) yang tidak punya tim IT. Fitur bisa disesuaikan, harga bisa di-tierkan, dan pembeli bisa mendaftar sendiri.

---

## 2. MASALAH UTAMA: "Kebutuhan Random"

Karena pembeli berbeda-beda (warung, cafe, toko pakaian, klinik kecil), kita tidak bisa memprediksi semua fiturnya. Solusi: **modular + feature flags + plugin**.

---

## 3. ARSITEKTUR MARKETPLACE

```
┌─────────────────────────────────────────────────────────────┐
│                    MARKETPLACE LAYER                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │  Onboarding │  │  Billing    │  │  Admin      │          │
│  │  Self-Reg   │  │  Subskripsi │  │  Marketplace│          │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘          │
└─────────┼────────────────┼────────────────┼───────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────┐
│                    CORE POS SERVICE                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │ Transaksi   │  │ Multi-Tenant│  │ Plugin/     │          │
│  │ (Order/Pay) │  │ (RLS/DB)    │  │ Feature     │          │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘          │
└─────────┼────────────────┼────────────────┼───────────────────┘
          │                │                │
          ▼                ▼                ▼
┌─────────────────────────────────────────────────────────────┐
│                    CLIENT (PWA Vue 3)                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │ POS Kasir   │  │ Dashboard   │  │ White-Label │          │
│  │ (Offline)   │  │ Owner       │  │ Custom URL  │          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. KOMPONEN DESAIN (Komponen Wajib Ada)

### 4.1. Multi-Tenant & Isolasi
| Komponen | Implementasi Saat Ini | Perlu Ditambah |
|----------|----------------------|----------------|
| **DB RLS** | Prisma middleware (`outletId`) | ✅ Sudah ada |
| **Auth JWT** | Per aktif (Rani/Bima) | ✅ Perlu per outlet (bukan global) |
| **Onboarding** | Wayan (manual) | ❌ **Self-service register** |
| **Billing** | Tidak ada | ❌ **Midtrans/Gopay subscription** |

### 4.2. Plugin / Feature Flag System
```typescript
// Per tenant (setiap pembeli / outlet)
interface TenantFeatures {
  qris: boolean;              // Default: true
  inventory: boolean;         // Default: true
  loyalty: boolean;           // Default: false (premium?)
  whiteLabel: boolean;        // Default: false
  reporting: 'none' | 'daily' | 'weekly';
  language: 'id' | 'en';
  printerType: 'thermal' | 'inkjet';
  accountingIntegration: string[]; // ['xero', 'quickbooks']
}
```

**Karena random:** Pembeli bisa aktifkan/menonaktifkan fitur dari dashboard, bukan dari kode.

### 4.3. White-Label (Branding Pembeli)
- **Dashboard:** Ganti warna, logo, nama aplikasi (`toko-ku-pos` vs `POS Pro`)
- **Domain:** `toko-budi.pos.com` atau custom domain (CNAME)
- **Struk:** Footer struk bisa diubah teks/logo
- **Email:** Notifikasi dari domain pembeli, bukan `pos-system`

### 4.4. Onboarding Self-Service (Kritik)
Saat ini: `Wayan` (manual). Untuk marketplace:

```
1. Pembeli buka situs marketplace
2. Isi form (nama toko, email, lokasi)
3. Pilih plan (Free / Pro / Enterprise)
4. Bayar / Trial 7 hari
5. System buat outlet baru + admin user
6. Kirim email welcome + link login
7. Redirect ke PWA (POS) dengan data awal (produk kosong / template)
```

---

## 5. TIER / MODEL BISNIS (Usulan)

| Plan | Harga/Bulan | Fitur | Batas |
|------|-------------|-------|-------|
| **Starter** | Grace / Gratis | Order, Cash, QRIS dasar, 1 outlet, 1 kasir | 500 transaksi/bln |
| **Pro** | Rp149.000 / bln | + Inventory, Reports harian, 3 outlet, 3 kasir, White-label | 5.000 transaksi/bln |
| **Enterprise** | Rp499.000 / bln | + Loyalty, Accounting API, Custom domain, Support, Unlimited | Tidak terbatas |

---

## 6. KOMponen YANG HARUS DITAMBAH KE PROJECT

| Komponen | File Target | Status -> Tindak |
|----------|-------------|-------------------|
| **Self-Service Register** | `apps/api/src/modules/auth/register-marketplace.ts` | **Baru** |
| **Subscription/Billing** | `apps/api/src/modules/billing/` | **Baru** (Midtrans/Stripe Billing) |
| **Feature Toggle** | `apps/api/src/common/features/feature-flags.ts` | **Baru** |
| **Plugin System** | `apps/api/src/plugins/` | **Baru** (modular) |
| **White-Label Config** | `apps/api/src/modules/tenants/settings.ts` | **Perlu扩展** |
| **Admin Marketplace** | `apps/dashboard/src/views/AdminMarket/` | **Baru** |

---

## 7. RISIKO "RANDOM" & MITIGASI

| Risiko | Mitigasi |
|--------|----------|
| Pembeli minta fitur aneh | Feature flags — bisa on/off tanpa deploy ulang |
| Pembeli butuh integrasi X | Plugin API (webhook + REST) — pembeli bisa hubungkan sendiri |
| Pembeli butuh bahasa lain | `language` flag + i18n file (JSON lokal) |
| Pembeli ingin harga berbeda | Pricing engine: `plan_id` per tenant, bukan hard-code |
| Pembeli ingin print format berbeda | Template sistem (HTML/ESC + custom footer) |

---

## 8. REKOMENDASI LANGKAH BERIKUTNYA

1. **Pilih 1 tier** dulu (Starter/Pro) — jangan buat semua sekaligus
2. **Buat Feature Toggle sederhana** (tabel `tenant_features`) — ini paling penting karena menyelesaikan "random"
3. **Onboarding manual → otomatis** (form + payment) — bisa pakai Midtrans Subscription
4. **White-label** bisa ditunda (Phase 2), tapi struktur data harus siap (`brand_name`, `primary_color`, `logo_url`)

---

*Disusun untuk review oleh MaruIsHere. Dokumentasi ini disimpan di `/home/satriyadi/Projects/pos-golang-vue/docs/marketplace-design/` agar bisa diakses bersama.*
