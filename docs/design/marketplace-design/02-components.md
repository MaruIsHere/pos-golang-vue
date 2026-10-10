# Komponen Desain — Marketplace POS (Untuk Review MaruIsHere)

Status: DRAFT
Tanggal: 21 Sep 2026
Proyek: pos-golang-vue
Referensi desain utama: `docs/marketplace-design/01-marketplace-architecture.md`

---

## 1. KOMPONEN UTAMA (Modular)

### 1.1 Core POS (Selalu Aktif — Tidak Bisa Dinonaktifkan)
```
├── OrderEngine       # State machine: PENDING → PAID → COMPLETED → CANCELLED
├── PaymentHandler    # Cash + QRIS (Midtrans) + Other
├── ReceiptPrinter    # ESC/POS + HTML (template)
├── TaxEngine         # Persentase pajak (per outlet, default dari plan)
└── AuditTrail        # Log semua operasi (security)
```

**Catatan:** Ini harus selalu jalan, bahkan kalau semua plugin dimatikan.

---

### 1.2 Plugin / Feature Modul (Opsional — Per Tenant)

```
plugins/
├── inventory/          # Stok produk, alert stok rendah, import CSV
├── loyalty/            # Poin, diskon, membership
├── reporting-advanced/ # Laporan harian, mingguan, ekspor PDF/Excel
├── accounting/         # Integrasi Xero, QuickBooks (via webhook/API)
├── multi-outlet/       # 1 akun > 1 toko (sudah ada, tapi bisa on/off)
└── white-label/        # Ganti logo, warna, domain custom
```

**Setiap plugin punya:**
- `plugin.json` — metadata (nama, versi, harga tier)
- `src/` — kode Go (handler, service)
- `frontend/src/` — komponen Vue (jika perlu UI tambahan)
- `migrations/` — skema DB tambahan (jika perlu tabel baru)

---

### 1.3 Marketplace Layer (Baru — Untuk Penjualan)

```
marketplace/
├── onboarding/         # Form self-service, trial 7 hari
├── billing/             # Midtrans Subscription / Stripe Billing
├── tenant-manager/      # Admin marketplace (daftar semua tenant)
├── feature-flags/       # Toggle fitur per tenant
├── domain-router/       # Custom domain (CNAME) per tenant
└── analytics/           # Metrik penggunaan (transaksi/bln, active users)
```

---

## 2. DATA MODEL EKSTENSI (Multi-Tenant + Plugin)

### 2.1 Tabel Baru: `Tenant`
```
Tenant (baru — per pembeli / akun marketplace)
- id (cuid)
- name (string) — nama bisnis
- ownerId (user)
- plan (enum: free, pro, enterprise)
- status (enum: trial, active, suspended, cancelled)
- features (JSONB) — feature flags
- branding (JSONB) — logo_url, primary_color, custom_domain
- billingStatus (enum: paid, overdue, trial)
- subscriptionId (string) — referensi Midtrans/Stripe
```

### 2.2 Tabel Baru: `PluginInstance`
```
PluginInstance (baru — plugin aktif per tenant)
- id (cuid)
- tenantId
- pluginName (string) — 'inventory', 'loyalty', 'reporting'
- isActive (boolean) — toggle
- config (JSONB) — konfigurasi plugin (misal: diskon 10%)
- version (string)
```

### 2.3 Tabel Baru: `Subscription`
```
Subscription (baru — billing marketplace)
- id (cuid)
- tenantId
- planName (string)
- amount (int — dalam rupiah)
- billingCycle (enum: monthly, yearly)
- startDate, endDate
- status (enum: active, cancelled, failed)
- paymentGatewayRef (string) — Midtrans/Stripe
```

---

## 3. KOMPONEN FRONTEND (Vue 3 — Baru / Perlu Perubahan)

### 3.1 PWA POS (Kasir)
```
Saat ini sudah ada:
- OrdersView.vue, CartDrawer.vue, PaymentModal.vue, SettingsView.vue

Perlu ditambah (jika plugin aktif):
- InventoryPanel.vue        # Komponen stok (plugin inventory)
- LoyaltyCard.vue          # Kartu poin pelanggan (plugin loyalty)
- ReportingDashboard.vue   # Grafik penjualan (plugin reporting)
- CustomBrandHeader.vue    # Header dengan logo pembeli (plugin white-label)
```

