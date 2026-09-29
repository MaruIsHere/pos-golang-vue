package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Port       string `json:"port"`
	DbEngine   string `json:"db_engine"` // "sqlite" or "mysql"
	SqlitePath string `json:"sqlite_path"`
	MysqlDsn   string `json:"mysql_dsn"`
}

var AppConfig Config
var activeConfigDir string

func findConfigDir() string {
	for _, dir := range []string{".", "backend"} {
		if _, err := os.Stat(filepath.Join(dir, "config.local.json")); err == nil {
			return dir
		}
	}
	for _, dir := range []string{".", "backend"} {
		if _, err := os.Stat(filepath.Join(dir, "config.json")); err == nil {
			return dir
		}
	}
	return "."
}

func defaultConfig() Config {
	return Config{
		Port:       "8080",
		DbEngine:   "sqlite",
		SqlitePath: "pos.db",
		MysqlDsn:   "root:@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local",
	}
}

func resolveSQLitePath(cfg Config) Config {
	if cfg.SqlitePath == "" || cfg.SqlitePath == ":memory:" || filepath.IsAbs(cfg.SqlitePath) {
		return cfg
	}
	if len(cfg.SqlitePath) >= 5 && cfg.SqlitePath[:5] == "file:" {
		return cfg
	}
	path, err := filepath.Abs(filepath.Join(activeConfigDir, cfg.SqlitePath))
	if err == nil {
		cfg.SqlitePath = path
	}
	return cfg
}

func LoadConfig() Config {
	activeConfigDir = findConfigDir()
	localConfigFile := filepath.Join(activeConfigDir, "config.local.json")
	configFile := localConfigFile
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configFile = filepath.Join(activeConfigDir, "config.json")
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		AppConfig = defaultConfig()
		AppConfig = resolveSQLitePath(AppConfig)
		if saveErr := SaveConfig(AppConfig); saveErr != nil {
			log.Printf("Could not save local configuration: %v", saveErr)
		}
		return AppConfig
	}

	AppConfig = defaultConfig()
	if err := json.Unmarshal(data, &AppConfig); err != nil {
		log.Printf("Could not parse configuration %s: %v", configFile, err)
		AppConfig = defaultConfig()
	}
	AppConfig = resolveSQLitePath(AppConfig)
	if configFile != localConfigFile {
		if saveErr := SaveConfig(AppConfig); saveErr != nil {
			log.Printf("Could not save local configuration: %v", saveErr)
		}
	}
	return AppConfig
}

func SaveConfig(cfg Config) error {
	if activeConfigDir == "" {
		activeConfigDir = findConfigDir()
	}
	cfg = resolveSQLitePath(cfg)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(activeConfigDir, "config.local.json"), data, 0644); err != nil {
		return err
	}
	AppConfig = cfg
	return nil
}
