package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderStatus represents order status
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID                 uint        `json:"id" gorm:"primaryKey"`
	UserID             uint        `json:"user_id" gorm:"not null"`
	OrderNumber        string      `json:"order_number" gorm:"type:varchar(50);uniqueIndex;not null"`
	Status             OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	TotalAmount        float64     `json:"total_amount" gorm:"type:decimal(12,2);not null"`
	CouponID           *uint       `json:"coupon_id" gorm:"default:null"`
	CouponCode         string      `json:"coupon_code" gorm:"type:varchar(50)"`
	DiscountAmount     float64     `json:"discount_amount" gorm:"type:decimal(12,2);default:0"`
	FinalAmount        float64     `json:"final_amount" gorm:"type:decimal(12,2);not null"`
	ShippingAddress    string      `json:"shipping_address" gorm:"type:text"`
	ShippingCity       string      `json:"shipping_city" gorm:"type:varchar(100)"`
	ShippingProvince   string      `json:"shipping_province" gorm:"type:varchar(100)"`
	ShippingPostalCode string      `json:"shipping_postal_code" gorm:"type:varchar(10)"`
	PaymentMethod      string      `json:"payment_method" gorm:"type:varchar(50)"`
	PaymentStatus      string      `json:"payment_status" gorm:"type:varchar(20);default:'unpaid'"`
	Notes              string      `json:"notes" gorm:"type:text"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`

	// Relations
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Coupon     *Coupon     `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

// OrderInput for creating order
type OrderInput struct {
	ShippingAddress    string           `json:"shipping_address" binding:"required"`
	ShippingCity       string           `json:"shipping_city" binding:"required"`
	ShippingProvince   string           `json:"shipping_province" binding:"required"`
	ShippingPostalCode string           `json:"shipping_postal_code" binding:"required"`
	PaymentMethod      string           `json:"payment_method" binding:"required"`
	Notes              string           `json:"notes"`
	CouponCode         string           `json:"coupon_code"`
	Items              []OrderItemInput `json:"items" binding:"required,min=1,dive"`
}

// OrderUpdateStatusInput for updating order status
type OrderUpdateStatusInput struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending processing shipped delivered cancelled"`
}

// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}

// OrderRepository handles database operations for orders
type OrderRepository struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

// Create adds a new order to database
func (r *OrderRepository) Create(order *Order) error {
	return r.db.Create(order).Error
}

// GetAll returns all orders from database
func (r *OrderRepository) GetAll() ([]Order, error) {
	var orders []Order
	err := r.db.Preload("User").Preload("OrderItems.Product").Order("id desc").Find(&orders).Error
	return orders, err
}

// GetByID returns an order by ID from database
func (r *OrderRepository) GetByID(id uint) (*Order, error) {
	var order Order
	err := r.db.Preload("User").Preload("OrderItems.Product").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByOrderNumber returns an order by order number from database
func (r *OrderRepository) GetByOrderNumber(orderNumber string) (*Order, error) {
	var order Order
	err := r.db.Preload("User").Preload("OrderItems.Product").Where("order_number = ?", orderNumber).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetByUserID returns orders by user ID from database
func (r *OrderRepository) GetByUserID(userID uint) ([]Order, error) {
	var orders []Order
	err := r.db.Preload("User").Preload("OrderItems.Product").Where("user_id = ?", userID).Order("id desc").Find(&orders).Error
	return orders, err
}

// Update updates an order in database
func (r *OrderRepository) Update(order *Order) error {
	return r.db.Save(order).Error
}

// UpdateStatus updates order status
func (r *OrderRepository) UpdateStatus(id uint, status OrderStatus) error {
	return r.db.Model(&Order{}).Where("id = ?", id).Update("status", status).Error
}

// Delete removes an order from database
func (r *OrderRepository) Delete(id uint) error {
	return r.db.Delete(&Order{}, id).Error
}
