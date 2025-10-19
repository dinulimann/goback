package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" binding:"required" gorm:"type:varchar(100);not null;uniqueIndex"`
	Description string    `json:"description" gorm:"type:text"`
	Slug        string    `json:"slug" gorm:"type:varchar(100);uniqueIndex;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

// CategoryInput for creating/updating category
type CategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Slug        string `json:"slug" binding:"required"`
	IsActive    *bool  `json:"is_active"`
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "categories"
}

// CategoryRepository handles database operations for categories
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

// Create adds a new category to database
func (r *CategoryRepository) Create(category *Category) error {
	return r.db.Create(category).Error
}

// GetAll returns all categories from database
func (r *CategoryRepository) GetAll() ([]Category, error) {
	var categories []Category
	err := r.db.Preload("Products").Order("id desc").Find(&categories).Error
	return categories, err
}

// GetByID returns a category by ID from database
func (r *CategoryRepository) GetByID(id uint) (*Category, error) {
	var category Category
	err := r.db.Preload("Products").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// GetBySlug returns a category by slug from database
func (r *CategoryRepository) GetBySlug(slug string) (*Category, error) {
	var category Category
	err := r.db.Preload("Products").Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// Update updates a category in database
func (r *CategoryRepository) Update(category *Category) error {
	return r.db.Save(category).Error
}

// Delete removes a category from database
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&Category{}, id).Error
}
