# POS Golang + Vue 3 - Project Status & TODO

## ✅ Selesai (Fase 1 & 2: Setup & Refactoring UI Admin)
- [x] Migrasi sistem *styling* ke standar Tailwind CSS v4 murni.
- [x] Pembangunan UI Kit (*AppButton, AppInput, AppCard, AppBadge*).
- [x] Refactoring `ProductsView.vue` (*Split-Screen* presisi & bebas *error Invalid end tag*).
- [x] Refactoring `SettingsView.vue` (Desain kartu M3, tipografi elegan, lencana interaktif).
- [x] Refactoring `OrdersView.vue` (Tabel modern dengan sudut melengkung 24px).
- [x] Refactoring `ReportsView.vue` (Distandardisasi menggunakan *AppCard*).
- [x] Pembersihan total `glass-panel`, `.btn`, dan *legacy CSS* lainnya dari seluruh komponen.
- [x] Perbaikan *bug Overflow* Barcode Scanner (Menggunakan kontrol segmen ala iOS).

## 🚀 Selanjutnya (Fase 3: Autentikasi & Sinkronisasi)
- [ ] Desain & Implementasi `LoginView.vue` (Opsi UI: Form PIN Numpad Kasir).
- [ ] Desain & Implementasi `RegisterView.vue`.
- [ ] Setup Pinia Store untuk Manajemen Sesi Kasir (JWT Token).
- [ ] Pembuatan Endpoint Backend Golang (Login, Register, Autorisasi Role Kasir/Admin).
- [ ] Arsitektur *Offline-First* & *Background Sync Worker* (Sinkronisasi SQLite lokal ke MySQL jarak jauh saat internet tersedia).
