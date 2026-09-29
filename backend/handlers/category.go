package handlers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
	"strings"

	"github.com/gin-gonic/gin"
)

// === CATEGORY HANDLERS ===

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.Name = strings.TrimSpace(cat.Name)
	if cat.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori wajib diisi"})
		return
	}
	var count int64
	database.DB.Model(&models.Category{}).Where("LOWER(name) = LOWER(?)", cat.Name).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Kategori sudah tersedia"})
		return
	}
	if err := database.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori"})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func UpdateCategory(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori wajib diisi"})
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" || len(input.Name) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama kategori wajib diisi (maksimal 100 karakter)"})
		return
	}

	var category models.Category
	if err := database.DB.First(&category, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	var count int64
	database.DB.Model(&models.Category{}).Where("LOWER(name) = LOWER(?) AND id <> ?", input.Name, category.ID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Nama kategori sudah digunakan"})
		return
	}
	if err := database.DB.Model(&category).Update("name", input.Name).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui kategori"})
		return
	}
	category.Name = input.Name
	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	var category models.Category
	if err := database.DB.First(&category, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	var productCount int64
	if err := database.DB.Model(&models.Product{}).Where("category_id = ?", category.ID).Count(&productCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa produk dalam kategori"})
		return
	}
	if productCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Kategori masih digunakan produk. Pindahkan atau hapus produknya terlebih dahulu."})
		return
	}
	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}
