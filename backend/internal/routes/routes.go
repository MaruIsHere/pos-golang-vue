package routes

import (
	"net/http"
	"pos-backend/internal/config"
	"pos-backend/internal/handlers"
	"pos-backend/internal/websocket"
	"pos-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/settings", handlers.GetSettings)
	{
		// WebSocket & Webhooks
		api.GET("/ws", websocket.HandleWebSocket)
		api.POST("/webhooks/qris", handlers.HandleQrisWebhook)

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
			auth.POST("/login", handlers.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/profile", middleware.RoleMiddleware("owner", "admin"), handlers.GetProfile)
			protected.PUT("/profile", middleware.RoleMiddleware("owner", "admin"), handlers.UpdateProfile)

			// Reusable product catalogs
			protected.GET("/artists", handlers.GetArtists)
			protected.POST("/artists", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.CreateArtist)
			protected.PUT("/artists/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.UpdateArtist)
			protected.DELETE("/artists/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.DeleteArtist)
			protected.GET("/product-types", handlers.GetProductTypes)
			protected.POST("/product-types", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.CreateProductType)
			protected.PUT("/product-types/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.UpdateProductType)
			protected.DELETE("/product-types/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.DeleteProductType)

			// Categories
			protected.GET("/categories", handlers.GetCategories)
			protected.POST("/categories", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.CreateCategory)
			protected.PUT("/categories/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.UpdateCategory)
			protected.DELETE("/categories/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.DeleteCategory)

			// Products
			protected.GET("/products", handlers.GetProducts)
			protected.GET("/products/filters", handlers.GetProductFilters)
			protected.POST("/products", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.CreateProduct)
			protected.PUT("/products/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.UpdateProduct)
			protected.DELETE("/products/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.DeleteProduct)
			protected.POST("/products/images", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.UploadProductImage)

			// Orders / Transactions
			protected.POST("/orders", handlers.CreateOrder)
			protected.POST("/orders/payment-proof", handlers.UploadPaymentProof)
			protected.GET("/orders", handlers.GetOrders)
			protected.GET("/orders/:id", handlers.GetOrderById)
			protected.POST("/orders/:id/refund", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.RefundOrder)

			// Customers
			protected.GET("/customers", handlers.GetCustomers)
			protected.POST("/customers", handlers.CreateCustomer)
			protected.PUT("/customers/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.UpdateCustomer)
			protected.DELETE("/customers/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.DeleteCustomer)

			// Stock Movements (Inventory - Receive & Issue)
			protected.GET("/stock-movements", handlers.GetStockMovements)
			protected.POST("/stock-movements", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.CreateStockMovement)

			// Reports & Dashboard
			protected.GET("/reports/dashboard", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.GetDashboardStats)

			// Store & DB Settings
			protected.PUT("/settings", middleware.RoleMiddleware("admin", "owner"), handlers.UpdateSettings)
			protected.POST("/settings/switch-db", middleware.RoleMiddleware("admin", "owner"), handlers.SwitchDatabase)

			// Store Management (Registrasi & Kelola Toko)
			protected.GET("/stores", handlers.GetStores)
			protected.POST("/stores", middleware.RoleMiddleware("administrator", "owner", "admin"), handlers.CreateStore)
			protected.PUT("/stores/:id", middleware.RoleMiddleware("administrator", "owner", "admin"), handlers.UpdateStore)
			protected.DELETE("/stores/:id", middleware.RoleMiddleware("administrator", "owner", "admin"), handlers.DeleteStore)
			protected.POST("/stores/assign-user", middleware.RoleMiddleware("administrator", "owner", "admin"), handlers.AssignUserStore)
			protected.POST("/stores/import-product", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.ImportMasterProductToStore)
			protected.POST("/stores/import-batch", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.ImportBatchMasterProductsToStore)

			// Vouchers
			protected.GET("/vouchers", handlers.GetVouchers)
			protected.POST("/vouchers", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.CreateVoucher)
			protected.DELETE("/vouchers/:id", middleware.RoleMiddleware("administrator", "owner", "admin", "kepala_kasir"), handlers.DeleteVoucher)

			// User Management (Khusus Owner & Admin)
			protected.GET("/users", middleware.RoleMiddleware("admin", "owner"), handlers.GetUsers)
			protected.GET("/users/:id", middleware.RoleMiddleware("admin", "owner"), handlers.GetUserProfile)
			protected.PUT("/users/:id", middleware.RoleMiddleware("admin", "owner"), handlers.UpdateUserProfile)
			protected.POST("/users", middleware.RoleMiddleware("administrator", "owner"), handlers.Register)
			protected.PUT("/users/:id/password", middleware.RoleMiddleware("admin", "owner"), handlers.ChangeUserPassword)
			protected.DELETE("/users/:id", middleware.RoleMiddleware("admin", "owner"), handlers.DeleteUser)
		}
	}
}
