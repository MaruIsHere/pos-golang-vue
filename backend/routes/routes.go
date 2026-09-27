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
	api.GET("/settings", handlers.GetSettings)
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
			protected.POST("/categories", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.CreateCategory)
			protected.DELETE("/categories/:id", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.DeleteCategory)

			// Products
			protected.GET("/products", handlers.GetProducts)
			protected.GET("/products/filters", handlers.GetProductFilters)
			protected.POST("/products", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.CreateProduct)
			protected.PUT("/products/:id", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.UpdateProduct)
			protected.DELETE("/products/:id", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.DeleteProduct)

			// Orders / Transactions
			protected.POST("/orders", handlers.CreateOrder)
			protected.GET("/orders", handlers.GetOrders)
			protected.GET("/orders/:id", handlers.GetOrderById)
			protected.POST("/orders/:id/refund", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.RefundOrder)

			// Customers
			protected.GET("/customers", handlers.GetCustomers)
			protected.POST("/customers", handlers.CreateCustomer)
			protected.PUT("/customers/:id", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.UpdateCustomer)
			protected.DELETE("/customers/:id", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.DeleteCustomer)

			// Stock Movements (Inventory - Receive & Issue)
			protected.GET("/stock-movements", handlers.GetStockMovements)
			protected.POST("/stock-movements", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.CreateStockMovement)

			// Reports & Dashboard
			protected.GET("/reports/dashboard", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.GetDashboardStats)

			// Store & DB Settings
			protected.GET("/settings", middleware.RoleMiddleware("admin", "owner"), handlers.GetSettings)
			protected.PUT("/settings", middleware.RoleMiddleware("admin", "owner"), handlers.UpdateSettings)
			protected.POST("/settings/switch-db", middleware.RoleMiddleware("admin", "owner"), handlers.SwitchDatabase)

			// Vouchers
			protected.GET("/vouchers", handlers.GetVouchers)
			protected.POST("/vouchers", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.CreateVoucher)
			protected.DELETE("/vouchers/:id", middleware.RoleMiddleware("admin", "owner", "kepala_kasir"), handlers.DeleteVoucher)
		}
	}
}
