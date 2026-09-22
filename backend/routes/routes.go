package routes

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Status endpoint
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":    "online",
				"system":    "POS Kasir System Golang Backend",
				"db_engine": config.AppConfig.DbEngine,
			})
		})

		// Categories
		api.GET("/categories", handlers.GetCategories)
		api.POST("/categories", handlers.CreateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)

		// Products
		api.GET("/products", handlers.GetProducts)
		api.POST("/products", handlers.CreateProduct)
		api.PUT("/products/:id", handlers.UpdateProduct)
		api.DELETE("/products/:id", handlers.DeleteProduct)

		// Orders / Transactions
		api.POST("/orders", handlers.CreateOrder)
		api.GET("/orders", handlers.GetOrders)
		api.GET("/orders/:id", handlers.GetOrderById)
		api.POST("/orders/:id/refund", handlers.RefundOrder)

		// Customers
		api.GET("/customers", handlers.GetCustomers)
		api.POST("/customers", handlers.CreateCustomer)
		api.PUT("/customers/:id", handlers.UpdateCustomer)
		api.DELETE("/customers/:id", handlers.DeleteCustomer)

		// Stock Movements (Inventory - Receive & Issue)
		api.GET("/stock-movements", handlers.GetStockMovements)
		api.POST("/stock-movements", handlers.CreateStockMovement)

		// Reports & Dashboard
		api.GET("/reports/dashboard", handlers.GetDashboardStats)

		// Store & DB Settings
		api.GET("/settings", handlers.GetSettings)
		api.PUT("/settings", handlers.UpdateSettings)
		api.POST("/settings/switch-db", handlers.SwitchDatabase)

		// Vouchers
		api.GET("/vouchers", handlers.GetVouchers)
		api.POST("/vouchers", handlers.CreateVoucher)
		api.DELETE("/vouchers/:id", handlers.DeleteVoucher)
	}
}
