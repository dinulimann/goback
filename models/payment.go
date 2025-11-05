package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment represents a payment transaction
type Payment struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	OrderID           uint           `gorm:"not null;index" json:"order_id"`
	Order             Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	PaymentMethod     string         `gorm:"type:varchar(50);not null" json:"payment_method"` // midtrans, bank_transfer, cod, etc.
	TransactionID     string         `gorm:"type:varchar(255);uniqueIndex" json:"transaction_id"` // Midtrans transaction_id
	SnapToken         string         `gorm:"type:text" json:"snap_token,omitempty"` // Midtrans Snap token
	GrossAmount       float64        `gorm:"type:decimal(15,2);not null" json:"gross_amount"`
	PaymentStatus     string         `gorm:"type:varchar(50);not null;default:'pending'" json:"payment_status"` // pending, paid, failed, expired
	PaymentProvider   string         `gorm:"type:varchar(50)" json:"payment_provider"` // midtrans, manual, etc.
	PaidAt            *time.Time     `json:"paid_at,omitempty"`
	ExpiredAt         *time.Time     `json:"expired_at,omitempty"`

	// Midtrans specific fields
	TransactionStatus string         `gorm:"type:varchar(50)" json:"transaction_status,omitempty"` // capture, settlement, pending, deny, cancel, expire
	FraudStatus       string         `gorm:"type:varchar(50)" json:"fraud_status,omitempty"` // accept, challenge, deny
	VANumber          string         `gorm:"type:varchar(100)" json:"va_number,omitempty"` // Virtual Account number
	Bank              string         `gorm:"type:varchar(50)" json:"bank,omitempty"` // BCA, BNI, Mandiri, etc.

	// Additional info
	Metadata          string         `gorm:"type:text" json:"metadata,omitempty"` // JSON string for additional data
	ErrorMessage      string         `gorm:"type:text" json:"error_message,omitempty"`

	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// PaymentRepository defines payment data access methods
type PaymentRepository interface {
	Create(payment *Payment) error
	GetByID(id uint) (*Payment, error)
	GetByOrderID(orderID uint) (*Payment, error)
	GetByTransactionID(transactionID string) (*Payment, error)
	Update(payment *Payment) error
	Delete(id uint) error
	GetAll(limit, offset int) ([]Payment, int64, error)
	GetByStatus(status string, limit, offset int) ([]Payment, int64, error)
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// Create creates a new payment
func (r *paymentRepository) Create(payment *Payment) error {
	return r.db.Create(payment).Error
}

// GetByID retrieves a payment by ID
func (r *paymentRepository) GetByID(id uint) (*Payment, error) {
	var payment Payment
	err := r.db.Preload("Order").First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByOrderID retrieves a payment by order ID
func (r *paymentRepository) GetByOrderID(orderID uint) (*Payment, error) {
	var payment Payment
	err := r.db.Where("order_id = ?", orderID).Preload("Order").First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByTransactionID retrieves a payment by transaction ID
func (r *paymentRepository) GetByTransactionID(transactionID string) (*Payment, error) {
	var payment Payment
	err := r.db.Where("transaction_id = ?", transactionID).Preload("Order").First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// Update updates a payment
func (r *paymentRepository) Update(payment *Payment) error {
	return r.db.Save(payment).Error
}

// Delete soft deletes a payment
func (r *paymentRepository) Delete(id uint) error {
	return r.db.Delete(&Payment{}, id).Error
}

// GetAll retrieves all payments with pagination
func (r *paymentRepository) GetAll(limit, offset int) ([]Payment, int64, error) {
	var payments []Payment
	var total int64

	if err := r.db.Model(&Payment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Order").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error

	return payments, total, err
}

// GetByStatus retrieves payments by status with pagination
func (r *paymentRepository) GetByStatus(status string, limit, offset int) ([]Payment, int64, error) {
	var payments []Payment
	var total int64

	query := r.db.Model(&Payment{}).Where("payment_status = ?", status)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Order").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error

	return payments, total, err
}
