package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);not null"` // Price at time of order
	Subtotal  float64   `json:"subtotal" gorm:"type:decimal(12,2);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Order   Order   `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

// OrderItemInput for creating order item
type OrderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

// TableName specifies the table name for OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderItemRepository handles database operations for order items
type OrderItemRepository struct {
	db *gorm.DB
}

// NewOrderItemRepository creates a new order item repository
func NewOrderItemRepository(db *gorm.DB) *OrderItemRepository {
	return &OrderItemRepository{
		db: db,
	}
}

// Create adds a new order item to database
func (r *OrderItemRepository) Create(orderItem *OrderItem) error {
	return r.db.Create(orderItem).Error
}

// GetByOrderID returns order items by order ID from database
func (r *OrderItemRepository) GetByOrderID(orderID uint) ([]OrderItem, error) {
	var orderItems []OrderItem
	err := r.db.Preload("Product").Where("order_id = ?", orderID).Find(&orderItems).Error
	return orderItems, err
}

// GetByID returns an order item by ID from database
func (r *OrderItemRepository) GetByID(id uint) (*OrderItem, error) {
	var orderItem OrderItem
	err := r.db.Preload("Product").Preload("Order").First(&orderItem, id).Error
	if err != nil {
		return nil, err
	}
	return &orderItem, nil
}

// Update updates an order item in database
func (r *OrderItemRepository) Update(orderItem *OrderItem) error {
	return r.db.Save(orderItem).Error
}

// Delete removes an order item from database
func (r *OrderItemRepository) Delete(id uint) error {
	return r.db.Delete(&OrderItem{}, id).Error
}

// DeleteByOrderID removes all order items for a specific order
func (r *OrderItemRepository) DeleteByOrderID(orderID uint) error {
	return r.db.Where("order_id = ?", orderID).Delete(&OrderItem{}).Error
}
