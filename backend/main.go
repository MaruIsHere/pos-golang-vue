package main

import (
	"log"
	"net/http"

	"pos-backend/config"
	"pos-backend/database"
	"pos-backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.Printf("Starting POS Golang Backend on port %s...", cfg.Port)
	log.Printf("Database Engine: %s", cfg.DbEngine)

	_, err := database.InitDB(cfg.DbEngine, cfg.MysqlDsn, cfg.SqlitePath)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	router := gin.Default()

	// Enable CORS for frontend integration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// API Route Group
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

		// Reports & Dashboard
		api.GET("/reports/dashboard", handlers.GetDashboardStats)

		// Store & DB Settings
		api.GET("/settings", handlers.GetSettings)
		api.PUT("/settings", handlers.UpdateSettings)
		api.POST("/settings/switch-db", handlers.SwitchDatabase)
	}

	// Serve frontend static files if dist folder exists
	router.Static("/assets", "../frontend/dist/assets")
	router.NoRoute(func(c *gin.Context) {
		c.File("../frontend/dist/index.html")
	})

	log.Printf("Server running at http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
