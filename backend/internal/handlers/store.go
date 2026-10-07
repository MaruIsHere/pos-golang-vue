package handlers

import (
	"net/http"
	"pos-backend/internal/database"
	"pos-backend/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// OutletWithStats Response Wrapper
type OutletWithStats struct {
	models.Outlet
	TotalCashiers int64 `json:"total_cashiers"`
	TotalProducts int64 `json:"total_products"`
}

// GetStores - Get list of all registered outlets (lapak/gudang)
func GetStores(c *gin.Context) {
	var outlets []models.Outlet
	if err := database.DB.Order("created_at desc").Find(&outlets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data lapak/outlet"})
		return
	}

	var result []OutletWithStats
	for _, s := range outlets {
		var cashierCount int64
		var productCount int64
		database.DB.Model(&models.User{}).Where("outlet_id = ?", s.ID).Count(&cashierCount)
		database.DB.Model(&models.Product{}).Where("outlet_id = ?", s.ID).Count(&productCount)

		result = append(result, OutletWithStats{
			Outlet:        s,
			TotalCashiers: cashierCount,
			TotalProducts: productCount,
		})
	}

	c.JSON(http.StatusOK, result)
}

// CreateOutletPayload Input structure
type CreateOutletPayload struct {
	MerchantID  string `json:"merchant_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	IsWarehouse bool   `json:"is_warehouse"`
}

// CreateStore - Register a new Outlet
func CreateStore(c *gin.Context) {
	var payload CreateOutletPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid: " + err.Error()})
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama outlet tidak boleh kosong"})
		return
	}

	merchantID, err := uuid.Parse(payload.MerchantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID tidak valid"})
		return
	}

	code := strings.TrimSpace(payload.Code)
	if code == "" {
		code = "OUTLET-" + uuid.New().String()[:8]
	}

	var existing int64
	database.DB.Model(&models.Outlet{}).Where("LOWER(code) = LOWER(?) AND merchant_id = ?", code, merchantID).Count(&existing)
	if existing > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Kode outlet sudah terdaftar di merchant ini"})
		return
	}

	outlet := models.Outlet{
		MerchantID:  merchantID,
		Name:        payload.Name,
		Code:        code,
		Address:     payload.Address,
		Phone:       payload.Phone,
		IsWarehouse: payload.IsWarehouse,
		IsActive:    true,
	}

	if err := database.DB.Create(&outlet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan outlet baru: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, outlet)
}

// UpdateStore - Update Outlet Details
func UpdateStore(c *gin.Context) {
	id := c.Param("id")
	var outlet models.Outlet
	if err := database.DB.First(&outlet, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
		return
	}

	var payload CreateOutletPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	outlet.Name = strings.TrimSpace(payload.Name)
	if payload.Code != "" {
		outlet.Code = strings.TrimSpace(payload.Code)
	}
	outlet.Address = payload.Address
	outlet.Phone = payload.Phone
	outlet.IsWarehouse = payload.IsWarehouse

	if err := database.DB.Save(&outlet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal meperbarui data outlet"})
		return
	}

	c.JSON(http.StatusOK, outlet)
}

// DeleteStore - Deactivate Outlet
func DeleteStore(c *gin.Context) {
	id := c.Param("id")
	var outlet models.Outlet
	if err := database.DB.First(&outlet, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
		return
	}

	if err := database.DB.Model(&outlet).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menonaktifkan outlet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Outlet berhasil dinonaktifkan"})
}

// AssignUserOutletPayload
type AssignUserOutletPayload struct {
	UserID   string `json:"user_id" binding:"required"`
	OutletID string `json:"outlet_id"` // Using string for UUID
}

// AssignUserStore - Assign cashier to outlet
func AssignUserStore(c *gin.Context) {
	var payload AssignUserOutletPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID wajib diisi"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", payload.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	
	if payload.OutletID != "" {
		outletID, err := uuid.Parse(payload.OutletID)
		if err == nil {
			user.OutletID = &outletID
		}
	} else {
		user.OutletID = nil
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menetapkan outlet pada kasir"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kasir berhasil ditugaskan ke outlet"})
}

// ImportMasterProductPayload
type ImportMasterProductPayload struct {
	MasterProductID string   `json:"master_product_id" binding:"required"`
	OutletID        string   `json:"outlet_id" binding:"required"`
	Stock           *float64 `json:"stock"`
	Price           *float64 `json:"price"`
}

// ImportMasterProductToStore - Add Master Product to Outlet catalog
func ImportMasterProductToStore(c *gin.Context) {
	var payload ImportMasterProductPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Master product ID dan Outlet ID wajib diisi"})
		return
	}

	var masterProd models.Product
	if err := database.DB.First(&masterProd, "id = ?", payload.MasterProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Master Produk tidak ditemukan"})
		return
	}

	var outlet models.Outlet
	if err := database.DB.First(&outlet, "id = ?", payload.OutletID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
		return
	}

	var existing models.Product
	err := database.DB.Where("outlet_id = ? AND is_master = ? AND master_product_id = ?", outlet.ID, false, masterProd.ID).First(&existing).Error
	if err == nil {
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
	outletID := outlet.ID
	
	storeProd := models.Product{
		MerchantID:      outlet.MerchantID,
		OutletID:        &outletID,
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
	}

	if err := database.DB.Create(&storeProd).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan produk ke kasir lapak"})
		return
	}

	c.JSON(http.StatusCreated, storeProd)
}

// ImportBatchPayload
type ImportBatchPayload struct {
	OutletID         string   `json:"outlet_id" binding:"required"`
	MasterProductIDs []string `json:"master_product_ids"`
}

// ImportBatchMasterProductsToStore
func ImportBatchMasterProductsToStore(c *gin.Context) {
	var payload ImportBatchPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Outlet ID wajib diisi"})
		return
	}

	var outlet models.Outlet
	if err := database.DB.First(&outlet, "id = ?", payload.OutletID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Outlet tidak ditemukan"})
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
		database.DB.Model(&models.Product{}).Where("outlet_id = ? AND is_master = ? AND master_product_id = ?", outlet.ID, false, masterProd.ID).Count(&existingCount)
		if existingCount > 0 {
			continue
		}

		masterID := masterProd.ID
		outletID := outlet.ID
		storeProd := models.Product{
			MerchantID:      outlet.MerchantID,
			OutletID:        &outletID,
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
		}
		if err := database.DB.Create(&storeProd).Error; err == nil {
			importedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengimpor data produk ke lapak",
		"count":   importedCount,
	})
}
