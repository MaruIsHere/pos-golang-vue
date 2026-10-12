package database

import (
	"log"

	"pos-backend/internal/models"

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
	var pending []models.Product
	if err := db.Where("is_master = ? AND master_product_id IS NULL", false).Find(&pending).Error; err != nil {
		log.Printf("Backfill master produk: gagal membaca daftar produk toko: %v", err)
		return
	}
	if len(pending) == 0 {
		return
	}

	created, linked := 0, 0
	for _, storeProd := range pending {
		var master models.Product
		err := db.Where("is_master = ? AND merchant_id = ? AND LOWER(name) = LOWER(?)",
			true, storeProd.MerchantID, storeProd.Name).First(&master).Error
		if err != nil {
			master = models.Product{
				MerchantID:       storeProd.MerchantID,
				OutletID:         nil,
				OriginalOutletID: storeProd.OutletID,
				CategoryID:       storeProd.CategoryID,
				Name:             storeProd.Name,
				Merk:             storeProd.Merk,
				ProductType:      storeProd.ProductType,
				Price:            storeProd.Price,
				CostPrice:        storeProd.CostPrice,
				Stock:            storeProd.Stock,
				Unit:             storeProd.Unit,
				Barcode:          storeProd.Barcode,
				ImageURL:         storeProd.ImageURL,
				IsActive:         storeProd.IsActive,
				IsMaster:         true,
			}
			if err := db.Create(&master).Error; err != nil {
				log.Printf("Backfill master produk: gagal membuat master untuk %q: %v", storeProd.Name, err)
				continue
			}
			created++
		}

		masterID := master.ID
		if err := db.Model(&models.Product{}).Where("id = ?", storeProd.ID).
			Update("master_product_id", masterID).Error; err != nil {
			log.Printf("Backfill master produk: gagal menautkan produk toko %q: %v", storeProd.Name, err)
			continue
		}
		linked++
	}

	log.Printf("Backfill master produk selesai: %d master dibuat, %d produk toko ditautkan.", created, linked)
}
