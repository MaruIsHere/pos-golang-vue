package handlers

import (
	"net/http"
	"pos-backend/internal/database"
	"pos-backend/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type catalogInput struct {
	Name string `json:"name" binding:"required"`
}

func normalizeCatalogInput(c *gin.Context) (catalogInput, bool) {
	var input catalogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama wajib diisi"})
		return input, false
	}
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" || len(input.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama wajib diisi (maksimal 100 karakter)"})
		return input, false
	}
	return input, true
}

func GetMerks(c *gin.Context) {
	var merks []models.Merk
	if err := database.DB.Order("name").Find(&merks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar merk"})
		return
	}
	c.JSON(http.StatusOK, merks)
}

func CreateMerk(c *gin.Context) {
	input, ok := normalizeCatalogInput(c)
	if !ok {
		return
	}
	var count int64
	if err := database.DB.Model(&models.Merk{}).Where("LOWER(name) = LOWER(?)", input.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa nama merk"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Merk sudah tersedia"})
		return
	}
	merk := models.Merk{Name: input.Name}
	if err := database.DB.Create(&merk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan merk"})
		return
	}
	c.JSON(http.StatusCreated, merk)
}

func UpdateMerk(c *gin.Context) {
	input, ok := normalizeCatalogInput(c)
	if !ok {
		return
	}
	var merk models.Merk
	if err := database.DB.Where("id = ?", c.Param("id")).First(&merk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merk tidak ditemukan"})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Merk{}).Where("LOWER(name) = LOWER(?) AND id <> ?", input.Name, merk.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa nama merk"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Nama merk sudah digunakan"})
		return
	}
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Where("merk = ?", merk.Name).Update("merk", input.Name).Error; err != nil {
			return err
		}
		return tx.Model(&merk).Update("name", input.Name).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui merk"})
		return
	}
	merk.Name = input.Name
	c.JSON(http.StatusOK, merk)
}

func DeleteMerk(c *gin.Context) {
	var merk models.Merk
	if err := database.DB.Where("id = ?", c.Param("id")).First(&merk).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merk tidak ditemukan"})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Product{}).Where("merk = ?", merk.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa penggunaan merk"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Merk masih digunakan produk"})
		return
	}
	if err := database.DB.Delete(&merk).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus merk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Merk berhasil dihapus"})
}

// Backward-compatible aliases for legacy routes
var GetArtists = GetMerks
var CreateArtist = CreateMerk
var UpdateArtist = UpdateMerk
var DeleteArtist = DeleteMerk

func GetProductTypes(c *gin.Context) {
	var productTypes []models.ProductType
	if err := database.DB.Order("name").Find(&productTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar tipe produk"})
		return
	}
	c.JSON(http.StatusOK, productTypes)
}

func CreateProductType(c *gin.Context) {
	input, ok := normalizeCatalogInput(c)
	if !ok {
		return
	}
	var count int64
	if err := database.DB.Model(&models.ProductType{}).Where("LOWER(name) = LOWER(?)", input.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa tipe produk"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Tipe produk sudah tersedia"})
		return
	}
	productType := models.ProductType{Name: input.Name}
	if err := database.DB.Create(&productType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tipe produk"})
		return
	}
	c.JSON(http.StatusCreated, productType)
}

func UpdateProductType(c *gin.Context) {
	input, ok := normalizeCatalogInput(c)
	if !ok {
		return
	}
	var productType models.ProductType
	if err := database.DB.Where("id = ?", c.Param("id")).First(&productType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tipe produk tidak ditemukan"})
		return
	}
	var count int64
	if err := database.DB.Model(&models.ProductType{}).Where("LOWER(name) = LOWER(?) AND id <> ?", input.Name, productType.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa tipe produk"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Nama tipe produk sudah digunakan"})
		return
	}
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Where("product_type = ?", productType.Name).Update("product_type", input.Name).Error; err != nil {
			return err
		}
		return tx.Model(&productType).Update("name", input.Name).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui tipe produk"})
		return
	}
	productType.Name = input.Name
	c.JSON(http.StatusOK, productType)
}

func DeleteProductType(c *gin.Context) {
	var productType models.ProductType
	if err := database.DB.Where("id = ?", c.Param("id")).First(&productType).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tipe produk tidak ditemukan"})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Product{}).Where("product_type = ?", productType.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa penggunaan tipe produk"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Tipe produk masih digunakan produk"})
		return
	}
	if err := database.DB.Delete(&productType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus tipe produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tipe produk berhasil dihapus"})
}