### 3.2 Dashboard (Owner)
```
Saat ini sudah ada: SettingsView.vue, Reports? (baru di commit second-commit)

Perlu ditambah (marketplace):
- SubscriptionPanel.vue     # Info paket, upgrade/downgrade, invoice
- PluginMarketplace.vue     # Daftar plugin + toggle on/off
- BrandingConfig.vue        # Upload logo, pilih warna, domain
- TrialBanner.vue           # Banner "X hari tersisa trial"
```

---

## 4. FLOW ONBOARDING MARKETPLACE (Visual)

```
[Pengunjung Website Marketplace]
        │
        ▼
[Form: Nama Bisnis, Email, Lokasi, Nama Toko]
        │
        ▼
[Pilih Plan: Starter (Free) / Pro / Enterprise]
        │
        ▼
[Bayar / Mulai Trial 7 Hari] ─── Midtrans / Stripe
        │
        ▼
[Server: Buat Tenant + Outlet + Admin User + Feature Flags Default]
        │
        ▼
[Kir Email: Link Login + Data Awal + Dokumentasi]
        │
        ▼
[Redirect ke PWA] ─── POS Kasir siap pakai (offline-ready)
```

---

## 5. FLOW FEATURE PLUGIN (Contoh: Inventory)

```
1. Admin Marketplace approve plugin 'inventory' untuk Pro/Enterprise
2. Tenant (pembeli) buka Dashboard → Plugin Marketplace
3. Toggle: "Inventory" → ON
4. Server: Buat tabel `plugin_inventory_items`
5. Server: Update `Tenant.features` → inventory: true
6. PWA: Component `InventoryPanel.vue` muncul (karena feature flag aktif)
7. PWA: Data baru (stok) tersimpan di IndexedDB + sync ke server
```

---

## 6. KEAMANAN MARKETPLACE (Tambahan dari Multi-Tenant)

| Ancaman | Mitigasi |
|---------|----------|
| Pembeli A akses data B | RLS PostgreSQL (`outlet_id` + `tenant_id`) + JWT claim (`tenantId`) |
| Plugin berbahaya (malware) | Plugin hanya bisa dibuat oleh admin marketplace; review manual sebelum publish |
| Feature bypass (pembeli edit JS) | Validasi feature flag di server (bukan hanya frontend) — setiap API call cek `Tenant.features` |
| Data leakage via domain custom | Setiap request ke custom domain divalidasi `Host` header → `tenantId` |

---

## 7. CATATAN IMPLEMENTASI (Untuk Developer)

### Prioritas Implementasi (Urutan Yang Logis)
```
Fase 1 (Kritik):
  1. Feature Toggle System (tabel + middleware + API)
  2. Plugin System dasar (struktur folder + loader)
  3. White-Label (logo, warna, domain config per tenant)

Fase 2 (Marketplace):
  4. Self-Service Onboarding (form + trial + billing)
  5. Subscription Management (billing cycle, cancel, retry)
  6. Plugin Marketplace UI (Dashboard → pilih plugin)

Fase 3 (Scale):
  7. Analytics (penggunaan per tenant, conversion rate)
  8. Admin Marketplace Panel (review plugin baru, suspend tenant)
  9. Multi-instance / On-prem option (Docker image untuk pembeli enterprise)
```

---

## 8. REFERENSI FILE (Untuk Rekan MaruIsHere)

| File | Lokasi | Keterangan |
|------|--------|------------|
| Dokumen ini | `docs/marketplace-design/01-marketplace-architecture.md` | Arsitektur lengkap |
| Komponen ini | `docs/marketplace-design/02-components.md` | Dokumen ini |
| Design utama | `docs/marketplace-design/` | Semua dokumen desain |
| ADR Index | `docs/decisions/adr-index.md` | 10 keputusan arsitektur |
| API Reference | `docs/api/reference.md` | Endpoint lengkap (termasuk sync) |

---

*Disusun oleh Hermes Agent (sesuai diskusi dengan satriyadi) — 21 Sep 2026*
*Untuk review dan implementasi oleh MaruIsHere dan tim.*
