package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" binding:"required" gorm:"type:varchar(200);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" binding:"required" gorm:"type:decimal(10,2);not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	SKU         string    `json:"sku" gorm:"type:varchar(100);uniqueIndex"`
	CategoryID  uint      `json:"category_id" binding:"required" gorm:"not null"`
	ImageURL    string    `json:"image_url" gorm:"type:varchar(500)"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Category   Category    `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
}

// ProductInput for creating/updating product
type ProductInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required,gt=0"`
	Stock       int      `json:"stock" binding:"gte=0"`
	SKU         string   `json:"sku"`
	CategoryID  uint     `json:"category_id" binding:"required"`
	ImageURL    string   `json:"image_url"`
	IsActive    *bool    `json:"is_active"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

// ProductRepository handles database operations for products
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Create adds a new product to database
func (r *ProductRepository) Create(product *Product) error {
	return r.db.Create(product).Error
}

// GetAll returns all products from database
func (r *ProductRepository) GetAll() ([]Product, error) {
	var products []Product
	err := r.db.Preload("Category").Order("id desc").Find(&products).Error
	return products, err
}

// GetByID returns a product by ID from database
func (r *ProductRepository) GetByID(id uint) (*Product, error) {
	var product Product
	err := r.db.Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetBySKU returns a product by SKU from database
func (r *ProductRepository) GetBySKU(sku string) (*Product, error) {
	var product Product
	err := r.db.Preload("Category").Where("sku = ?", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetByCategory returns products by category ID from database
func (r *ProductRepository) GetByCategory(categoryID uint) ([]Product, error) {
	var products []Product
	err := r.db.Preload("Category").Where("category_id = ?", categoryID).Order("id desc").Find(&products).Error
	return products, err
}

// Update updates a product in database
func (r *ProductRepository) Update(product *Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product from database
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}

// UpdateStock updates product stock
func (r *ProductRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&Product{}).Where("id = ?", id).Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
