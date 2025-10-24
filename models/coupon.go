package models

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Coupon represents a discount coupon
type Coupon struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Code          string    `gorm:"uniqueIndex;not null" json:"code"`
	Description   string    `gorm:"type:text" json:"description"`
	DiscountType  string    `gorm:"type:varchar(20);not null" json:"discount_type"` // "percentage" or "fixed"
	DiscountValue float64   `gorm:"not null" json:"discount_value"`
	MinPurchase   float64   `gorm:"default:0" json:"min_purchase"`
	MaxDiscount   float64   `gorm:"default:0" json:"max_discount"` // Max discount for percentage type
	UsageLimit    int       `gorm:"default:0" json:"usage_limit"`  // 0 means unlimited
	UsedCount     int       `gorm:"default:0" json:"used_count"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CouponRepository handles database operations for coupons
type CouponRepository struct {
	db *gorm.DB
}

// NewCouponRepository creates a new CouponRepository
func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{db: db}
}

// Create creates a new coupon
func (r *CouponRepository) Create(coupon *Coupon) error {
	// Convert code to uppercase for consistency
	coupon.Code = strings.ToUpper(coupon.Code)
	return r.db.Create(coupon).Error
}

// GetByID retrieves a coupon by ID
func (r *CouponRepository) GetByID(id uint) (*Coupon, error) {
	var coupon Coupon
	err := r.db.First(&coupon, id).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetByCode retrieves a coupon by code
func (r *CouponRepository) GetByCode(code string) (*Coupon, error) {
	var coupon Coupon
	code = strings.ToUpper(code)
	err := r.db.Where("code = ?", code).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetAll retrieves all coupons
func (r *CouponRepository) GetAll() ([]Coupon, error) {
	var coupons []Coupon
	err := r.db.Order("created_at DESC").Find(&coupons).Error
	return coupons, err
}

// GetActive retrieves all active coupons
func (r *CouponRepository) GetActive() ([]Coupon, error) {
	var coupons []Coupon
	now := time.Now()
	err := r.db.Where("is_active = ? AND start_date <= ? AND end_date >= ?", true, now, now).
		Order("created_at DESC").
		Find(&coupons).Error
	return coupons, err
}

// Update updates a coupon
func (r *CouponRepository) Update(coupon *Coupon) error {
	coupon.Code = strings.ToUpper(coupon.Code)
	return r.db.Save(coupon).Error
}

// Delete deletes a coupon
func (r *CouponRepository) Delete(id uint) error {
	return r.db.Delete(&Coupon{}, id).Error
}

// IncrementUsage increments the usage count of a coupon
func (r *CouponRepository) IncrementUsage(id uint) error {
	return r.db.Model(&Coupon{}).Where("id = ?", id).
		UpdateColumn("used_count", gorm.Expr("used_count + ?", 1)).Error
}

// ValidateCoupon validates if a coupon can be used
func (r *CouponRepository) ValidateCoupon(code string, orderAmount float64) (*Coupon, error) {
	coupon, err := r.GetByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("coupon not found")
		}
		return nil, err
	}

	// Check if coupon is active
	if !coupon.IsActive {
		return nil, errors.New("coupon is not active")
	}

	// Check date validity
	now := time.Now()
	if now.Before(coupon.StartDate) {
		return nil, errors.New("coupon is not yet valid")
	}
	if now.After(coupon.EndDate) {
		return nil, errors.New("coupon has expired")
	}

	// Check minimum purchase
	if orderAmount < coupon.MinPurchase {
		return nil, errors.New("order amount does not meet minimum purchase requirement")
	}

	// Check usage limit
	if coupon.UsageLimit > 0 && coupon.UsedCount >= coupon.UsageLimit {
		return nil, errors.New("coupon usage limit reached")
	}

	return coupon, nil
}

// CalculateDiscount calculates the discount amount based on coupon
func (coupon *Coupon) CalculateDiscount(orderAmount float64) float64 {
	var discount float64

	if coupon.DiscountType == "percentage" {
		discount = orderAmount * (coupon.DiscountValue / 100)
		// Apply max discount if set
		if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
			discount = coupon.MaxDiscount
		}
	} else if coupon.DiscountType == "fixed" {
		discount = coupon.DiscountValue
		// Discount cannot exceed order amount
		if discount > orderAmount {
			discount = orderAmount
		}
	}

	return discount
}
