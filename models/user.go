package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserStore is an in-memory store for users (in production, use a database)
type UserStore struct {
	users  map[int]*User
	nextID int
}

// NewUserStore creates a new user store
func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

// Create adds a new user
func (s *UserStore) Create(user *User) *User {
	user.ID = s.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	s.nextID++
	return user
}

// GetAll returns all users
func (s *UserStore) GetAll() []*User {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// GetByID returns a user by ID
func (s *UserStore) GetByID(id int) (*User, bool) {
	user, exists := s.users[id]
	return user, exists
}

// Update updates a user
func (s *UserStore) Update(id int, updatedUser *User) (*User, bool) {
	user, exists := s.users[id]
	if !exists {
		return nil, false
	}

	user.Name = updatedUser.Name
	user.Email = updatedUser.Email
	user.UpdatedAt = time.Now()

	return user, true
}

// Delete removes a user
func (s *UserStore) Delete(id int) bool {
	_, exists := s.users[id]
	if !exists {
		return false
	}
	delete(s.users, id)
	return true
}
