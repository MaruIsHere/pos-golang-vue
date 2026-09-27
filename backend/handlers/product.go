package handlers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// === PRODUCT HANDLERS ===

func GetProducts(c *gin.Context) {
	var products []models.Product
	query := database.DB.Preload("Category")

	catID := c.Query("category_id")
	if catID != "" {
		query = query.Where("category_id = ?", catID)
	}

	if artist := c.Query("artist"); artist != "" {
		query = query.Where("artist = ?", artist)
	}

	if productType := c.Query("product_type"); productType != "" {
		query = query.Where("product_type = ?", productType)
	}

	search := c.Query("search")
	if search != "" {
		s := "%" + search + "%"
		query = query.Where("name LIKE ? OR barcode LIKE ? OR artist LIKE ? OR product_type LIKE ?", s, s, s, s)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func GetProductFilters(c *gin.Context) {
	var artists []string
	var productTypes []string

	database.DB.Model(&models.Product{}).Where("artist IS NOT NULL AND artist != ''").Distinct("artist").Pluck("artist", &artists)
	database.DB.Model(&models.Product{}).Where("product_type IS NOT NULL AND product_type != ''").Distinct("product_type").Pluck("product_type", &productTypes)

	c.JSON(http.StatusOK, gin.H{
		"artists":       artists,
		"product_types": productTypes,
	})
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	database.DB.Preload("Category").First(&product, product.ID)
	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&product)
	database.DB.Preload("Category").First(&product, product.ID)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}
