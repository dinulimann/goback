package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" binding:"required" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" binding:"required,email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"` // Hidden from JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterInput for user registration
type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput for user login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// UserRepository handles database operations for users
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create adds a new user to database
func (r *UserRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

// GetAll returns all users from database
func (r *UserRepository) GetAll() ([]User, error) {
	var users []User
	err := r.db.Order("id desc").Find(&users).Error
	return users, err
}

// GetByID returns a user by ID from database
func (r *UserRepository) GetByID(id uint) (*User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user in database
func (r *UserRepository) Update(user *User) error {
	return r.db.Save(user).Error
}

// Delete removes a user from database
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&User{}, id).Error
}

// GetByEmail returns a user by email from database
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
