package handlers

import (
	"net/http"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WishlistHandler struct {
	repo *models.WishlistRepository
}

func NewWishlistHandler(db *gorm.DB) *WishlistHandler {
	return &WishlistHandler{
		repo: models.NewWishlistRepository(db),
	}
}

// AddToWishlist adds a product to user's wishlist
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	var input models.WishlistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	// Check if already in wishlist
	exists, err := h.repo.Exists(userID.(uint), input.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check wishlist",
		})
		return
	}

	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Product already in wishlist",
		})
		return
	}

	wishlist := &models.Wishlist{
		UserID:    userID.(uint),
		ProductID: input.ProductID,
	}

	if err := h.repo.Create(wishlist); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to add to wishlist",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product added to wishlist",
		"data":    wishlist,
	})
}

// GetWishlist retrieves user's wishlist
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	wishlists, err := h.repo.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve wishlist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": wishlists,
	})
}

// RemoveFromWishlist removes a product from user's wishlist
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	var input models.WishlistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	if err := h.repo.Delete(userID.(uint), input.ProductID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove from wishlist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product removed from wishlist",
	})
}

// CheckWishlist checks if a product is in user's wishlist
func (h *WishlistHandler) CheckWishlist(c *gin.Context) {
	var input models.WishlistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	inWishlist, err := h.repo.Exists(userID.(uint), input.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check wishlist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"in_wishlist": inWishlist,
	})
}
