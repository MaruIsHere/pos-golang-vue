package handlers

import (
	"fmt"
	"net/http"
	"pos-backend/internal/database"
	"pos-backend/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// === STOCK MOVEMENT HANDLERS ===

func GetStockMovements(c *gin.Context) {
	var movements []models.StockMovement
	query := database.DB.Preload("Product").Order("created_at desc")

	if merchantID, exists := c.Get("merchant_id"); exists {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if outletID, exists := c.Get("outlet_id"); exists {
		query = query.Where("outlet_id = ?", outletID)
	}

	if err := query.Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movements)
}

type CreateStockMovementInput struct {
	MerchantID *uuid.UUID `json:"merchant_id"`
	OutletID   *uuid.UUID `json:"outlet_id"`
	ProductID  uuid.UUID  `json:"product_id"`
	Type       string     `json:"type"` // "in" or "out"
	Quantity   float64    `json:"quantity"`
	Reason     string     `json:"reason"`
	Notes      string     `json:"notes"`
}

func CreateStockMovement(c *gin.Context) {
	var input CreateStockMovementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ProductID == uuid.Nil || input.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Produk dan kuantitas harus valid"})
		return
	}

	if input.Type != "in" && input.Type != "out" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe mutasi stok harus 'in' atau 'out'"})
		return
	}

	tx := database.DB.Begin()

	var prod models.Product
	if err := tx.First(&prod, input.ProductID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	if !isValidStockQuantity(input.Quantity, prod.Unit) {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas harus bilangan bulat untuk pcs dan maksimal 3 angka desimal untuk gram/liter"})
		return
	}

	if input.Type == "in" {
		prod.Stock += input.Quantity
	} else {
		if prod.Stock < input.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok %s tidak mencukupi (sisa %s %s)", prod.Name, strconv.FormatFloat(prod.Stock, 'f', -1, 64), prod.Unit)})
			return
		}
		prod.Stock -= input.Quantity
	}

	if err := tx.Save(&prod).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parsing IDs from input or context
	merchantID := prod.MerchantID
	if input.MerchantID != nil {
		merchantID = *input.MerchantID
	}

	var outletID uuid.UUID
	if prod.OutletID != nil {
		outletID = *prod.OutletID
	} else if input.OutletID != nil {
		outletID = *input.OutletID
	} else if oidRaw, exists := c.Get("outlet_id"); exists {
		outletID, _ = uuid.Parse(fmt.Sprintf("%v", oidRaw))
	}

	movement := models.StockMovement{
		MerchantID: merchantID,
		OutletID:   outletID,
		ProductID:  input.ProductID,
		Type:       input.Type,
		Quantity:   input.Quantity,
		Unit:       prod.Unit,
		Reason:     input.Reason,
		Notes:      input.Notes,
	}

	if err := tx.Create(&movement).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	database.DB.Preload("Product").First(&movement, movement.ID)
	c.JSON(http.StatusCreated, movement)
}
