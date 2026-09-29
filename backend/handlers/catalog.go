package handlers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
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

func GetArtists(c *gin.Context) {
	var artists []models.Artist
	if err := database.DB.Order("name").Find(&artists).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar artist"})
		return
	}
	c.JSON(http.StatusOK, artists)
}

func CreateArtist(c *gin.Context) {
	input, ok := normalizeCatalogInput(c)
	if !ok {
		return
	}
	var count int64
	if err := database.DB.Model(&models.Artist{}).Where("LOWER(name) = LOWER(?)", input.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa nama artist"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Artist sudah tersedia"})
		return
	}
	artist := models.Artist{Name: input.Name}
	if err := database.DB.Create(&artist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan artist"})
		return
	}
	c.JSON(http.StatusCreated, artist)
}

func UpdateArtist(c *gin.Context) {
	input, ok := normalizeCatalogInput(c)
	if !ok {
		return
	}
	var artist models.Artist
	if err := database.DB.First(&artist, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist tidak ditemukan"})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Artist{}).Where("LOWER(name) = LOWER(?) AND id <> ?", input.Name, artist.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa nama artist"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Nama artist sudah digunakan"})
		return
	}
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Where("artist = ?", artist.Name).Update("artist", input.Name).Error; err != nil {
			return err
		}
		return tx.Model(&artist).Update("name", input.Name).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui artist"})
		return
	}
	artist.Name = input.Name
	c.JSON(http.StatusOK, artist)
}

func DeleteArtist(c *gin.Context) {
	var artist models.Artist
	if err := database.DB.First(&artist, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist tidak ditemukan"})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Product{}).Where("artist = ?", artist.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa penggunaan artist"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Artist masih digunakan produk"})
		return
	}
	if err := database.DB.Delete(&artist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus artist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Artist berhasil dihapus"})
}

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
	if err := database.DB.First(&productType, c.Param("id")).Error; err != nil {
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
	if err := database.DB.First(&productType, c.Param("id")).Error; err != nil {
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