package handlers

import (
	"net/http"
	"strconv"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	repo *models.ProfileRepository
}

func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{
		repo: models.NewProfileRepository(db),
	}
}

// GetMyProfile retrieves the authenticated user's profile
func (h *ProfileHandler) GetMyProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	profile, err := h.repo.GetByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

// GetProfile retrieves a profile by user ID (admin)
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	profile, err := h.repo.GetByUserID(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Profile not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": profile,
	})
}

// CreateOrUpdateMyProfile creates or updates the authenticated user's profile
func (h *ProfileHandler) CreateOrUpdateMyProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	var input models.ProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if profile exists
	existingProfile, err := h.repo.GetByUserID(userID.(uint))

	if err == gorm.ErrRecordNotFound {
		// Create new profile
		profile := &models.Profile{
			UserID:      userID.(uint),
			PhoneNumber: input.PhoneNumber,
			Address:     input.Address,
			City:        input.City,
			Province:    input.Province,
			PostalCode:  input.PostalCode,
			Country:     input.Country,
			Avatar:      input.Avatar,
			Bio:         input.Bio,
			DateOfBirth: input.DateOfBirth,
			Gender:      input.Gender,
		}

		if err := h.repo.Create(profile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create profile",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Profile created successfully",
			"data":    profile,
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve profile",
		})
		return
	}

	// Update existing profile
	existingProfile.PhoneNumber = input.PhoneNumber
	existingProfile.Address = input.Address
	existingProfile.City = input.City
	existingProfile.Province = input.Province
	existingProfile.PostalCode = input.PostalCode
	existingProfile.Country = input.Country
	existingProfile.Avatar = input.Avatar
	existingProfile.Bio = input.Bio
	existingProfile.DateOfBirth = input.DateOfBirth
	existingProfile.Gender = input.Gender

	if err := h.repo.Update(existingProfile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    existingProfile,
	})
}

// DeleteProfile deletes a profile by ID (admin)
func (h *ProfileHandler) DeleteProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid profile ID",
		})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile deleted successfully",
	})
}
