package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"pos-backend/config"
	"pos-backend/database"
	"pos-backend/models"

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
	if err := database.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
}

// === PRODUCT HANDLERS ===

func GetProducts(c *gin.Context) {
	var products []models.Product
	query := database.DB.Preload("Category")

	catID := c.Query("category_id")
	if catID != "" {
		query = query.Where("category_id = ?", catID)
	}

	search := c.Query("search")
	if search != "" {
		query = query.Where("name LIKE ? OR barcode LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
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

// === ORDER / TRANSAKSI HANDLERS ===

type CreateOrderInput struct {
	CustomerName  string             `json:"customer_name"`
	PaymentMethod string             `json:"payment_method"`
	PaidAmount    float64            `json:"paid_amount"`
	Discount      float64            `json:"discount"`
	Tax           float64            `json:"tax"`
	Items         []CreateOrderItem  `json:"items"`
}

type CreateOrderItem struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Notes     string  `json:"notes"`
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

// === REPORT & DASHBOARD HANDLERS ===

func GetDashboardStats(c *gin.Context) {
	var totalOrders int64
	var totalRevenue float64
	var totalItemsSold int64

	database.DB.Model(&models.Order{}).Where("status = ?", "completed").Count(&totalOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", "completed").Select("COALESCE(SUM(grand_total), 0)").Scan(&totalRevenue)
	database.DB.Model(&models.OrderItem{}).Select("COALESCE(SUM(quantity), 0)").Scan(&totalItemsSold)

	// Top Selling Products
	type TopProduct struct {
		ProductName string  `json:"product_name"`
		TotalQty    int     `json:"total_qty"`
		TotalSales  float64 `json:"total_sales"`
	}
	var topProducts []TopProduct
	database.DB.Table("order_items").
		Select("product_name, SUM(quantity) as total_qty, SUM(subtotal) as total_sales").
		Group("product_name").
		Order("total_qty desc").
		Limit(5).
		Scan(&topProducts)

	// Recent Orders
	var recentOrders []models.Order
	database.DB.Order("created_at desc").Limit(5).Find(&recentOrders)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":     totalOrders,
		"total_revenue":    totalRevenue,
		"total_items_sold": totalItemsSold,
		"top_products":     topProducts,
		"recent_orders":    recentOrders,
	})
}

// === STORE & DATABASE SETTINGS HANDLERS ===

func GetSettings(c *gin.Context) {
	var setting models.StoreSetting
	if err := database.DB.First(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, setting)
}

func UpdateSettings(c *gin.Context) {
	var setting models.StoreSetting
	if err := database.DB.First(&setting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input models.StoreSetting
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setting.StoreName = input.StoreName
	setting.Address = input.Address
	setting.Phone = input.Phone
	setting.ReceiptFooter = input.ReceiptFooter
	setting.TaxPercentage = input.TaxPercentage
	setting.QrisImageUrl = input.QrisImageUrl

	database.DB.Save(&setting)
	c.JSON(http.StatusOK, setting)
}

// === VOUCHER HANDLERS ===

func GetVouchers(c *gin.Context) {
	var vouchers []models.Voucher
	if err := database.DB.Order("created_at desc").Find(&vouchers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vouchers)
}

func CreateVoucher(c *gin.Context) {
	var voucher models.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if voucher.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kode voucher tidak boleh kosong"})
		return
	}

	voucher.IsActive = true
	if err := database.DB.Create(&voucher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat voucher (kode mungkin sudah ada)"})
		return
	}
	c.JSON(http.StatusCreated, voucher)
}

func DeleteVoucher(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Voucher{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voucher berhasil dihapus"})
}

type SwitchDbInput struct {
	Engine   string `json:"engine"`    // "sqlite" or "mysql"
	MysqlDsn string `json:"mysql_dsn"` // optional DSN
}

func SwitchDatabase(c *gin.Context) {
	var input SwitchDbInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Engine != "sqlite" && input.Engine != "mysql" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Engine database harus 'sqlite' atau 'mysql'"})
		return
	}

	cfg := config.AppConfig
	cfg.DbEngine = input.Engine
	if input.MysqlDsn != "" {
		cfg.MysqlDsn = input.MysqlDsn
	}

	newDB, err := database.InitDB(cfg.DbEngine, cfg.MysqlDsn, cfg.SqlitePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Gagal berpindah ke %s: %v", input.Engine, err)})
		return
	}

	database.DB = newDB
	config.SaveConfig(cfg)

	// Update setting record
	var setting models.StoreSetting
	if err := database.DB.First(&setting).Error; err == nil {
		setting.DbEngine = cfg.DbEngine
		setting.MysqlDsn = cfg.MysqlDsn
		database.DB.Save(&setting)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("Berhasil berpindah ke basis data %s", input.Engine),
		"db_engine": input.Engine,
	})
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

// === CUSTOMER HANDLERS ===

func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	if err := database.DB.Order("created_at desc").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func CreateCustomer(c *gin.Context) {
	var cust models.Customer
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if cust.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama pelanggan tidak boleh kosong"})
		return
	}
	if err := database.DB.Create(&cust).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cust)
}

func UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var cust models.Customer
	if err := database.DB.First(&cust, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pelanggan tidak ditemukan"})
		return
	}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&cust)
	c.JSON(http.StatusOK, cust)
}

func DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Customer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Pelanggan berhasil dihapus"})
}

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
	ProductID uint   `json:"product_id"`
	Type      string `json:"type"`   // "in" or "out"
	Quantity  int    `json:"quantity"`
	Reason    string `json:"reason"`
	Notes     string `json:"notes"`
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

	if input.Type == "in" {
		prod.Stock += input.Quantity
	} else {
		if prod.Stock < input.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Stok %s tidak mencukupi (sisa %d)", prod.Name, prod.Stock)})
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
