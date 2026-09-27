package handlers

import (
	"net/http"
	"pos-backend/database"
	"pos-backend/models"

	"github.com/gin-gonic/gin"
)

// === REPORT & DASHBOARD HANDLERS ===

func GetDashboardStats(c *gin.Context) {
	var totalOrders int64
	var totalRevenue float64
	var totalItemsSold int64

	database.DB.Model(&models.Order{}).Where("status = ?", "completed").Count(&totalOrders)
	database.DB.Model(&models.Order{}).Where("status = ?", "completed").Select("COALESCE(SUM(grand_total), 0)").Scan(&totalRevenue)
	database.DB.Model(&models.OrderItem{}).Select("COALESCE(SUM(quantity), 0)").Scan(&totalItemsSold)

	type ProductSalesStat struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		Artist      string  `json:"artist"`
		ProductType string  `json:"product_type"`
		TotalQty    int     `json:"total_qty"`
		TotalSales  float64 `json:"total_sales"`
		Stock       int     `json:"stock"`
		Price       float64 `json:"price"`
	}

	// Top Selling Products (Barang Paling Laku)
	var topProducts []ProductSalesStat
	database.DB.Table("products p").
		Select("p.id as product_id, p.name as product_name, COALESCE(NULLIF(p.artist, ''), 'Umum') as artist, COALESCE(NULLIF(p.product_type, ''), 'Umum') as product_type, COALESCE(SUM(oi.quantity), 0) as total_qty, COALESCE(SUM(oi.subtotal), 0) as total_sales, p.stock, p.price").
		Joins("INNER JOIN order_items oi ON oi.product_id = p.id").
		Group("p.id, p.name, p.artist, p.product_type, p.stock, p.price").
		Order("total_qty desc").
		Limit(10).
		Scan(&topProducts)

	// Slow Moving / Least Sold Products (Barang Kurang Laku)
	var leastProducts []ProductSalesStat
	database.DB.Table("products p").
		Select("p.id as product_id, p.name as product_name, COALESCE(NULLIF(p.artist, ''), 'Umum') as artist, COALESCE(NULLIF(p.product_type, ''), 'Umum') as product_type, COALESCE(SUM(oi.quantity), 0) as total_qty, COALESCE(SUM(oi.subtotal), 0) as total_sales, p.stock, p.price").
		Joins("LEFT JOIN order_items oi ON oi.product_id = p.id").
		Where("p.is_active = ?", true).
		Group("p.id, p.name, p.artist, p.product_type, p.stock, p.price").
		Order("total_qty asc, p.stock desc").
		Limit(10).
		Scan(&leastProducts)

	// All Sold Products List (List Seluruh Barang Laku)
	var allSoldProducts []ProductSalesStat
	database.DB.Table("order_items").
		Select("product_id, product_name, COALESCE(NULLIF(artist, ''), 'Umum') as artist, COALESCE(NULLIF(product_type, ''), 'Umum') as product_type, SUM(quantity) as total_qty, SUM(subtotal) as total_sales, MAX(product_price) as price").
		Group("product_id, product_name, artist, product_type").
		Order("total_qty desc").
		Scan(&allSoldProducts)

	// Sub-category Sales Summaries
	type SubGroupStat struct {
		Name       string  `json:"name"`
		TotalQty   int     `json:"total_qty"`
		TotalSales float64 `json:"total_sales"`
	}

	var salesByArtist []SubGroupStat
	database.DB.Table("order_items").
		Select("COALESCE(NULLIF(artist, ''), 'Lainnya / Umum') as name, SUM(quantity) as total_qty, SUM(subtotal) as total_sales").
		Group("name").
		Order("total_sales desc").
		Scan(&salesByArtist)

	var salesByType []SubGroupStat
	database.DB.Table("order_items").
		Select("COALESCE(NULLIF(product_type, ''), 'Lainnya / Umum') as name, SUM(quantity) as total_qty, SUM(subtotal) as total_sales").
		Group("name").
		Order("total_sales desc").
		Scan(&salesByType)

	// Recent Orders
	var recentOrders []models.Order
	database.DB.Order("created_at desc").Limit(5).Find(&recentOrders)

	// Store Setting for Report Header
	var storeSetting models.StoreSetting
	database.DB.First(&storeSetting)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":      totalOrders,
		"total_revenue":     totalRevenue,
		"total_items_sold":  totalItemsSold,
		"top_products":      topProducts,
		"least_products":    leastProducts,
		"all_sold_products": allSoldProducts,
		"sales_by_artist":   salesByArtist,
		"sales_by_type":     salesByType,
		"recent_orders":     recentOrders,
		"store_setting":     storeSetting,
	})
}
