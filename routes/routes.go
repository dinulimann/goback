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
	categoryHandler := handlers.NewCategoryHandler(db)
	productHandler := handlers.NewProductHandler(db)
	profileHandler := handlers.NewProfileHandler(db)
	orderHandler := handlers.NewOrderHandler(db)
	reviewHandler := handlers.NewReviewHandler(db)
	wishlistHandler := handlers.NewWishlistHandler(db)

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

		// Public category routes
		categories := v1.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)
		}

		// Public product routes
		products := v1.Group("/products")
		{
			products.GET("", productHandler.GetAllProducts) // Now supports search & filters
			products.GET("/:id", productHandler.GetProduct)
			products.GET("/sku/:sku", productHandler.GetProductBySKU)
			products.GET("/category/:categoryId", productHandler.GetProductsByCategory)
		}

		// Public review routes
		reviews := v1.Group("/reviews")
		{
			reviews.GET("/product/:productId", reviewHandler.GetProductReviews)
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

			// Category management routes (protected)
			categoryManagement := protected.Group("/categories")
			{
				categoryManagement.POST("", categoryHandler.CreateCategory)
				categoryManagement.PUT("/:id", categoryHandler.UpdateCategory)
				categoryManagement.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// Product management routes (protected)
			productManagement := protected.Group("/products")
			{
				productManagement.POST("", productHandler.CreateProduct)
				productManagement.PUT("/:id", productHandler.UpdateProduct)
				productManagement.DELETE("/:id", productHandler.DeleteProduct)
			}

			// Profile routes (protected)
			profile := protected.Group("/profile")
			{
				profile.GET("", profileHandler.GetMyProfile)
				profile.POST("", profileHandler.CreateOrUpdateMyProfile)
				profile.PUT("", profileHandler.CreateOrUpdateMyProfile)
				profile.GET("/user/:userId", profileHandler.GetProfile) // Admin only
				profile.DELETE("/:id", profileHandler.DeleteProfile)    // Admin only
			}

			// Order routes (protected)
			orders := protected.Group("/orders")
			{
				orders.POST("", orderHandler.CreateOrder)
				orders.GET("", orderHandler.GetMyOrders)
				orders.GET("/all", orderHandler.GetAllOrders) // Admin only
				orders.GET("/:id", orderHandler.GetOrder)
				orders.PUT("/:id/status", orderHandler.UpdateOrderStatus) // Admin only
				orders.POST("/:id/cancel", orderHandler.CancelOrder)
			}

			// Review routes (protected)
			reviewManagement := protected.Group("/reviews")
			{
				reviewManagement.POST("", reviewHandler.CreateReview)
				reviewManagement.GET("/my", reviewHandler.GetUserReviews)
				reviewManagement.PUT("/:id", reviewHandler.UpdateReview)
				reviewManagement.DELETE("/:id", reviewHandler.DeleteReview)
			}

			// Wishlist routes (protected)
			wishlist := protected.Group("/wishlist")
			{
				wishlist.POST("", wishlistHandler.AddToWishlist)
				wishlist.GET("", wishlistHandler.GetWishlist)
				wishlist.DELETE("", wishlistHandler.RemoveFromWishlist)
				wishlist.POST("/check", wishlistHandler.CheckWishlist)
			}
		}
	}
}
