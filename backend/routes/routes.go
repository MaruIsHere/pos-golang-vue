package routes

import (
	"net/http"
	"pos-backend/config"
	"pos-backend/handlers"
	"pos-backend/middleware"

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

		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Categories
			protected.GET("/categories", handlers.GetCategories)
			protected.POST("/categories", handlers.CreateCategory)
			protected.DELETE("/categories/:id", handlers.DeleteCategory)

			// Products
			protected.GET("/products", handlers.GetProducts)
			protected.GET("/products/filters", handlers.GetProductFilters)
			protected.POST("/products", handlers.CreateProduct)
			protected.PUT("/products/:id", handlers.UpdateProduct)
			protected.DELETE("/products/:id", handlers.DeleteProduct)

			// Orders / Transactions
			protected.POST("/orders", handlers.CreateOrder)
			protected.GET("/orders", handlers.GetOrders)
			protected.GET("/orders/:id", handlers.GetOrderById)
			protected.POST("/orders/:id/refund", middleware.RoleMiddleware("admin"), handlers.RefundOrder)

			// Customers
			protected.GET("/customers", handlers.GetCustomers)
			protected.POST("/customers", handlers.CreateCustomer)
			protected.PUT("/customers/:id", handlers.UpdateCustomer)
			protected.DELETE("/customers/:id", handlers.DeleteCustomer)

			// Stock Movements (Inventory - Receive & Issue)
			protected.GET("/stock-movements", handlers.GetStockMovements)
			protected.POST("/stock-movements", handlers.CreateStockMovement)

			// Reports & Dashboard
			protected.GET("/reports/dashboard", handlers.GetDashboardStats)

			// Store & DB Settings
			protected.GET("/settings", handlers.GetSettings)
			protected.PUT("/settings", middleware.RoleMiddleware("admin"), handlers.UpdateSettings)
			protected.POST("/settings/switch-db", middleware.RoleMiddleware("admin"), handlers.SwitchDatabase)

			// Vouchers
			protected.GET("/vouchers", handlers.GetVouchers)
			protected.POST("/vouchers", middleware.RoleMiddleware("admin"), handlers.CreateVoucher)
			protected.DELETE("/vouchers/:id", middleware.RoleMiddleware("admin"), handlers.DeleteVoucher)
		}
	}
}
