package handlers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// === VOUCHER HANDLERS ===

func GetVouchers(c *gin.Context) {
	var vouchers []models.Voucher
	if err := database.DB.Order("created_at desc").Find(&vouchers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vouchers)
}

func CreateVoucher(c *gin.Context) {
	var voucher models.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if voucher.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode voucher tidak boleh kosong"})
		return
	}

	voucher.IsActive = true
	if err := database.DB.Create(&voucher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat voucher (kode mungkin sudah ada)"})
		return
	}
	c.JSON(http.StatusCreated, voucher)
}

func DeleteVoucher(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Voucher{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voucher berhasil dihapus"})
}
