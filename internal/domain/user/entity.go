package user

import (
	"clean-arch-go/internal/errors"
	"strings"
	"time"
)

// User represents a user entity
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate validates the user entity
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.NewValidationError("email", "Email is required")
	}
	if u.Password == "" {
		return errors.NewValidationError("password", "Password is required")
	}
	if u.Name == "" {
		return errors.NewValidationError("name", "Name is required")
	}
	return nil
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenResponse represents the response for login
type TokenResponse struct {
	Token string `json:"token"`
}

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// Validate validates the create user input
func (i *CreateUserInput) Validate() error {
	if i.Email == "" {
		return errors.NewValidationError("email", "Email is required")
	}
	if i.Password == "" || len(i.Password) < 6 {
		return errors.NewValidationError("password", "Password must be at least 6 characters")
	}
	if i.Name == "" {
		return errors.NewValidationError("name", "Name is required")
	}
	return nil
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	ID    string `json:"id" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
	Name  string `json:"name" binding:"omitempty"`
}

// Validate validates the update user input
func (i *UpdateUserInput) Validate() error {
	if i.Email != "" && !isValidEmail(i.Email) {
		return errors.NewValidationError("email", "Invalid email format")
	}
	return nil
}

func isValidEmail(email string) bool {
	// Simple email validation
	if len(email) < 3 || !strings.Contains(email, "@") {
		return false
	}
	return true
}
