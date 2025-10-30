package handlers

import (
	"net/http"
	"strconv"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddressHandler handles address-related requests
type AddressHandler struct {
	addressRepo *models.AddressRepository
}

// NewAddressHandler creates a new AddressHandler
func NewAddressHandler(db *gorm.DB) *AddressHandler {
	return &AddressHandler{
		addressRepo: models.NewAddressRepository(db),
	}
}

// CreateAddress creates a new address for the authenticated user
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var input models.AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := &models.Address{
		UserID:        userID.(uint),
		Label:         input.Label,
		RecipientName: input.RecipientName,
		PhoneNumber:   input.PhoneNumber,
		Address:       input.Address,
		City:          input.City,
		Province:      input.Province,
		PostalCode:    input.PostalCode,
		IsDefault:     input.IsDefault,
		Notes:         input.Notes,
	}

	if err := h.addressRepo.Create(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Address created successfully",
		"address": address,
	})
}

// GetMyAddresses retrieves all addresses for the authenticated user
func (h *AddressHandler) GetMyAddresses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	addresses, err := h.addressRepo.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"addresses": addresses,
		"count":     len(addresses),
	})
}

// GetDefaultAddress retrieves the default address for the authenticated user
func (h *AddressHandler) GetDefaultAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	address, err := h.addressRepo.GetDefaultByUserID(userID.(uint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No default address found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve default address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

// GetAddress retrieves a specific address by ID
func (h *AddressHandler) GetAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	// Check ownership
	owns, err := h.addressRepo.CheckOwnership(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	address, err := h.addressRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"address": address})
}

// UpdateAddress updates an existing address
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	// Check ownership
	owns, err := h.addressRepo.CheckOwnership(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	address, err := h.addressRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	var input models.AddressInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	address.Label = input.Label
	address.RecipientName = input.RecipientName
	address.PhoneNumber = input.PhoneNumber
	address.Address = input.Address
	address.City = input.City
	address.Province = input.Province
	address.PostalCode = input.PostalCode
	address.IsDefault = input.IsDefault
	address.Notes = input.Notes

	if err := h.addressRepo.Update(address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Address updated successfully",
		"address": address,
	})
}

// DeleteAddress deletes an address
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	// Check ownership
	owns, err := h.addressRepo.CheckOwnership(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.addressRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

// SetDefaultAddress sets an address as the default
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
		return
	}

	// Check ownership
	owns, err := h.addressRepo.CheckOwnership(uint(id), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify ownership"})
		return
	}
	if !owns {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.addressRepo.SetDefault(uint(id), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default address updated successfully"})
}
