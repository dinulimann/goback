package handlers

import (
	"net/http"
	"strconv"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewHandler struct {
	repo *models.ReviewRepository
}

func NewReviewHandler(db *gorm.DB) *ReviewHandler {
	return &ReviewHandler{
		repo: models.NewReviewRepository(db),
	}
}

// CreateReview creates a new review
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	review := &models.Review{
		ProductID: input.ProductID,
		UserID:    userID.(uint),
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	if err := h.repo.Create(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create review",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review created successfully",
		"data":    review,
	})
}

// GetProductReviews retrieves reviews for a product
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	reviews, err := h.repo.GetByProductID(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve reviews",
		})
		return
	}

	// Get average rating and count
	avgRating, _ := h.repo.GetAverageRating(uint(productID))
	count, _ := h.repo.CountByProductID(uint(productID))

	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
		"stats": gin.H{
			"average_rating": avgRating,
			"total_reviews":  count,
		},
	})
}

// GetUserReviews retrieves reviews by a user
func (h *ReviewHandler) GetUserReviews(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	reviews, err := h.repo.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve reviews",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
	})
}

// UpdateReview updates a review
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid review ID",
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

	review, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Review not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve review",
		})
		return
	}

	// Check if user owns this review
	if review.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only update your own reviews",
		})
		return
	}

	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	review.Rating = input.Rating
	review.Comment = input.Comment

	if err := h.repo.Update(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update review",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review updated successfully",
		"data":    review,
	})
}

// DeleteReview deletes a review
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid review ID",
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

	review, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Review not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve review",
		})
		return
	}

	// Check if user owns this review
	if review.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You can only delete your own reviews",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete review",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Review deleted successfully",
	})
}
