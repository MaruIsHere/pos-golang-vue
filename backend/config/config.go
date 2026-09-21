package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port     string `json:"port"`
	DbEngine string `json:"db_engine"` // "sqlite" or "mysql"
	SqlitePath string `json:"sqlite_path"`
	MysqlDsn string `json:"mysql_dsn"`
}

var AppConfig Config

func LoadConfig() Config {
	configFile := "config.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		AppConfig = Config{
			Port:       "8080",
			DbEngine:   "sqlite",
			SqlitePath: "pos.db",
			MysqlDsn:   "root:@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local",
		}
		SaveConfig(AppConfig)
		return AppConfig
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		AppConfig = Config{
			Port:       "8080",
			DbEngine:   "sqlite",
			SqlitePath: "pos.db",
			MysqlDsn:   "root:@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local",
		}
		return AppConfig
	}

	json.Unmarshal(data, &AppConfig)
	return AppConfig
}

func SaveConfig(cfg Config) error {
	AppConfig = cfg
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}
