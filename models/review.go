package models

import (
	"time"

	"gorm.io/gorm"
)

// Review represents a product review
type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	Rating    int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment   string    `json:"comment" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ReviewInput for creating/updating review
type ReviewInput struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}

// TableName specifies the table name for Review model
func (Review) TableName() string {
	return "reviews"
}

// ReviewRepository handles database operations for reviews
type ReviewRepository struct {
	db *gorm.DB
}

// NewReviewRepository creates a new review repository
func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

// Create adds a new review to database
func (r *ReviewRepository) Create(review *Review) error {
	return r.db.Create(review).Error
}

// GetByProductID returns reviews by product ID
func (r *ReviewRepository) GetByProductID(productID uint) ([]Review, error) {
	var reviews []Review
	err := r.db.Preload("User").Where("product_id = ?", productID).Order("created_at desc").Find(&reviews).Error
	return reviews, err
}

// GetByUserID returns reviews by user ID
func (r *ReviewRepository) GetByUserID(userID uint) ([]Review, error) {
	var reviews []Review
	err := r.db.Preload("Product").Where("user_id = ?", userID).Order("created_at desc").Find(&reviews).Error
	return reviews, err
}

// GetByID returns a review by ID
func (r *ReviewRepository) GetByID(id uint) (*Review, error) {
	var review Review
	err := r.db.Preload("Product").Preload("User").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// Update updates a review
func (r *ReviewRepository) Update(review *Review) error {
	return r.db.Save(review).Error
}

// Delete removes a review
func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&Review{}, id).Error
}

// GetAverageRating gets average rating for a product
func (r *ReviewRepository) GetAverageRating(productID uint) (float64, error) {
	var avg float64
	err := r.db.Model(&Review{}).Where("product_id = ?", productID).Select("AVG(rating)").Scan(&avg).Error
	return avg, err
}

// CountByProductID counts reviews for a product
func (r *ReviewRepository) CountByProductID(productID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Review{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}
