package main

import (
	"log"

	"github.com/dinulimann/goback/config"
	"github.com/dinulimann/goback/database"
	"github.com/dinulimann/goback/models"
	"github.com/dinulimann/goback/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize config
	cfg := config.NewConfig()

	// Refuse to start without a signing key: a built-in default would let
	// anyone who read this code forge login tokens.
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set (see .env.example)")
	}

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Get database instance
	db := database.GetDB()

	// Auto migrate database schema
	if err := db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Category{},
		&models.Product{},
		&models.Coupon{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
		&models.Wishlist{},
		&models.Payment{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Setup routes with database and config
	routes.SetupRoutes(router, db, cfg)

	// Start server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
