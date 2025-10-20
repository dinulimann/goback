package models

import (
	"time"

	"gorm.io/gorm"
)

// Wishlist represents a user's wishlist item
type Wishlist struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// WishlistInput for adding to wishlist
type WishlistInput struct {
	ProductID uint `json:"product_id" binding:"required"`
}

// TableName specifies the table name for Wishlist model
func (Wishlist) TableName() string {
	return "wishlists"
}

// WishlistRepository handles database operations for wishlists
type WishlistRepository struct {
	db *gorm.DB
}

// NewWishlistRepository creates a new wishlist repository
func NewWishlistRepository(db *gorm.DB) *WishlistRepository {
	return &WishlistRepository{
		db: db,
	}
}

// Create adds a product to wishlist
func (r *WishlistRepository) Create(wishlist *Wishlist) error {
	return r.db.Create(wishlist).Error
}

// GetByUserID returns wishlist items by user ID
func (r *WishlistRepository) GetByUserID(userID uint) ([]Wishlist, error) {
	var wishlists []Wishlist
	err := r.db.Preload("Product.Category").Where("user_id = ?", userID).Order("created_at desc").Find(&wishlists).Error
	return wishlists, err
}

// Delete removes a product from wishlist
func (r *WishlistRepository) Delete(userID, productID uint) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&Wishlist{}).Error
}

// Exists checks if product is in user's wishlist
func (r *WishlistRepository) Exists(userID, productID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Wishlist{}).Where("user_id = ? AND product_id = ?", userID, productID).Count(&count).Error
	return count > 0, err
}
