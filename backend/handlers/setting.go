package handlers

import (
	"fmt"
	"net/http"
	"pos-backend/config"
	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// === STORE & DATABASE SETTINGS HANDLERS ===

func GetSettings(c *gin.Context) {
	var setting models.StoreSetting
	if err := database.DB.First(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, setting)
}

func UpdateSettings(c *gin.Context) {
	var setting models.StoreSetting
	if err := database.DB.First(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input models.StoreSetting
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting.StoreName = input.StoreName
	setting.Address = input.Address
	setting.Phone = input.Phone
	setting.ReceiptFooter = input.ReceiptFooter
	setting.TaxPercentage = input.TaxPercentage
	setting.QrisImageUrl = input.QrisImageUrl

	database.DB.Save(&setting)
	c.JSON(http.StatusOK, setting)
}

type SwitchDbInput struct {
	Engine   string `json:"engine"`    // "sqlite" or "mysql"
	MysqlDsn string `json:"mysql_dsn"` // optional DSN
}

func SwitchDatabase(c *gin.Context) {
	var input SwitchDbInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Engine != "sqlite" && input.Engine != "mysql" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Engine database harus 'sqlite' atau 'mysql'"})
		return
	}

	cfg := config.AppConfig
	cfg.DbEngine = input.Engine
	if input.MysqlDsn != "" {
		cfg.MysqlDsn = input.MysqlDsn
	}

	newDB, err := database.InitDB(cfg.DbEngine, cfg.MysqlDsn, cfg.SqlitePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Gagal berpindah ke %s: %v", input.Engine, err)})
		return
	}

	database.DB = newDB
	config.SaveConfig(cfg)

	// Update setting record
	var setting models.StoreSetting
	if err := database.DB.First(&setting).Error; err == nil {
		setting.DbEngine = cfg.DbEngine
		setting.MysqlDsn = cfg.MysqlDsn
		database.DB.Save(&setting)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("Berhasil berpindah ke basis data %s", input.Engine),
		"db_engine": input.Engine,
	})
}
