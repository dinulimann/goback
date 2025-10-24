package handlers

import (
	"net/http"
	"strconv"

	"github.com/dinulimann/goback/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CouponHandler handles coupon-related requests
type CouponHandler struct {
	couponRepo *models.CouponRepository
}

// NewCouponHandler creates a new CouponHandler
func NewCouponHandler(db *gorm.DB) *CouponHandler {
	return &CouponHandler{
		couponRepo: models.NewCouponRepository(db),
	}
}

// CreateCoupon creates a new coupon
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var coupon models.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate discount type
	if coupon.DiscountType != "percentage" && coupon.DiscountType != "fixed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "discount_type must be 'percentage' or 'fixed'"})
		return
	}

	// Validate discount value
	if coupon.DiscountType == "percentage" && (coupon.DiscountValue <= 0 || coupon.DiscountValue > 100) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "percentage discount must be between 0 and 100"})
		return
	}

	if coupon.DiscountType == "fixed" && coupon.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fixed discount must be greater than 0"})
		return
	}

	// Validate dates
	if coupon.EndDate.Before(coupon.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	if err := h.couponRepo.Create(&coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Coupon created successfully",
		"coupon":  coupon,
	})
}

// GetAllCoupons retrieves all coupons
func (h *CouponHandler) GetAllCoupons(c *gin.Context) {
	coupons, err := h.couponRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coupons": coupons,
		"count":   len(coupons),
	})
}

// GetActiveCoupons retrieves all active coupons
func (h *CouponHandler) GetActiveCoupons(c *gin.Context) {
	coupons, err := h.couponRepo.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve active coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"coupons": coupons,
		"count":   len(coupons),
	})
}

// GetCouponByID retrieves a coupon by ID
func (h *CouponHandler) GetCouponByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	coupon, err := h.couponRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"coupon": coupon})
}

// UpdateCoupon updates a coupon
func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	coupon, err := h.couponRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	var updateData models.Coupon
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate discount type if provided
	if updateData.DiscountType != "" && updateData.DiscountType != "percentage" && updateData.DiscountType != "fixed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "discount_type must be 'percentage' or 'fixed'"})
		return
	}

	// Update fields
	if updateData.Code != "" {
		coupon.Code = updateData.Code
	}
	if updateData.Description != "" {
		coupon.Description = updateData.Description
	}
	if updateData.DiscountType != "" {
		coupon.DiscountType = updateData.DiscountType
	}
	if updateData.DiscountValue > 0 {
		coupon.DiscountValue = updateData.DiscountValue
	}
	coupon.MinPurchase = updateData.MinPurchase
	coupon.MaxDiscount = updateData.MaxDiscount
	coupon.UsageLimit = updateData.UsageLimit
	coupon.IsActive = updateData.IsActive
	if !updateData.StartDate.IsZero() {
		coupon.StartDate = updateData.StartDate
	}
	if !updateData.EndDate.IsZero() {
		coupon.EndDate = updateData.EndDate
	}

	// Validate dates
	if coupon.EndDate.Before(coupon.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	if err := h.couponRepo.Update(coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon updated successfully",
		"coupon":  coupon,
	})
}

// DeleteCoupon deletes a coupon
func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	if err := h.couponRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupon deleted successfully"})
}

// ValidateCouponRequest represents the request body for validating a coupon
type ValidateCouponRequest struct {
	Code        string  `json:"code" binding:"required"`
	OrderAmount float64 `json:"order_amount" binding:"required,gt=0"`
}

// ValidateCoupon validates a coupon and returns discount information
func (h *CouponHandler) ValidateCoupon(c *gin.Context) {
	var req ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon, err := h.couponRepo.ValidateCoupon(req.Code, req.OrderAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discount := coupon.CalculateDiscount(req.OrderAmount)
	finalAmount := req.OrderAmount - discount

	c.JSON(http.StatusOK, gin.H{
		"valid":         true,
		"coupon":        coupon,
		"order_amount":  req.OrderAmount,
		"discount":      discount,
		"final_amount":  finalAmount,
		"message":       "Coupon applied successfully",
	})
}

// ToggleActiveCoupon toggles the active status of a coupon
func (h *CouponHandler) ToggleActiveCoupon(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	coupon, err := h.couponRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	coupon.IsActive = !coupon.IsActive

	if err := h.couponRepo.Update(coupon); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle coupon status"})
		return
	}

	status := "inactive"
	if coupon.IsActive {
		status = "active"
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon status updated successfully",
		"coupon":  coupon,
		"status":  status,
	})
}
