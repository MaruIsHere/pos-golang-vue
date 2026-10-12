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
			if dir := filepath.Dir(cfg.SqlitePath); dir != "" && dir != "." {
				_ = os.MkdirAll(dir, 0755)
			}
			for _, target := range []string{cfg.SqlitePath, cfg.SqlitePath + "-wal", cfg.SqlitePath + "-shm"} {
				if _, err := os.Stat(target); err == nil {
					if err := os.Remove(target); err != nil {
						log.Printf("Peringatan: Gagal menghapus %s: %v", filepath.Base(target), err)
					} else {
						log.Printf("✔️ File SQLite '%s' berhasil dihapus.", filepath.Base(target))
					}
				}
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
