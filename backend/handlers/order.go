package handlers

import (
	"fmt"
	"net/http"
	"pos-backend/database"
	"pos-backend/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// === ORDER / TRANSAKSI HANDLERS ===

type CreateOrderInput struct {
	CustomerName  string            `json:"customer_name"`
	PaymentMethod string            `json:"payment_method"`
	PaidAmount    float64           `json:"paid_amount"`
	Discount      float64           `json:"discount"`
	Tax           float64           `json:"tax"`
	Items         []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Notes     string `json:"notes"`
}

func CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang belanja tidak boleh kosong"})
		return
	}

	tx := database.DB.Begin()

	var totalAmount float64
	var orderItems []models.OrderItem

	for _, itemInput := range input.Items {
		var prod models.Product
		if err := tx.First(&prod, itemInput.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Produk ID %d tidak ditemukan", itemInput.ProductID)})
			return
		}

		if prod.Stock < itemInput.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok produk %s tidak mencukupi (sisa %d)", prod.Name, prod.Stock)})
			return
		}

		// Reduce Stock
		prod.Stock -= itemInput.Quantity
		if err := tx.Save(&prod).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		subtotal := prod.Price * float64(itemInput.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:    prod.ID,
			ProductName:  prod.Name,
			ProductPrice: prod.Price,
			Quantity:     itemInput.Quantity,
			Subtotal:     subtotal,
			Notes:        itemInput.Notes,
		})
	}

	grandTotal := (totalAmount - input.Discount) + input.Tax
	if grandTotal < 0 {
		grandTotal = 0
	}

	changeAmount := input.PaidAmount - grandTotal
	if changeAmount < 0 {
		changeAmount = 0
	}

	invoiceNo := fmt.Sprintf("INV-%s-%d", time.Now().Format("20060102-150405"), time.Now().Nanosecond()%1000)

	custName := input.CustomerName
	if custName == "" {
		custName = "Umum"
	}

	order := models.Order{
		InvoiceNo:     invoiceNo,
		TotalAmount:   totalAmount,
		Discount:      input.Discount,
		Tax:           input.Tax,
		GrandTotal:    grandTotal,
		PaidAmount:    input.PaidAmount,
		ChangeAmount:  changeAmount,
		PaymentMethod: input.PaymentMethod,
		Status:        "completed",
		CashierName:   "Kasir 1",
		CustomerName:  custName,
		OrderItems:    orderItems,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	database.DB.Preload("OrderItems").First(&order, order.ID)
	c.JSON(http.StatusCreated, order)
}

func GetOrders(c *gin.Context) {
	var orders []models.Order
	query := database.DB.Preload("OrderItems").Order("created_at desc")

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	query = query.Limit(limit)

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func GetOrderById(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("OrderItems").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// === REFUND ORDER HANDLER ===

func RefundOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order
	if err := database.DB.Preload("OrderItems").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	if order.Status == "refunded" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaksi ini sudah diretur sebelumnya"})
		return
	}

	tx := database.DB.Begin()

	// Update status
	order.Status = "refunded"
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Restore product stocks & record stock movements
	for _, item := range order.OrderItems {
		var prod models.Product
		if err := tx.First(&prod, item.ProductID).Error; err == nil {
			prod.Stock += item.Quantity
			tx.Save(&prod)
		}

		movement := models.StockMovement{
			ProductID: item.ProductID,
			Type:      "in",
			Quantity:  item.Quantity,
			Reason:    "retur_penjualan",
			Notes:     fmt.Sprintf("Retur Transaksi Invoice #%s", order.InvoiceNo),
		}
		tx.Create(&movement)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil diretur dan stok telah dipulihkan",
		"order":   order,
	})
}
