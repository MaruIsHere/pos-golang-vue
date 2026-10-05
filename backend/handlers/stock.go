package handlers

import (
	"fmt"
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// === STOCK MOVEMENT HANDLERS ===

func GetStockMovements(c *gin.Context) {
	var movements []models.StockMovement
	if err := database.DB.Preload("Product").Order("created_at desc").Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movements)
}

type CreateStockMovementInput struct {
	ProductID uint    `json:"product_id"`
	Type      string  `json:"type"` // "in" or "out"
	Quantity  float64 `json:"quantity"`
	Reason    string  `json:"reason"`
	Notes     string  `json:"notes"`
}

func CreateStockMovement(c *gin.Context) {
	var input CreateStockMovementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ProductID == 0 || input.Quantity <= 0 {
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

	movement := models.StockMovement{
		ProductID: input.ProductID,
		Type:      input.Type,
		Quantity:  input.Quantity,
		Unit:      prod.Unit,
		Reason:    input.Reason,
		Notes:     input.Notes,
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
