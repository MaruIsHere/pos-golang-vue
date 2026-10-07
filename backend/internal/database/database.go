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
		log.Printf("Connecting to SQLite Database file: %s", sqlitePath)
		dialector = sqlite.Open(sqlitePath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
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
		catMinuman := models.Category{MerchantID: &defaultMerchant.ID, Name: "Minuman", Icon: "coffee"}
		db.Create(&catMakanan)
		db.Create(&catMinuman)

		// Seed Products
		products := []models.Product{
			{
				MerchantID:  defaultMerchant.ID,
				OutletID:    &defaultOutlet.ID,
				CategoryID:  &catMakanan.ID,
				Name:        "Nasi Goreng Special",
				ProductType: "FNB",
				Price:       28000,
				CostPrice:   15000,
				Stock:       50,
			},
		}
		for _, prod := range products {
			db.Create(&prod)
		}
	}
}
