package handlers

import (
	"net/http"
	"pos-backend/internal/database"
	"pos-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// === REPORT & DASHBOARD HANDLERS ===

func GetDashboardStats(c *gin.Context) {
	var totalOrders int64
	var totalRevenue float64
	var totalItemsSold float64

	merchantID, _ := c.Get("merchant_id")
	outletID, _ := c.Get("outlet_id")

	orderQuery := database.DB.Model(&models.Order{}).Where("status = ?", "completed")
	if merchantID != nil {
		orderQuery = orderQuery.Where("merchant_id = ?", merchantID)
	}
	if outletID != nil {
		orderQuery = orderQuery.Where("outlet_id = ?", outletID)
	}

	getOrderItemQuery := func() *gorm.DB {
		q := database.DB.Model(&models.OrderItem{})
		if merchantID != nil || outletID != nil {
			q = q.Joins("INNER JOIN orders ON orders.id = order_items.order_id")
			if merchantID != nil {
				q = q.Where("orders.merchant_id = ?", merchantID)
			}
			if outletID != nil {
				q = q.Where("orders.outlet_id = ?", outletID)
			}
		}
		return q
	}

	orderQuery.Count(&totalOrders)

	var paymentStats struct {
		TotalRevenue     float64
		TotalCash        float64
		TotalQris        float64
		TotalDebit       float64
		TotalCredit      float64
		TotalPlatformFee float64
	}
	
	orderQuery.Select(`
		COALESCE(SUM(grand_total), 0) as total_revenue,
		COALESCE(SUM(CASE WHEN payment_method = 'cash' THEN grand_total ELSE 0 END), 0) as total_cash,
		COALESCE(SUM(CASE WHEN payment_method = 'qris' THEN grand_total ELSE 0 END), 0) as total_qris,
		COALESCE(SUM(CASE WHEN payment_method = 'debit' THEN grand_total ELSE 0 END), 0) as total_debit,
		COALESCE(SUM(CASE WHEN payment_method = 'credit' THEN grand_total ELSE 0 END), 0) as total_credit,
		COALESCE(SUM(platform_fee), 0) as total_platform_fee
	`).Scan(&paymentStats)

	totalRevenue = paymentStats.TotalRevenue
	getOrderItemQuery().Select("COALESCE(SUM(order_items.quantity), 0)").Scan(&totalItemsSold)

	type ProductSalesStat struct {
		ProductID   uuid.UUID `json:"product_id"`
		ProductName string    `json:"product_name"`
		Artist      string    `json:"artist"`
		ProductType string    `json:"product_type"`
		TotalQty    float64   `json:"total_qty"`
		TotalSales  float64   `json:"total_sales"`
		Stock       float64   `json:"stock"`
		Price       float64   `json:"price"`
		Unit        string    `json:"unit"`
	}

	// Top Selling Products (Barang Paling Laku)
	topProducts := []ProductSalesStat{}
	topProdQ := database.DB.Table("products p").
		Select("p.id as product_id, p.name as product_name, COALESCE(NULLIF(p.artist, ''), 'Umum') as artist, COALESCE(NULLIF(p.product_type, ''), 'Umum') as product_type, COALESCE(SUM(oi.quantity), 0) as total_qty, COALESCE(SUM(oi.subtotal), 0) as total_sales, p.stock, p.price, p.unit").
		Joins("INNER JOIN order_items oi ON oi.product_id = p.id").
		Joins("INNER JOIN orders o ON o.id = oi.order_id").
		Where("o.status = 'completed'")
	if merchantID != nil { topProdQ = topProdQ.Where("p.merchant_id = ?", merchantID) }
	if outletID != nil { topProdQ = topProdQ.Where("p.outlet_id = ?", outletID) }
	topProdQ.Group("p.id, p.name, p.artist, p.product_type, p.stock, p.price, p.unit").
		Order("total_qty desc").
		Limit(10).
		Scan(&topProducts)

	// Slow Moving / Least Sold Products (Barang Kurang Laku)
	leastProducts := []ProductSalesStat{}
	leastProdQ := database.DB.Table("products p").
		Select("p.id as product_id, p.name as product_name, COALESCE(NULLIF(p.artist, ''), 'Umum') as artist, COALESCE(NULLIF(p.product_type, ''), 'Umum') as product_type, COALESCE(SUM(oi.quantity), 0) as total_qty, COALESCE(SUM(oi.subtotal), 0) as total_sales, p.stock, p.price, p.unit").
		Joins("LEFT JOIN order_items oi ON oi.product_id = p.id").
		Where("p.is_active = ?", true)
	if merchantID != nil { leastProdQ = leastProdQ.Where("p.merchant_id = ?", merchantID) }
	if outletID != nil { leastProdQ = leastProdQ.Where("p.outlet_id = ?", outletID) }
	leastProdQ.Group("p.id, p.name, p.artist, p.product_type, p.stock, p.price, p.unit").
		Order("total_qty asc, p.stock desc").
		Limit(10).
		Scan(&leastProducts)

	// All Sold Products List (List Seluruh Barang Laku)
	allSoldProducts := []ProductSalesStat{}
	getOrderItemQuery().
		Select("order_items.product_id, order_items.product_name, COALESCE(NULLIF(order_items.artist, ''), 'Umum') as artist, COALESCE(NULLIF(order_items.product_type, ''), 'Umum') as product_type, SUM(order_items.quantity) as total_qty, SUM(order_items.subtotal) as total_sales, MAX(order_items.product_price) as price, MAX(order_items.unit) as unit").
		Where("orders.status = 'completed'").
		Group("order_items.product_id, order_items.product_name, order_items.artist, order_items.product_type, order_items.unit").
		Order("total_qty desc").
		Scan(&allSoldProducts)

	// Sub-category Sales Summaries
	type SubGroupStat struct {
		Name       string  `json:"name"`
		TotalQty   float64   `json:"total_qty"`
		TotalSales float64   `json:"total_sales"`
	}

	salesByArtist := []SubGroupStat{}
	getOrderItemQuery().
		Select("COALESCE(NULLIF(order_items.artist, ''), 'Lainnya / Umum') as name, SUM(order_items.quantity) as total_qty, SUM(order_items.subtotal) as total_sales").
		Where("orders.status = 'completed'").
		Group("name").
		Order("total_sales desc").
		Scan(&salesByArtist)

	salesByType := []SubGroupStat{}
	getOrderItemQuery().
		Select("COALESCE(NULLIF(order_items.product_type, ''), 'Lainnya / Umum') as name, SUM(order_items.quantity) as total_qty, SUM(order_items.subtotal) as total_sales").
		Where("orders.status = 'completed'").
		Group("name").
		Order("total_sales desc").
		Scan(&salesByType)

	// Recent Orders
	var recentOrders []models.Order
	recentOrderQ := database.DB.Order("created_at desc").Limit(5)
	if merchantID != nil { recentOrderQ = recentOrderQ.Where("merchant_id = ?", merchantID) }
	if outletID != nil { recentOrderQ = recentOrderQ.Where("outlet_id = ?", outletID) }
	recentOrderQ.Find(&recentOrders)

	// Store Setting for Report Header
	var storeSetting models.StoreSetting
	settingQ := database.DB.Model(&models.StoreSetting{})
	if merchantID != nil { settingQ = settingQ.Where("merchant_id = ?", merchantID) }
	settingQ.First(&storeSetting)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":       totalOrders,
		"total_revenue":      totalRevenue,
		"total_cash":         paymentStats.TotalCash,
		"total_qris":         paymentStats.TotalQris,
		"total_debit":        paymentStats.TotalDebit,
		"total_credit":       paymentStats.TotalCredit,
		"total_platform_fee": paymentStats.TotalPlatformFee,
		"total_items_sold":   totalItemsSold,
		"top_products":       topProducts,
		"least_products":     leastProducts,
		"all_sold_products":  allSoldProducts,
		"sales_by_artist":    salesByArtist,
		"sales_by_type":      salesByType,
		"recent_orders":      recentOrders,
		"store_setting":      storeSetting,
	})
}
