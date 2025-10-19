package models

import (
	"time"

	"gorm.io/gorm"
)

// Profile represents extended user profile information
type Profile struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(20)"`
	Address     string    `json:"address" gorm:"type:text"`
	City        string    `json:"city" gorm:"type:varchar(100)"`
	Province    string    `json:"province" gorm:"type:varchar(100)"`
	PostalCode  string    `json:"postal_code" gorm:"type:varchar(10)"`
	Country     string    `json:"country" gorm:"type:varchar(100);default:'Indonesia'"`
	Avatar      string    `json:"avatar" gorm:"type:varchar(500)"`
	Bio         string    `json:"bio" gorm:"type:text"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender" gorm:"type:varchar(10)"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ProfileInput for creating/updating profile
type ProfileInput struct {
	PhoneNumber string     `json:"phone_number"`
	Address     string     `json:"address"`
	City        string     `json:"city"`
	Province    string     `json:"province"`
	PostalCode  string     `json:"postal_code"`
	Country     string     `json:"country"`
	Avatar      string     `json:"avatar"`
	Bio         string     `json:"bio"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender" binding:"omitempty,oneof=male female other"`
}

// TableName specifies the table name for Profile model
func (Profile) TableName() string {
	return "profiles"
}

// ProfileRepository handles database operations for profiles
type ProfileRepository struct {
	db *gorm.DB
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

// Create adds a new profile to database
func (r *ProfileRepository) Create(profile *Profile) error {
	return r.db.Create(profile).Error
}

// GetByUserID returns a profile by user ID from database
func (r *ProfileRepository) GetByUserID(userID uint) (*Profile, error) {
	var profile Profile
	err := r.db.Preload("User").Where("user_id = ?", userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// GetByID returns a profile by ID from database
func (r *ProfileRepository) GetByID(id uint) (*Profile, error) {
	var profile Profile
	err := r.db.Preload("User").First(&profile, id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// Update updates a profile in database
func (r *ProfileRepository) Update(profile *Profile) error {
	return r.db.Save(profile).Error
}

// Delete removes a profile from database
func (r *ProfileRepository) Delete(id uint) error {
	return r.db.Delete(&Profile{}, id).Error
}
