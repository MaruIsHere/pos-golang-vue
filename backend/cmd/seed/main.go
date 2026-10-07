package main

import (
	"log"

	"pos-backend/internal/config"
	"pos-backend/internal/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("Memulai proses Database Seeding...")

	cfg := config.LoadConfig()

	// InitDB otomatis memanggil koneksi dan menjalankan SeedInitialData
	// karena SeedInitialData di database.go menggunakan Count() == 0,
	// skrip ini sangat aman dijalankan berkali-kali tanpa menduplikasi data.
	db, err := database.InitDB(cfg.DbEngine, cfg.MysqlDsn, cfg.SqlitePath)
	if err != nil {
		log.Fatalf("Gagal inisialisasi database: %v", err)
	}

	log.Println("Memeriksa dan mengeksekusi Seeder...")
	
	// Kita bisa panggil secara eksplisit lagi jika ingin log tambahan
	database.SeedInitialData(db, cfg.DbEngine, cfg.MysqlDsn)

	log.Println("🌱 Database Seeding BERHASIL! Akun default dan data awal telah siap.")
}
