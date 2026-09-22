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
