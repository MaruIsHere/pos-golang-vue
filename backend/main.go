package main

import (
	"net/http"
	"log"
	"pos-backend/config"
	"pos-backend/database"
	"pos-backend/routes"
	"strings"

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
		AllowCredentials: false,
	}))

	// API Route Group
	routes.RegisterAPIRoutes(router)

	// Serve frontend static files if dist folder exists
	router.Static("/assets", "../frontend/dist/assets")
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API tidak ditemukan. Pastikan backend memakai versi terbaru."})
			return
		}
		c.File("../frontend/dist/index.html")
	})

	log.Printf("Server running at http://localhost:%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
