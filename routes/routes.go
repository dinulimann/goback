package routes

import (
	"github.com/dinulimann/goback/handlers"
	"github.com/dinulimann/goback/middleware"
	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Apply global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// Initialize repositories
	userRepo := models.NewUserRepository(db)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(userRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	// Health check endpoint
	router.GET("/health", healthHandler.Check)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public auth routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.JWTAuth())
		{
			// Auth profile
			protected.GET("/auth/profile", authHandler.GetProfile)

			// User routes (now protected)
			users := protected.Group("/users")
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.GetUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}
}
