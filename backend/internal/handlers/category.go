package handlers

import (
	"net/http"
	"pos-backend/internal/database"
	"pos-backend/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// === CATEGORY HANDLERS ===

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Where("parent_id IS NULL").Order("name").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var input struct {
		MerchantID string `json:"merchant_id"` // Required for multi-tenant
		Name       string `json:"name" binding:"required"`
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
	
	// Default to a placeholder if auth merchant parsing is not implemented yet
	var mID *uuid.UUID
	if input.MerchantID != "" {
		parsed, err := uuid.Parse(input.MerchantID)
		if err == nil {
			mID = &parsed
		}
	}

	var count int64
	if err := database.DB.Model(&models.Category{}).Where("LOWER(name) = LOWER(?) AND parent_id IS NULL", input.Name).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa kategori"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Kategori sudah tersedia"})
		return
	}
	category := models.Category{Name: input.Name, MerchantID: mID}
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan kategori"})
		return
	}
	c.JSON(http.StatusCreated, category)
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
	if err := database.DB.First(&category, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	if category.ParentID != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Kategori lama sudah dipisahkan. Pilih kategori utama untuk mengubahnya."})
		return
	}
	var count int64
	if err := database.DB.Model(&models.Category{}).Where("LOWER(name) = LOWER(?) AND id <> ? AND parent_id IS NULL", input.Name, category.ID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa nama kategori"})
		return
	}
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
	if err := database.DB.First(&category, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
		return
	}
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat subkategori"})
		return
	}
	
	categoryIDs := []uuid.UUID{category.ID}
	knownIDs := map[uuid.UUID]struct{}{category.ID: {}}
	
	for index := 0; index < len(categoryIDs); index++ {
		parentID := categoryIDs[index]
		for _, child := range categories {
			if child.ParentID == nil || *child.ParentID != parentID {
				continue
			}
			if _, exists := knownIDs[child.ID]; exists {
				continue
			}
			knownIDs[child.ID] = struct{}{}
			categoryIDs = append(categoryIDs, child.ID)
		}
	}
	var productCount int64
	if err := database.DB.Model(&models.Product{}).Where("category_id IN ?", categoryIDs).Count(&productCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa produk dalam kategori"})
		return
	}
	if productCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Kategori atau turunannya masih digunakan produk. Pindahkan atau hapus produknya terlebih dahulu."})
		return
	}
	if err := database.DB.Where("id IN ?", categoryIDs).Delete(&models.Category{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori dan seluruh turunannya berhasil dihapus"})
}
