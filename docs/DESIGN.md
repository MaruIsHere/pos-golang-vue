---
version: "alpha"
name: "Strict Rigid POS Dashboard (Tailwind v4, Heroicons, MD3)"
description: "High-density POS dashboard built on strict grid layout constraints, fully optimized for Tailwind CSS v4."
colors:
  primary: "#000000" # Pure Black / Charcoal (Untuk Teks Utama & Border Bersih)
  surface: "#FFFFFF" # Pure White (Canvas Utama / Card)
  surface-variant: "#FFFFFF" # Tetap White untuk konsistensi warna kontras
  outline: "#000000" # Clean Borders menggunakan Pure Black murni
  outline-variant: "#9CA3AF" # Gray-400 untuk Muted/Disabled Border
  accent-color: "#000000" # Single Primary Only (Aksen dominan hitam murni)
  tertiary: "#10B981" # Emerald-500 (Khusus Angka Finansial / Total Belanjaan)
  error: "#EF4444" # Red-500 (Error / Danger / Void Item)

typography:
  font-family: "sans-serif"
  display-large:
    class: "font-mono text-4xl font-bold tracking-tight" # MD3 Display untuk Angka Finansial
  title-medium:
    class: "font-sans text-base font-bold tracking-normal" # font-weight: 700 untuk judul
  body-medium:
    class: "font-sans text-sm font-normal tracking-normal" # font-weight: 400 untuk data tabel
  label-medium:
    class: "font-sans text-xs font-medium uppercase tracking-wider"

# Enforced Technical Specification Constants
tokens:
  max-width: "1200px" # Batas kaku kontainer utama
  gap: "gap-8" # Representasi strict gap: 2rem dalam Tailwind v4
  spacing: "p-8" # Representasi strict spacing: 2rem
  border-radius: "rounded-none" # --border-radius: 0px (Sharp edges)
  shadow: "shadow-none" # --shadow: none (No box-shadow unless necessary)
---

## Overview

Arsitektur UI POS dengan kepadatan tinggi yang dibangun berdasarkan batas teknis yang kaku: tata letak berbasis **CSS Grid dengan celah 2rem**, lebar maksimal **1200px**, sudut tajam tanpa bulat (**0px**), dan bebas dari efek bayangan (**no box-shadow**). Desain ini menyatukan kejelasan struktural Swiss Style dengan semantik status interaksi dari Material Design 3 (MD3).

- Density: 9/10 — Maximum Compact Data
- CSS Engine: Tailwind CSS v4 Core Utility System
- Grid Gap Constraints: Strict 2rem (`gap-8`)
- Global Shadow Rule: Enforced `shadow-none`

## Technical Specifications Implementation

Saat `agy CLI` menyusun komponen, ia harus mematuhi token CSS bawaan ini secara mutlak:

- **Layout Container:** Area POS harus dibungkus dalam kontainer dengan kelas `max-w-[1200px] mx-auto grid grid-cols-12 gap-8 p-8`.
- **Borders & Separation:** Pemisah zona antar komponen menggunakan **Clean Borders** murni (`border border-black` atau `divide-y divide-black`). Jangan gunakan elevasi bayangan blur.
- **Typography Matrix:** Seluruh elemen UI menggunakan `font-sans` dengan variasi berat antara `font-normal` (400) hingga `font-bold` (700). Angka hitungan kasir dikecualikan ke `font-mono` demi kerapian vertikal.

## Material Design 3 State Layer Rules

Setiap komponen tombol aksi wajib merespons input kasir secara instan (Toleransi transisi ≤ 75ms) menggunakan state modifiers Tailwind v4:

- **Hover:** Efek overlay tipis (`hover:bg-black/5` atau `hover:invert` untuk tombol kontras tinggi).
- **Focus:** Navigasi scanner/keyboard wajib memicu outline yang tegas (`focus-visible:outline-4 focus-visible:outline-black focus-visible:outline-offset-2`).
- **Pressed (Active):** Efek taktil instan (`active:scale-[0.98]`).
- **Disabled:** Reduksi visual penuh (`disabled:opacity-35 disabled:pointer-events-none`).

## Heroicons Specifications

- **Ukuran Ikon:** Gunakan `size-5` (20px) untuk tombol aksi mikro di dalam tabel keranjang, dan `size-6` (24px) untuk navigasi menu utama.
- **Aturan Varian:**
  - **Heroicons Outline:** Navigasi samping, ikon pencarian, status bar jaringan.
  - **Heroicons Solid:** Tombol penambah/pengurang jumlah barang (`PlusIcon`, `MinusIcon`) dan tombol hapus (`TrashIcon`).

## Component Blueprints for Tailwind v4

- **Main Grid Wrapper:** `w-dvw h-dvh overflow-hidden bg-white text-black font-sans`
- **POS Tri-Grid Zones (Dalam Max-Width 1200px):**
  1. **Cart & Input Scanner Bar (Kiri/Tengah - `col-span-5 flex flex-col gap-8`):** Daftar pesanan dengan baris data setinggi `h-12 border-b border-black`.
  2. **Product Catalog Showcase (Kanan - `col-span-7 grid grid-cols-3 gap-8`):** Grid kartu produk yang dibatasi oleh ukuran celah strict `2rem`.
- **Cart Data Rows:** Menggunakan kelas `font-mono text-right tabular-nums` untuk seluruh angka nominal harga produk agar sejajar lurus ke bawah.

## Do's and Don'ts

- Jangan gunakan ukuran teks adaptif (fluid text) — Gunakan utilitas standar yang pasti (`text-xs`, `text-sm`, `text-base`, `text-4xl`).
- Jangan gunakan bayangan (shadow/elevation blur) — Elevasi MD3 diratakan menjadi garis struktural murni (`border border-black`) untuk menjaga keterbacaan murni.
- Jangan gunakan sudut membulat lebar — Terapkan `rounded-none` secara massal untuk memaksimalkan area hit-box klik/sentuh jari kasir.
- Jangan gunakan kelas lawas `h-screen` atau `w-screen` — Wajib gunakan unit dinamis Tailwind v4 (`h-dvh`, `w-dvw`).
- Lakukan perataan kanan (text-right) untuk semua nilai finansial dan metrik kuantitas.
