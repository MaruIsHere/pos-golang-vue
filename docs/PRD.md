# Product Requirements Document (PRD)
## POS Kasir Pro

### 1. Overview
POS Kasir Pro is a modern, web-based Point of Sale (POS) system built with Go (backend) and Vue.js/TailwindCSS (frontend). It supports multi-platform deployment (Windows, Linux, macOS) through standalone binaries and features a Progressive Web App (PWA) frontend.

### 2. Core Features
- **Product Management:** Add, edit, delete products with barcode support.
- **Inventory & Stock Tracking:** Real-time stock movements with history.
- **Order Processing:** Quick checkout with cart, discount, and tax calculation.
- **AI OCR Integration:** Automatically scan and extract payment amounts from QRIS/Bank Transfer receipts using an embedded Rust-based AI engine (`ocrs`).
- **Role-based Access Control:** Strict permission system (Owner, Administrator, Admin, Kepala Kasir, Kasir).
- **Multi-tenant/Multi-outlet Support:** Designed with merchant_id and outlet_id for scalability.
- **Dynamic Database Selection:** Seamless switching between SQLite (local) and MySQL (production).

### 3. User Experience
- Responsive UI optimized for desktop and mobile touchscreens.
- Skeleton loading states for smooth data transitions.
- PWA support for offline caching and standalone installation.

### 4. Technical Constraints
- Must run efficiently on low-resource hardware (e.g., AMD A4 processors).
- Backend compiled as a single binary with zero external dependencies (CGO disabled).
- OCR module runs locally (offline) via compiled `ocrs` binary.
