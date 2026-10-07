package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"pos-backend/internal/config"
	"pos-backend/internal/database"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 1. Load config
	cfg := config.LoadConfig()

	log.Println("Memulai proses Migrate:Fresh...")

	// 2. Drop & Recreate Database berdasarkan Engine
	if cfg.DbEngine == "mysql" {
		log.Printf("Engine terdeteksi: MySQL. Mencoba me-reset database...")
		
		// Parse DSN untuk mendapatkan nama database
		// Format DSN: user:pass@tcp(host:port)/dbname?params
		parts := strings.SplitN(cfg.MysqlDsn, "/", 2)
		if len(parts) != 2 {
			log.Fatalf("DSN MySQL tidak valid: %s", cfg.MysqlDsn)
		}
		
		baseDSN := parts[0] + "/"
		dbAndParams := parts[1]
		subParts := strings.SplitN(dbAndParams, "?", 2)
		dbName := subParts[0]

		// Konek ke MySQL tanpa spesifik database
		db, err := sql.Open("mysql", baseDSN)
		if err != nil {
			log.Fatalf("Gagal konek ke MySQL server: %v", err)
		}
		defer db.Close()

		// Drop database
		_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", dbName))
		if err != nil {
			log.Fatalf("Gagal drop database %s: %v", dbName, err)
		}
		log.Printf("✔️ Database '%s' berhasil di-drop.", dbName)

		// Create database
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
		if err != nil {
			log.Fatalf("Gagal membuat database %s: %v", dbName, err)
		}
		log.Printf("✔️ Database '%s' berhasil dibuat ulang.", dbName)

	} else if cfg.DbEngine == "sqlite" {
		log.Printf("Engine terdeteksi: SQLite. Mencoba menghapus file...")
		
		if cfg.SqlitePath != "" && cfg.SqlitePath != ":memory:" {
			// Jika file ada, hapus
			if _, err := os.Stat(cfg.SqlitePath); err == nil {
				err := os.Remove(cfg.SqlitePath)
				if err != nil {
					log.Fatalf("Gagal menghapus file SQLite %s: %v", cfg.SqlitePath, err)
				}
				log.Printf("✔️ File SQLite '%s' berhasil dihapus.", filepath.Base(cfg.SqlitePath))
			} else {
				log.Printf("File SQLite '%s' belum ada. Aman.", filepath.Base(cfg.SqlitePath))
			}
		}
	} else {
		log.Fatalf("Engine database '%s' tidak didukung untuk skrip ini.", cfg.DbEngine)
	}

	log.Println("Menjalankan GORM AutoMigrate & Seeding awal...")

	// 3. Panggil InitDB yang di dalamnya sudah ada gorm.AutoMigrate dan Seeding data awal
	_, err := database.InitDB(cfg.DbEngine, cfg.MysqlDsn, cfg.SqlitePath)
	if err != nil {
		log.Fatalf("Gagal inisialisasi / migrasi database: %v", err)
	}

	log.Println("🚀 Migrate Fresh BERHASIL! Database telah direset ke skema terbaru.")
}
