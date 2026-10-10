package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"pos-backend/internal/database"
	"pos-backend/internal/models"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// === ORDER / TRANSAKSI HANDLERS ===

type CreateOrderInput struct {
	MerchantID       *uuid.UUID        `json:"merchant_id"`
	OutletID         *uuid.UUID        `json:"outlet_id"`
	CustomerName     string            `json:"customer_name"`
	TableNumber      string            `json:"table_number"`
	CashierName      string            `json:"cashier_name"`
	PaymentMethod    string            `json:"payment_method"`
	PaymentProof     string            `json:"payment_proof"`
	PaymentReference string            `json:"payment_reference"`
	PlatformFee      float64           `json:"platform_fee"`
	PaidAmount       float64           `json:"paid_amount"`
	Discount         float64           `json:"discount"`
	Tax              float64           `json:"tax"`
	Items            []CreateOrderItem `json:"items"`
}

type CreateOrderItem struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  float64   `json:"quantity"`
	Notes     string    `json:"notes"`
}

func UploadPaymentProof(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10*1024*1024)
	if err := c.Request.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran foto bukti pembayaran terlalu besar (maksimal 10 MB)"})
		return
	}
	if c.Request.MultipartForm != nil {
		defer c.Request.MultipartForm.RemoveAll()
	}

	file, header, err := c.Request.FormFile("proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pilih file bukti pembayaran terlebih dahulu"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal membaca bukti pembayaran"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		ext = ".jpg"
	}

	var fileID [16]byte
	if _, err := rand.Read(fileID[:]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat nama file"})
		return
	}
	const uploadDir = "uploads/proofs"
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan direktori bukti pembayaran"})
		return
	}
	filename := fmt.Sprintf("proof_%s%s", hex.EncodeToString(fileID[:]), ext)
	filePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(filePath, data, 0640); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file bukti pembayaran"})
		return
	}

	// Mendeteksi OS untuk memilih file binary OCR yang tepat
	binaryName := "ocrs"
	if runtime.GOOS == "windows" {
		binaryName = "ocrs.exe"
	}
	
	// Coba cari di folder bin/ pada root project (jika run dari root atau backend/)
	ocrsPath := filepath.Join("..", "bin", binaryName)
	if _, err := os.Stat(ocrsPath); os.IsNotExist(err) {
		ocrsPath = filepath.Join(".", "bin", binaryName)
		if _, err := os.Stat(ocrsPath); os.IsNotExist(err) {
			// Fallback ke ~/.cargo/bin/ocrs untuk development lokal
			homeDir, _ := os.UserHomeDir()
			ocrsPath = filepath.Join(homeDir, ".cargo", "bin", "ocrs")
		}
	}

	var detectedText string
	var detectedAmount float64

	cmd := exec.Command(ocrsPath, filePath)
	out, err := cmd.CombinedOutput()
	if err == nil {
		detectedText = string(out)

		// Regex sederhana untuk mencari nominal uang
		// Contoh: Rp 50.000 atau Rp50.000 atau sekadar 50000
		// Ini adalah implementasi awal yang bisa disempurnakan nanti
		re := regexp.MustCompile(`(?i)(?:rp|idr)\s*\.?\s*([0-9.,]+)`)
		matches := re.FindAllStringSubmatch(detectedText, -1)

		var maxAmount float64
		for _, match := range matches {
			if len(match) > 1 {
				// Bersihkan titik dan koma
				numStr := strings.ReplaceAll(match[1], ".", "")
				numStr = strings.ReplaceAll(numStr, ",", "")

				if val, err := strconv.ParseFloat(numStr, 64); err == nil {
					// Cari nilai terbesar (biasanya total pembayaran)
					if val > maxAmount {
						maxAmount = val
					}
				}
			}
		}
		detectedAmount = maxAmount
	} else {
		// Log error jika OCR gagal, tapi jangan gagalkan upload
		fmt.Printf("OCR Error: %v, Output: %s\n", err, string(out))
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment_proof":   fmt.Sprintf("/uploads/proofs/%s", filename),
		"detected_amount": detectedAmount,
		"ocr_raw_text":    detectedText, // Bisa dihapus nanti kalau sudah tidak butuh debugging
	})
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
		if itemInput.Quantity <= 0 || math.Abs(itemInput.Quantity*1000-math.Round(itemInput.Quantity*1000)) > 1e-7 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas produk harus lebih dari 0 dan maksimal 3 angka desimal"})
			return
		}

		var prod models.Product
		if err := tx.First(&prod, itemInput.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Produk ID %s tidak ditemukan", itemInput.ProductID)})
			return
		}
		if !isValidStockQuantity(itemInput.Quantity, prod.Unit) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas harus bilangan bulat untuk pcs dan maksimal 3 angka desimal untuk gram/liter"})
			return
		}

		if prod.Stock < itemInput.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok produk %s tidak mencukupi (sisa %s %s)", prod.Name, strconv.FormatFloat(prod.Stock, 'f', -1, 64), prod.Unit)})
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
			Artist:       prod.Artist,
			ProductType:  prod.ProductType,
			ProductPrice: prod.Price,
			Quantity:     itemInput.Quantity,
			Unit:         prod.Unit,
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

	cashierName := input.CashierName
	if cashierName == "" {
		if userID, exists := c.Get("user_id"); exists {
			var u models.User
			if err := database.DB.Where("id = ?", userID).First(&u).Error; err == nil {
				if u.Name != "" {
					cashierName = u.Name
				} else if u.Username != "" {
					cashierName = u.Username
				}
			}
		}
	}
	if cashierName == "" {
		cashierName = "Kasir"
	}

	merchantID := uuid.Nil
	if input.MerchantID != nil {
		merchantID = *input.MerchantID
	} else if mid, exists := c.Get("merchant_id"); exists {
		merchantID, _ = uuid.Parse(fmt.Sprintf("%v", mid))
	}

	outletID := uuid.Nil
	if input.OutletID != nil {
		outletID = *input.OutletID
	} else if oid, exists := c.Get("outlet_id"); exists {
		outletID, _ = uuid.Parse(fmt.Sprintf("%v", oid))
	}

	order := models.Order{
		MerchantID:       merchantID,
		OutletID:         outletID,
		InvoiceNo:        invoiceNo,
		TotalAmount:      totalAmount,
		Discount:         input.Discount,
		Tax:              input.Tax,
		GrandTotal:       grandTotal,
		PaidAmount:       input.PaidAmount,
		ChangeAmount:     changeAmount,
		PaymentMethod:    input.PaymentMethod,
		PaymentProof:     input.PaymentProof,
		PaymentReference: input.PaymentReference,
		PlatformFee:      input.PlatformFee,
		Status:           "completed",
		CashierName:      cashierName,
		CustomerName:     custName,
		TableNumber:      input.TableNumber,
		OrderItems:       orderItems,
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

	// Multi-tenant isolation
	if merchantID, exists := c.Get("merchant_id"); exists {
		query = query.Where("merchant_id = ?", merchantID)
	}
	if outletID, exists := c.Get("outlet_id"); exists {
		query = query.Where("outlet_id = ?", outletID)
	}

	// Optional filtering from query params
	if qMerchantID := c.Query("merchant_id"); qMerchantID != "" {
		query = query.Where("merchant_id = ?", qMerchantID)
	}
	if qOutletID := c.Query("outlet_id"); qOutletID != "" {
		query = query.Where("outlet_id = ?", qOutletID)
	}

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
			MerchantID: order.MerchantID,
			OutletID:   order.OutletID,
			ProductID:  item.ProductID,
			Type:       "in",
			Quantity:   item.Quantity,
			Unit:       item.Unit,
			Reason:     "retur_penjualan",
			Notes:      fmt.Sprintf("Retur Transaksi Invoice #%s", order.InvoiceNo),
		}
		tx.Create(&movement)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil diretur dan stok telah dipulihkan",
		"order":   order,
	})
}
