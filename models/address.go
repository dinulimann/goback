package models

import (
	"time"

	"gorm.io/gorm"
)

// Address represents a saved shipping address
type Address struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	Label        string    `gorm:"type:varchar(50)" json:"label"` // e.g., "Home", "Office"
	RecipientName string   `gorm:"type:varchar(100);not null" json:"recipient_name"`
	PhoneNumber  string    `gorm:"type:varchar(20);not null" json:"phone_number"`
	Address      string    `gorm:"type:text;not null" json:"address"`
	City         string    `gorm:"type:varchar(100);not null" json:"city"`
	Province     string    `gorm:"type:varchar(100);not null" json:"province"`
	PostalCode   string    `gorm:"type:varchar(10);not null" json:"postal_code"`
	IsDefault    bool      `gorm:"default:false" json:"is_default"`
	Notes        string    `gorm:"type:text" json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// AddressInput for creating/updating address
type AddressInput struct {
	Label         string `json:"label"`
	RecipientName string `json:"recipient_name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	Address       string `json:"address" binding:"required"`
	City          string `json:"city" binding:"required"`
	Province      string `json:"province" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	IsDefault     bool   `json:"is_default"`
	Notes         string `json:"notes"`
}

// TableName specifies the table name for Address model
func (Address) TableName() string {
	return "addresses"
}

// AddressRepository handles database operations for addresses
type AddressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(db *gorm.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

// Create creates a new address
func (r *AddressRepository) Create(address *Address) error {
	// If this is set as default, unset other defaults for this user
	if address.IsDefault {
		if err := r.UnsetDefaultForUser(address.UserID); err != nil {
			return err
		}
	}
	return r.db.Create(address).Error
}

// GetByID retrieves an address by ID
func (r *AddressRepository) GetByID(id uint) (*Address, error) {
	var address Address
	err := r.db.First(&address, id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByUserID retrieves all addresses for a user
func (r *AddressRepository) GetByUserID(userID uint) ([]Address, error) {
	var addresses []Address
	err := r.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	return addresses, err
}

// GetDefaultByUserID retrieves the default address for a user
func (r *AddressRepository) GetDefaultByUserID(userID uint) (*Address, error) {
	var address Address
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Update updates an address
func (r *AddressRepository) Update(address *Address) error {
	// If this is set as default, unset other defaults for this user
	if address.IsDefault {
		if err := r.UnsetDefaultForUser(address.UserID); err != nil {
			return err
		}
	}
	return r.db.Save(address).Error
}

// Delete deletes an address
func (r *AddressRepository) Delete(id uint) error {
	return r.db.Delete(&Address{}, id).Error
}

// SetDefault sets an address as default and unsets others
func (r *AddressRepository) SetDefault(id uint, userID uint) error {
	// Start transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Unset all defaults for this user
	if err := tx.Model(&Address{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set this address as default
	if err := tx.Model(&Address{}).Where("id = ? AND user_id = ?", id, userID).Update("is_default", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UnsetDefaultForUser unsets all default addresses for a user
func (r *AddressRepository) UnsetDefaultForUser(userID uint) error {
	return r.db.Model(&Address{}).Where("user_id = ?", userID).Update("is_default", false).Error
}

// CheckOwnership checks if an address belongs to a user
func (r *AddressRepository) CheckOwnership(addressID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Address{}).Where("id = ? AND user_id = ?", addressID, userID).Count(&count).Error
	return count > 0, err
}
