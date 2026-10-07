package handlers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// StoreWithStats Response Wrapper
type StoreWithStats struct {
	models.Store
	TotalCashiers int64 `json:"total_cashiers"`
	TotalProducts int64 `json:"total_products"`
}

// GetStores - Get list of all registered stores
func GetStores(c *gin.Context) {
	var stores []models.Store
	if err := database.DB.Order("created_at desc").Find(&stores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data toko"})
		return
	}

	var result []StoreWithStats
	for _, s := range stores {
		var cashierCount int64
		var productCount int64
		database.DB.Model(&models.User{}).Where("store_id = ?", s.ID).Count(&cashierCount)
		database.DB.Model(&models.Product{}).Where("store_id = ?", s.ID).Count(&productCount)

		result = append(result, StoreWithStats{
			Store:         s,
			TotalCashiers: cashierCount,
			TotalProducts: productCount,
		})
	}

	c.JSON(http.StatusOK, result)
}

// CreateStorePayload Input structure
type CreateStorePayload struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

// CreateStore - Register a new Store (Registrasi Toko Baru)
func CreateStore(c *gin.Context) {
	var payload CreateStorePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama toko wajib diisi"})
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama toko tidak boleh kosong"})
		return
	}

	code := strings.TrimSpace(payload.Code)
	if code == "" {
		var count int64
		database.DB.Model(&models.Store{}).Count(&count)
		code = "STORE-" + strconv.FormatInt(count+1, 10)
	}

	// Check unique code
	var existing int64
	database.DB.Model(&models.Store{}).Where("LOWER(code) = LOWER(?)", code).Count(&existing)
	if existing > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Kode toko sudah terdaftar, gunakan kode lain"})
		return
	}

	store := models.Store{
		Name:     payload.Name,
		Code:     code,
		Address:  payload.Address,
		Phone:    payload.Phone,
		IsActive: true,
	}

	if err := database.DB.Create(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan toko baru: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, store)
}

// UpdateStore - Update Store Details
func UpdateStore(c *gin.Context) {
	id := c.Param("id")
	var store models.Store
	if err := database.DB.First(&store, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	var payload CreateStorePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	store.Name = strings.TrimSpace(payload.Name)
	if payload.Code != "" {
		store.Code = strings.TrimSpace(payload.Code)
	}
	store.Address = payload.Address
	store.Phone = payload.Phone

	if err := database.DB.Save(&store).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meperbarui data toko"})
		return
	}

	c.JSON(http.StatusOK, store)
}

// DeleteStore - Deactivate or Delete Store
func DeleteStore(c *gin.Context) {
	id := c.Param("id")
	var store models.Store
	if err := database.DB.First(&store, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	// Deactivate store
	if err := database.DB.Model(&store).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menonaktifkan toko"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Toko berhasil dinonaktifkan"})
}

// AssignUserStorePayload
type AssignUserStorePayload struct {
	UserID  uint  `json:"user_id" binding:"required"`
	StoreID *uint `json:"store_id"`
}

// AssignUserStore - Assign user/cashier to a store
func AssignUserStore(c *gin.Context) {
	var payload AssignUserStorePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID wajib diisi"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, payload.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if err := database.DB.Model(&user).Update("store_id", payload.StoreID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menetapkan toko pada kasir"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kasir berhasil ditugaskan ke toko"})
}

// ImportMasterProductPayload
type ImportMasterProductPayload struct {
	MasterProductID uint     `json:"master_product_id" binding:"required"`
	StoreID         uint     `json:"store_id" binding:"required"`
	Stock           *float64 `json:"stock"`
	Price           *float64 `json:"price"`
}

// ImportMasterProductToStore - Add Master Product to Store POS catalog
func ImportMasterProductToStore(c *gin.Context) {
	var payload ImportMasterProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Master product ID dan Store ID wajib diisi"})
		return
	}

	var masterProd models.Product
	if err := database.DB.First(&masterProd, payload.MasterProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Master Produk tidak ditemukan"})
		return
	}

	var store models.Store
	if err := database.DB.First(&store, payload.StoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	// Check if store already has a store product (is_master = false) from this master product
	var existing models.Product
	err := database.DB.Where("store_id = ? AND is_master = ? AND master_product_id = ?", payload.StoreID, false, masterProd.ID).First(&existing).Error
	if err == nil {
		// Update existing store product stock & price if provided
		if payload.Stock != nil {
			existing.Stock = *payload.Stock
		}
		if payload.Price != nil {
			existing.Price = *payload.Price
		}
		existing.IsActive = true
		database.DB.Save(&existing)
		c.JSON(http.StatusOK, existing)
		return
	}

	stockVal := masterProd.Stock
	if payload.Stock != nil {
		stockVal = *payload.Stock
	}
	priceVal := masterProd.Price
	if payload.Price != nil {
		priceVal = *payload.Price
	}

	masterID := masterProd.ID
	storeProd := models.Product{
		CategoryID:      masterProd.CategoryID,
		Name:            masterProd.Name,
		Artist:          masterProd.Artist,
		ProductType:     masterProd.ProductType,
		Price:           priceVal,
		CostPrice:       masterProd.CostPrice,
		Stock:           stockVal,
		Unit:            masterProd.Unit,
		Barcode:         masterProd.Barcode,
		ImageURL:        masterProd.ImageURL,
		IsActive:        true,
		IsMaster:        false,
		MasterProductID: &masterID,
		StoreID:         &payload.StoreID,
	}

	if err := database.DB.Create(&storeProd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan produk ke kasir toko"})
		return
	}

	c.JSON(http.StatusCreated, storeProd)
}

// ImportBatchPayload
type ImportBatchPayload struct {
	StoreID          uint   `json:"store_id" binding:"required"`
	MasterProductIDs []uint `json:"master_product_ids"`
}

// ImportBatchMasterProductsToStore - Copy multiple or all Master Products to a Store POS catalog
func ImportBatchMasterProductsToStore(c *gin.Context) {
	var payload ImportBatchPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Store ID wajib diisi"})
		return
	}

	var store models.Store
	if err := database.DB.First(&store, payload.StoreID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Toko tidak ditemukan"})
		return
	}

	var masterProducts []models.Product
	query := database.DB.Where("is_master = ?", true)
	if len(payload.MasterProductIDs) > 0 {
		query = query.Where("id IN ?", payload.MasterProductIDs)
	}

	if err := query.Find(&masterProducts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk utama"})
		return
	}

	importedCount := 0
	for _, masterProd := range masterProducts {
		var existingCount int64
		database.DB.Model(&models.Product{}).Where("store_id = ? AND is_master = ? AND master_product_id = ?", payload.StoreID, false, masterProd.ID).Count(&existingCount)
		if existingCount > 0 {
			continue
		}

		masterID := masterProd.ID
		storeProd := models.Product{
			CategoryID:      masterProd.CategoryID,
			Name:            masterProd.Name,
			Artist:          masterProd.Artist,
			ProductType:     masterProd.ProductType,
			Price:           masterProd.Price,
			CostPrice:       masterProd.CostPrice,
			Stock:           masterProd.Stock,
			Unit:            masterProd.Unit,
			Barcode:         masterProd.Barcode,
			ImageURL:        masterProd.ImageURL,
			IsActive:        true,
			IsMaster:        false,
			MasterProductID: &masterID,
			StoreID:         &payload.StoreID,
		}
		if err := database.DB.Create(&storeProd).Error; err == nil {
			importedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengimpor data produk utama ke kasir toko",
		"count":   importedCount,
	})
}
