package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"pos-backend/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ensureMySQLDatabaseExists(dsn string) {
	parts := strings.SplitN(dsn, "/", 2)
	if len(parts) == 2 {
		dbAndParams := parts[1]
		subParts := strings.SplitN(dbAndParams, "?", 2)
		dbName := subParts[0]
		if dbName != "" {
			params := ""
			if len(subParts) > 1 {
				params = "?" + subParts[1]
			}
			baseDSN := parts[0] + "/" + params
			db, err := sql.Open("mysql", baseDSN)
			if err == nil {
				defer db.Close()
				_, _ = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName))
			}
		}
	}
}

func InitDB(engine string, mysqlDsn string, sqlitePath string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	if engine == "mysql" {
		log.Printf("Connecting to MySQL Database with DSN: %s", mysqlDsn)
		ensureMySQLDatabaseExists(mysqlDsn)
		dialector = mysql.Open(mysqlDsn)
	} else {
		log.Printf("Connecting to SQLite Database file: %s", sqlitePath)
		dialector = sqlite.Open(sqlitePath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// Auto Migration
	err = db.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.StoreSetting{},
		&models.Voucher{},
		&models.Customer{},
		&models.StockMovement{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate models: %w", err)
	}

	DB = db
	SeedInitialData(db, engine, mysqlDsn)
	return db, nil
}

func SeedInitialData(db *gorm.DB, engine string, mysqlDsn string) {
	// Seed Customers if empty
	var custCount int64
	db.Model(&models.Customer{}).Count(&custCount)
	if custCount == 0 {
		initialCusts := []models.Customer{
			{Name: "Budi Santoso", Phone: "081234567890", Email: "budi@gmail.com", Address: "Jakarta", Points: 120},
			{Name: "Siti Rahma", Phone: "081987654321", Email: "siti@gmail.com", Address: "Bandung", Points: 85},
			{Name: "Dewi Lestari", Phone: "085612345678", Email: "dewi@gmail.com", Address: "Surabaya", Points: 210},
		}
		for _, c := range initialCusts {
			db.Create(&c)
		}
	}

	// Seed Vouchers if empty
	var voucherCount int64
	db.Model(&models.Voucher{}).Count(&voucherCount)
	if voucherCount == 0 {
		initialVouchers := []models.Voucher{
			{Code: "DISKON10", Type: "percent", Value: 10, Description: "Diskon 10% Semua Produk", IsActive: true},
			{Code: "DISKON20", Type: "percent", Value: 20, Description: "Diskon 20% Promo Spesial", IsActive: true},
			{Code: "HEMAT10K", Type: "flat", Value: 10000, Description: "Potongan Rp 10.000", IsActive: true},
			{Code: "HEMAT50K", Type: "flat", Value: 50000, Description: "Potongan Rp 50.000", IsActive: true},
			{Code: "POSHEMAT", Type: "percent", Value: 15, Description: "Diskon Kasir 15%", IsActive: true},
		}
		for _, v := range initialVouchers {
			db.Create(&v)
		}
	}

	// Seed Store Setting if not exists
	var count int64
	db.Model(&models.StoreSetting{}).Count(&count)
	if count == 0 {
		setting := models.StoreSetting{
			StoreName:     "KASIR COFFEE & BISTRO",
			Address:       "Jl. Boulevard Utama No. 88, Jakarta Selatan",
			Phone:         "0812-9876-5432",
			ReceiptFooter: "Terima kasih atas kunjungan Anda!\nSimpan struk ini sebagai bukti pembayaran resmi.",
			TaxPercentage: 10,
			DbEngine:      engine,
			MysqlDsn:      mysqlDsn,
		}
		db.Create(&setting)
	}

	// Seed Categories if empty
	db.Model(&models.Category{}).Count(&count)
	if count == 0 {
		categories := []models.Category{
			{Name: "Makanan Utama", Icon: "utensils"},
			{Name: "Minuman", Icon: "coffee"},
			{Name: "Cemilan & Snack", Icon: "cookie"},
			{Name: "Dessert", Icon: "ice-cream"},
		}
		for _, cat := range categories {
			db.Create(&cat)
		}

		// Seed Products
		var catMakanan, catMinuman, catSnack, catDessert models.Category
		db.Where("name = ?", "Makanan Utama").First(&catMakanan)
		db.Where("name = ?", "Minuman").First(&catMinuman)
		db.Where("name = ?", "Cemilan & Snack").First(&catSnack)
		db.Where("name = ?", "Dessert").First(&catDessert)

		products := []models.Product{
			{
				CategoryID: catMakanan.ID,
				Name:       "Nasi Goreng Special Egg",
				Price:      28000,
				CostPrice:  15000,
				Stock:      50,
				Barcode:    "8991001",
				ImageURL:   "https://images.unsplash.com/photo-1603133872878-684f208fb84b?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catMakanan.ID,
				Name:       "Ayam Geprek Sambal Korek",
				Price:      25000,
				CostPrice:  13000,
				Stock:      40,
				Barcode:    "8991002",
				ImageURL:   "https://images.unsplash.com/photo-1626082927389-6cd097cdc6ec?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catMakanan.ID,
				Name:       "Mie Goreng Seafood Premium",
				Price:      32000,
				CostPrice:  18000,
				Stock:      30,
				Barcode:    "8991003",
				ImageURL:   "https://images.unsplash.com/photo-1612927601601-6638404737ce?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catMinuman.ID,
				Name:       "Kopi Susu Gula Aren",
				Price:      18000,
				CostPrice:  8000,
				Stock:      100,
				Barcode:    "8992001",
				ImageURL:   "https://images.unsplash.com/photo-1541167760496-1628856ab772?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catMinuman.ID,
				Name:       "Es Teh Manis Jumbo",
				Price:      8000,
				CostPrice:  2500,
				Stock:      200,
				Barcode:    "8992002",
				ImageURL:   "https://images.unsplash.com/photo-1556679343-c7306c1976bc?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catMinuman.ID,
				Name:       "Matcha Latte Ice",
				Price:      22000,
				CostPrice:  11000,
				Stock:      60,
				Barcode:    "8992003",
				ImageURL:   "https://images.unsplash.com/photo-1536256263959-770b48d82b0a?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catSnack.ID,
				Name:       "French Fries Crispy",
				Price:      15000,
				CostPrice:  7000,
				Stock:      80,
				Barcode:    "8993001",
				ImageURL:   "https://images.unsplash.com/photo-1573080496219-bb080dd4f877?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catSnack.ID,
				Name:       "Roti Bakar Coklat Keju",
				Price:      18000,
				CostPrice:  8500,
				Stock:      45,
				Barcode:    "8993002",
				ImageURL:   "https://images.unsplash.com/photo-1584776296944-ab6fb57b0bdd?w=400&auto=format&fit=crop&q=80",
			},
			{
				CategoryID: catDessert.ID,
				Name:       "Choco Lava Cake",
				Price:      25000,
				CostPrice:  12000,
				Stock:      25,
				Barcode:    "8994001",
				ImageURL:   "https://images.unsplash.com/photo-1606313564200-e75d5e30476c?w=400&auto=format&fit=crop&q=80",
			},
		}

		for _, prod := range products {
			db.Create(&prod)
		}
	}
}
