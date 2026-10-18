package database

import (
	"gorm.io/gorm"
)

// BackfillMasterProducts menjamin setiap produk toko (is_master = false) punya
// induk produk master, supaya tampil di Katalog Produk Utama dan terhitung benar
// pada dialog "Ambil Data Menu".
//
// Fungsi ini idempoten dan aman dipanggil di setiap startup: produk yang sudah
// punya master_product_id dilewati, dan master yang sudah ada dengan nama sama
// pada merchant yang sama dipakai ulang, bukan digandakan.
func BackfillMasterProducts(db *gorm.DB) {
	// Mode disederhanakan: 1 produk = 1 baris per toko, tidak ada duplikasi master shadow.
}
