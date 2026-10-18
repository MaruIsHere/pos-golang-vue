package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"pos-backend/internal/models"

	"github.com/glebarez/sqlite"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
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

	// MySQL: DisableForeignKeyConstraintWhenMigrating mencegah GORM membuat FK
	// constraint baru saat AutoMigrate, sehingga kolom bebas dimodifikasi.
	gormCfg := &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: engine == "mysql",
	}

	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// SQLite: serialize semua akses DB lewat satu koneksi.
	if engine != "mysql" {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(1)
	}

	// Migrasi otomatis tabel artists -> merks dan kolom artist -> merk jika masih ada
	migrateArtistToMerk(db)

	// MySQL: hapus semua FK constraints yang sudah ada sebelum AutoMigrate.
	// Ini solusi definitif error 1025 (errno 150): MySQL tidak bisa ALTER TABLE
	// pada kolom yang terikat FK, bahkan dengan FOREIGN_KEY_CHECKS=0.
	// Karena DisableForeignKeyConstraintWhenMigrating=true, FK tidak akan dibuat ulang.
	if engine == "mysql" {
		dropAllMySQLForeignKeys(db)
	}

	// Auto Migration SaaS (UUID Multi-Tenant)
	err = db.AutoMigrate(
		&models.Merchant{},
		&models.Outlet{},
		&models.Expense{},
		&models.User{},
		&models.Category{},
		&models.Merk{},
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

// dropAllMySQLForeignKeys menghapus semua FK constraints di database aktif.
// Diperlukan sebelum AutoMigrate agar MySQL bisa ALTER TABLE bebas (error 1025 errno 150).
// FK tidak akan dibuat ulang karena DisableForeignKeyConstraintWhenMigrating=true.
func dropAllMySQLForeignKeys(db *gorm.DB) {
	var constraints []struct {
		TableName      string `gorm:"column:TABLE_NAME"`
		ConstraintName string `gorm:"column:CONSTRAINT_NAME"`
	}
	if err := db.Raw(`
		SELECT TABLE_NAME, CONSTRAINT_NAME
		FROM information_schema.TABLE_CONSTRAINTS
		WHERE CONSTRAINT_TYPE = 'FOREIGN KEY'
		AND TABLE_SCHEMA = DATABASE()
	`).Scan(&constraints).Error; err != nil {
		log.Printf("[warn] gagal query FK constraints: %v", err)
		return
	}
	for _, c := range constraints {
		sql := fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`", c.TableName, c.ConstraintName)
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("[warn] gagal drop FK %s.%s: %v", c.TableName, c.ConstraintName, err)
		} else {
			log.Printf("[info] FK constraint `%s`.`%s` berhasil dihapus", c.TableName, c.ConstraintName)
		}
	}
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

	// Seed Merks (Brands)
	var merkCount int64
	db.Model(&models.Merk{}).Count(&merkCount)
	if merkCount == 0 {
		defaultMerks := []string{"House Blend", "Indofood", "Nestle", "Sosro", "Umum"}
		for _, m := range defaultMerks {
			db.Create(&models.Merk{
				MerchantID: &defaultMerchant.ID,
				Name:       m,
			})
		}
	}

	// Seed Product Types
	var typeCount int64
	db.Model(&models.ProductType{}).Count(&typeCount)
	if typeCount == 0 {
		defaultTypes := []string{"FNB", "Retail", "Jasa"}
		for _, t := range defaultTypes {
			db.Create(&models.ProductType{
				MerchantID: &defaultMerchant.ID,
				Name:       t,
			})
		}
	}

	// Seed Categories
	var catCount int64
	db.Model(&models.Category{}).Count(&catCount)
	catIDs := map[string]uuid.UUID{}

	if catCount == 0 {
		cats := []models.Category{
			{MerchantID: &defaultMerchant.ID, Name: "Makanan Utama", Icon: "utensils"},
			{MerchantID: &defaultMerchant.ID, Name: "Snack & Camilan", Icon: "cookie"},
			{MerchantID: &defaultMerchant.ID, Name: "Kopi", Icon: "coffee"},
			{MerchantID: &defaultMerchant.ID, Name: "Non-Kopi", Icon: "cup-soda"},
			{MerchantID: &defaultMerchant.ID, Name: "Dessert", Icon: "cake-slice"},
			{MerchantID: &defaultMerchant.ID, Name: "Paket Hemat", Icon: "box"},
		}
		for i := range cats {
			if err := db.Create(&cats[i]).Error; err != nil {
				log.Printf("[seeder] Gagal membuat kategori %q: %v", cats[i].Name, err)
				continue
			}
			catIDs[cats[i].Name] = cats[i].ID
		}
	} else {
		var existingCats []models.Category
		db.Where("merchant_id = ?", defaultMerchant.ID).Find(&existingCats)
		for _, c := range existingCats {
			catIDs[c.Name] = c.ID
		}
	}

	// Seed Products (5 produk contoh standar, 1 per kategori)
	var prodCount int64
	db.Model(&models.Product{}).Count(&prodCount)
	if prodCount == 0 {
		var catMakanan, catSnack, catKopi, catNonKopi, catDessert models.Category
		db.Where("merchant_id = ? AND name = ?", defaultMerchant.ID, "Makanan Utama").First(&catMakanan)
		db.Where("merchant_id = ? AND name = ?", defaultMerchant.ID, "Snack & Camilan").First(&catSnack)
		db.Where("merchant_id = ? AND name = ?", defaultMerchant.ID, "Kopi").First(&catKopi)
		db.Where("merchant_id = ? AND name = ?", defaultMerchant.ID, "Non-Kopi").First(&catNonKopi)
		db.Where("merchant_id = ? AND name = ?", defaultMerchant.ID, "Dessert").First(&catDessert)

		products := []models.Product{
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catMakanan.ID, Name: "Nasi Goreng Special", Merk: "Umum", ProductType: "FNB", Price: 28000, CostPrice: 15000, Stock: 50, Unit: "pcs", IsActive: true},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catSnack.ID, Name: "French Fries", Merk: "Umum", ProductType: "FNB", Price: 18000, CostPrice: 7000, Stock: 50, Unit: "pcs", IsActive: true},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catKopi.ID, Name: "Kopi Susu Gula Aren", Merk: "House Blend", ProductType: "FNB", Price: 22000, CostPrice: 8000, Stock: 200, Unit: "pcs", IsActive: true},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catNonKopi.ID, Name: "Matcha Latte (Ice)", Merk: "Umum", ProductType: "FNB", Price: 30000, CostPrice: 13000, Stock: 100, Unit: "pcs", IsActive: true},
			{MerchantID: defaultMerchant.ID, OutletID: &defaultOutlet.ID, CategoryID: &catDessert.ID, Name: "Cheesecake Slice", Merk: "Umum", ProductType: "FNB", Price: 25000, CostPrice: 12000, Stock: 20, Unit: "pcs", IsActive: true},
		}
		for _, prod := range products {
			prod.IsMaster = false
			prod.OutletID = &defaultOutlet.ID
			if err := db.Create(&prod).Error; err != nil {
				log.Printf("Gagal membuat produk toko %q: %v", prod.Name, err)
			}
		}
		log.Println("[seeder] Seeding 5 produk contoh default berhasil.")
	}
}



func migrateArtistToMerk(db *gorm.DB) {
	// 1. Rename table 'artists' to 'merks'
	if db.Migrator().HasTable("artists") {
		if !db.Migrator().HasTable("merks") {
			log.Println("Migrasi: Mengubah nama tabel 'artists' menjadi 'merks'...")
			if err := db.Migrator().RenameTable("artists", "merks"); err != nil {
				log.Printf("Gagal rename table artists via migrator: %v, mencoba SQL mentah...", err)
				_ = db.Exec("ALTER TABLE artists RENAME TO merks").Error
			} else {
				log.Println("✔️ Tabel 'artists' berhasil diubah menjadi 'merks'.")
			}
		} else {
			// Tabel merks sudah ada, salin data jika merks masih kosong lalu drop artists
			var merksCount int64
			db.Table("merks").Count(&merksCount)
			if merksCount == 0 {
				_ = db.Exec("INSERT INTO merks SELECT * FROM artists").Error
			}
			_ = db.Migrator().DropTable("artists")
		}
	}

	// 2. Rename column 'artist' to 'merk' in 'products' table
	if db.Migrator().HasColumn("products", "artist") && !db.Migrator().HasColumn("products", "merk") {
		log.Println("Migrasi: Mengubah kolom 'artist' menjadi 'merk' pada tabel 'products'...")
		if err := db.Migrator().RenameColumn("products", "artist", "merk"); err != nil {
			_ = db.Exec("ALTER TABLE products CHANGE artist merk VARCHAR(100)").Error
			_ = db.Exec("ALTER TABLE products RENAME COLUMN artist TO merk").Error
			if !db.Migrator().HasColumn("products", "merk") {
				_ = db.Exec("ALTER TABLE products ADD COLUMN merk VARCHAR(100)").Error
				_ = db.Exec("UPDATE products SET merk = artist WHERE (merk IS NULL OR merk = '') AND artist IS NOT NULL").Error
			}
		} else {
			log.Println("✔️ Kolom 'artist' pada tabel 'products' berhasil diubah menjadi 'merk'.")
		}
	}

	// 3. Rename column 'artist' to 'merk' in 'order_items' table
	if db.Migrator().HasColumn("order_items", "artist") && !db.Migrator().HasColumn("order_items", "merk") {
		log.Println("Migrasi: Mengubah kolom 'artist' menjadi 'merk' pada tabel 'order_items'...")
		if err := db.Migrator().RenameColumn("order_items", "artist", "merk"); err != nil {
			_ = db.Exec("ALTER TABLE order_items CHANGE artist merk VARCHAR(100)").Error
			_ = db.Exec("ALTER TABLE order_items RENAME COLUMN artist TO merk").Error
			if !db.Migrator().HasColumn("order_items", "merk") {
				_ = db.Exec("ALTER TABLE order_items ADD COLUMN merk VARCHAR(100)").Error
				_ = db.Exec("UPDATE order_items SET merk = artist WHERE (merk IS NULL OR merk = '') AND artist IS NOT NULL").Error
			}
		} else {
			log.Println("✔️ Kolom 'artist' pada tabel 'order_items' berhasil diubah menjadi 'merk'.")
		}
	}
}
