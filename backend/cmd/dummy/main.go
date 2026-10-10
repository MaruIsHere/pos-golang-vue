package main

import (
	"log"
	"fmt"
	"pos-backend/internal/config"
	"pos-backend/internal/database"
	"pos-backend/internal/models"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.InitDB(cfg.DbEngine, cfg.MysqlDsn, "pos.db")
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}

	var merchant models.Merchant
	if err := db.First(&merchant).Error; err != nil {
		log.Fatal(err)
	}

	for i := 1; i <= 10; i++ {
		p := models.Product{
			Name:       fmt.Sprintf("Dummy Produk Master %d", i),
			MerchantID: merchant.ID,
			Price:      10000 + float64(i*1000),
			Stock:      100,
			Unit:       "pcs",
			IsMaster:   true,
		}
		if err := db.Create(&p).Error; err != nil {
			log.Println("Error:", err)
		} else {
			log.Println("Created:", p.Name)
		}
	}
}
