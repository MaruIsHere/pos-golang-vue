package handlers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
	"strings"

	"github.com/gin-gonic/gin"
)

type productInput struct {
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name"`
	Artist      string  `json:"artist"`
	ProductType string  `json:"product_type"`
	Price       float64 `json:"price"`
	CostPrice   float64 `json:"cost_price"`
	Stock       int     `json:"stock"`
	Barcode     string  `json:"barcode"`
	ImageURL    string  `json:"image_url"`
	IsActive    *bool   `json:"is_active"`
}

func validateProductInput(input *productInput) (string, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Artist = strings.TrimSpace(input.Artist)
	input.ProductType = strings.TrimSpace(input.ProductType)
	if input.Name == "" {
		return "Nama produk wajib diisi", nil
	}
	if input.CategoryID == 0 {
		return "Kategori produk wajib dipilih", nil
	}
	if input.Price < 0 || input.CostPrice < 0 || input.Stock < 0 {
		return "Harga dan stok tidak boleh negatif", nil
	}
	var categoryCount int64
	if err := database.DB.Model(&models.Category{}).Where("id = ?", input.CategoryID).Count(&categoryCount).Error; err != nil {
		return "", err
	}
	if categoryCount == 0 {
		return "Kategori produk tidak ditemukan", nil
	}
	if input.Artist != "" {
		var artistCount int64
		if err := database.DB.Model(&models.Artist{}).Where("LOWER(name) = LOWER(?)", input.Artist).Count(&artistCount).Error; err != nil {
			return "", err
		}
		if artistCount == 0 {
			return "Artist tidak ditemukan. Tambahkan artist melalui master artist terlebih dahulu", nil
		}
	}
	if input.ProductType != "" {
		var typeCount int64
		if err := database.DB.Model(&models.ProductType{}).Where("LOWER(name) = LOWER(?)", input.ProductType).Count(&typeCount).Error; err != nil {
			return "", err
		}
		if typeCount == 0 {
			return "Tipe produk tidak ditemukan. Tambahkan tipe produk melalui master terlebih dahulu", nil
		}
	}
	return "", nil
}

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

	database.DB.Model(&models.Artist{}).Order("name").Pluck("name", &artists)
	database.DB.Model(&models.ProductType{}).Order("name").Pluck("name", &productTypes)

	c.JSON(http.StatusOK, gin.H{
		"artists":       artists,
		"product_types": productTypes,
	})
}

func CreateProduct(c *gin.Context) {
	var input productInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if message, err := validateProductInput(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa kategori produk"})
		return
	} else if message != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}
	product := models.Product{
		CategoryID: input.CategoryID,
		Name:       input.Name,
		Artist:     input.Artist,
		ProductType: input.ProductType,
		Price:      input.Price,
		CostPrice:  input.CostPrice,
		Stock:      input.Stock,
		Barcode:    input.Barcode,
		ImageURL:   input.ImageURL,
	}
	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	}
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk"})
		return
	}
	if err := database.DB.Preload("Category").First(&product, product.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Produk tersimpan, tetapi gagal memuat kategorinya"})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var input productInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if message, err := validateProductInput(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa kategori produk"})
		return
	} else if message != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": message})
		return
	}
	updates := map[string]interface{}{
		"category_id":  input.CategoryID,
		"name":         input.Name,
		"artist":       input.Artist,
		"product_type": input.ProductType,
		"price":        input.Price,
		"cost_price":   input.CostPrice,
		"stock":        input.Stock,
		"barcode":      input.Barcode,
		"image_url":    input.ImageURL,
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if err := database.DB.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui produk"})
		return
	}
	if err := database.DB.Preload("Category").First(&product, product.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Produk diperbarui, tetapi gagal memuat kategorinya"})
		return
	}
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
