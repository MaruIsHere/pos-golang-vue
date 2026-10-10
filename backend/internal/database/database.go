package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"pos-backend/internal/models"

	"github.com/glebarez/sqlite"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
		sqliteDSN := sqlitePath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=synchronous(NORMAL)"
		log.Printf("Connecting to SQLite Database file: %s", sqliteDSN)
		dialector = sqlite.Open(sqliteDSN)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// SQLite: serialize all DB access through a single connection.
	// WAL mode allows concurrent reads, but writes must be serialized.
	// This lets Go's connection pool queue writes instead of SQLite rejecting them.
	if engine != "mysql" {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(1)
	}

	// Auto Migration SaaS (UUID Multi-Tenant)
	err = db.AutoMigrate(
		&models.Merchant{},
		&models.Outlet{},
		&models.Expense{},
		&models.User{},
		&models.Category{},
		&models.Artist{},
		&models.ProductType{},
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
	// Seed Default Merchant
	var merchantCount int64
	db.Model(&models.Merchant{}).Count(&merchantCount)
	var defaultMerchant models.Merchant

	if merchantCount == 0 {
		defaultMerchant = models.Merchant{
			Name:     "Owner Default",
			Email:    "admin@pos.com",
			Phone:    "0811111111",
			IsActive: true,
		}
		db.Create(&defaultMerchant)
	} else {
		db.First(&defaultMerchant)
	}

	// Seed Default Outlet (Gudang & Lapak)
	var outletCount int64
	db.Model(&models.Outlet{}).Count(&outletCount)
	var defaultOutlet models.Outlet

	if outletCount == 0 {
		defaultOutlet = models.Outlet{
			MerchantID:  defaultMerchant.ID,
			Name:        "Toko Utama",
			Code:        "STORE-001",
			Address:     "Jl. Merdeka Utama No. 1",
			Phone:       "0812-3456-7890",
			IsWarehouse: true,
			IsActive:    true,
		}
		db.Create(&defaultOutlet)
	} else {
		db.First(&defaultOutlet)
	}

	// Seed Default Accounts
	defaultUsers := []struct {
		username string
		password string
		role     string
	}{
		{"owner", "owner123", "OWNER"},
		{"kepalakasir", "kepala123", "KEPALA_KASIR"},
		{"kasir", "kasir123", "KASIR"},
		{"admin", "admin123", "ADMIN"},
	}

	for _, u := range defaultUsers {
		var count int64
		db.Model(&models.User{}).Where("username = ?", u.username).Count(&count)
		if count == 0 {
			hashed, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
			db.Create(&models.User{
				MerchantID: &defaultMerchant.ID,
				OutletID:   &defaultOutlet.ID,
				Username:   u.username,
				Password:   string(hashed),
				Role:       models.UserRole(u.role),
			})
		}
	}

	// Seed Customers
	var custCount int64
	db.Model(&models.Customer{}).Count(&custCount)
	if custCount == 0 {
		initialCusts := []models.Customer{
			{MerchantID: defaultMerchant.ID, Name: "Budi Santoso", Phone: "081234567890", Email: "budi@gmail.com", Address: "Jakarta", Points: 120},
		}
		for _, c := range initialCusts {
			db.Create(&c)
		}
	}

	// Seed Store Setting
	var settingCount int64
	db.Model(&models.StoreSetting{}).Count(&settingCount)
	if settingCount == 0 {
		db.Create(&models.StoreSetting{
			MerchantID:               defaultMerchant.ID,
			OutletID:                 &defaultOutlet.ID,
			StoreName:                "KASIR COFFEE & BISTRO",
			Address:                  "Jl. Boulevard Utama No. 88",
			Phone:                    "0812-9876-5432",
			ReceiptFooter:            "Terima kasih atas kunjungan Anda!",
			TaxPercentage:            10,
			MemberDiscountPercentage: 5,
			DbEngine:                 engine,
			MysqlDsn:                 mysqlDsn,
		})
	}

	// Seed Categories
	var catCount int64
	db.Model(&models.Category{}).Count(&catCount)
	if catCount == 0 {
		catMakanan := models.Category{MerchantID: &defaultMerchant.ID, Name: "Makanan Utama", Icon: "utensils"}
		catSnack := models.Category{MerchantID: &defaultMerchant.ID, Name: "Snack & Camilan", Icon: "cookie"}
		catKopi := models.Category{MerchantID: &defaultMerchant.ID, Name: "Kopi", Icon: "coffee"}
		catNonKopi := models.Category{MerchantID: &defaultMerchant.ID, Name: "Non-Kopi", Icon: "cup-soda"}
		catDessert := models.Category{MerchantID: &defaultMerchant.ID, Name: "Dessert", Icon: "cake-slice"}
		catPaket := models.Category{MerchantID: &defaultMerchant.ID, Name: "Paket Hemat", Icon: "box"}
		db.Create(&catMakanan)
		db.Create(&catSnack)
		db.Create(&catKopi)
		db.Create(&catNonKopi)
		db.Create(&catDessert)
		db.Create(&catPaket)

		// Seed Products
		products := []models.Product{
			// === MAKANAN UTAMA ===
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Nasi Goreng Special", ProductType: "FNB", Price: 28000, CostPrice: 15000, Stock: 50},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Nasi Goreng Seafood", ProductType: "FNB", Price: 32000, CostPrice: 18000, Stock: 30},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Mie Goreng Jawa", ProductType: "FNB", Price: 25000, CostPrice: 12000, Stock: 40},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Indomie Goreng Telur", ProductType: "FNB", Price: 18000, CostPrice: 8000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Indomie Kuah Telur", ProductType: "FNB", Price: 18000, CostPrice: 8000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Chicken Katsu Rice", ProductType: "FNB", Price: 35000, CostPrice: 20000, Stock: 25},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Beef Teriyaki Rice", ProductType: "FNB", Price: 38000, CostPrice: 22000, Stock: 20},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Sandwich Tuna Melt", ProductType: "FNB", Price: 30000, CostPrice: 16000, Stock: 15},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Croissant Ham & Cheese", ProductType: "FNB", Price: 28000, CostPrice: 14000, Stock: 20},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Spaghetti Bolognese", ProductType: "FNB", Price: 33000, CostPrice: 17000, Stock: 20},

			// === SNACK & CAMILAN ===
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "French Fries", ProductType: "FNB", Price: 18000, CostPrice: 7000, Stock: 50},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "Chicken Wings (5pcs)", ProductType: "FNB", Price: 25000, CostPrice: 13000, Stock: 30},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "Pisang Goreng Keju", ProductType: "FNB", Price: 15000, CostPrice: 6000, Stock: 40},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "Roti Bakar Coklat", ProductType: "FNB", Price: 15000, CostPrice: 5000, Stock: 35},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "Dimsum Ayam (4pcs)", ProductType: "FNB", Price: 20000, CostPrice: 10000, Stock: 25},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "Onion Rings", ProductType: "FNB", Price: 16000, CostPrice: 6000, Stock: 30},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "Kentang Wedges", ProductType: "FNB", Price: 20000, CostPrice: 8000, Stock: 30},

			// === KOPI ===
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Espresso", ProductType: "FNB", Price: 15000, CostPrice: 5000, Stock: 200},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Americano (Hot)", ProductType: "FNB", Price: 18000, CostPrice: 6000, Stock: 200},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Americano (Ice)", ProductType: "FNB", Price: 20000, CostPrice: 7000, Stock: 200},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Cafe Latte (Hot)", ProductType: "FNB", Price: 25000, CostPrice: 9000, Stock: 200},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Cafe Latte (Ice)", ProductType: "FNB", Price: 27000, CostPrice: 10000, Stock: 200},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Cappuccino", ProductType: "FNB", Price: 25000, CostPrice: 9000, Stock: 200},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Caramel Macchiato", ProductType: "FNB", Price: 30000, CostPrice: 12000, Stock: 150},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Mocha Latte", ProductType: "FNB", Price: 28000, CostPrice: 11000, Stock: 150},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Vanilla Latte", ProductType: "FNB", Price: 28000, CostPrice: 11000, Stock: 150},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Hazelnut Latte", ProductType: "FNB", Price: 28000, CostPrice: 11000, Stock: 150},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Kopi Susu Gula Aren", ProductType: "FNB", Price: 22000, CostPrice: 8000, Stock: 200},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "V60 Single Origin", ProductType: "FNB", Price: 35000, CostPrice: 15000, Stock: 50},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Affogato", ProductType: "FNB", Price: 30000, CostPrice: 12000, Stock: 80},

			// === NON-KOPI ===
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Matcha Latte (Hot)", ProductType: "FNB", Price: 28000, CostPrice: 12000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Matcha Latte (Ice)", ProductType: "FNB", Price: 30000, CostPrice: 13000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Coklat Panas", ProductType: "FNB", Price: 22000, CostPrice: 9000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Ice Chocolate", ProductType: "FNB", Price: 25000, CostPrice: 10000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Thai Tea", ProductType: "FNB", Price: 20000, CostPrice: 7000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Taro Milk", ProductType: "FNB", Price: 22000, CostPrice: 8000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Lemon Tea", ProductType: "FNB", Price: 15000, CostPrice: 5000, Stock: 100},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Fresh Orange Juice", ProductType: "FNB", Price: 20000, CostPrice: 8000, Stock: 80},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Strawberry Smoothie", ProductType: "FNB", Price: 25000, CostPrice: 10000, Stock: 60},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Air Mineral", ProductType: "FNB", Price: 5000, CostPrice: 2000, Stock: 300},

			// === DESSERT ===
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catDessert.ID, Name: "Brownies Coklat", ProductType: "FNB", Price: 18000, CostPrice: 8000, Stock: 30},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catDessert.ID, Name: "Cheesecake Slice", ProductType: "FNB", Price: 25000, CostPrice: 12000, Stock: 20},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catDessert.ID, Name: "Tiramisu", ProductType: "FNB", Price: 28000, CostPrice: 14000, Stock: 15},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catDessert.ID, Name: "Banana Split", ProductType: "FNB", Price: 25000, CostPrice: 10000, Stock: 20},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catDessert.ID, Name: "Waffle Ice Cream", ProductType: "FNB", Price: 30000, CostPrice: 13000, Stock: 20},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catDessert.ID, Name: "Pancake Maple Syrup", ProductType: "FNB", Price: 22000, CostPrice: 9000, Stock: 25},

			// === PAKET HEMAT ===
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catPaket.ID, Name: "Paket Nasi + Es Teh", ProductType: "FNB", Price: 35000, CostPrice: 18000, Stock: 50},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catPaket.ID, Name: "Paket Indomie + Kopi Susu", ProductType: "FNB", Price: 32000, CostPrice: 14000, Stock: 50},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catPaket.ID, Name: "Paket Snack Platter", ProductType: "FNB", Price: 45000, CostPrice: 22000, Stock: 30},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catPaket.ID, Name: "Paket Berdua (2 Kopi + 1 Snack)", ProductType: "FNB", Price: 55000, CostPrice: 25000, Stock: 30},
		}
		for _, prod := range products {
			db.Create(&prod)
		}
	}
}
