package repositories

import (
	"github.com/company/ga-ticketing/src/domain/entities"
)

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	// Create saves a new user to the database
	Create(user *entities.User) error

	// FindByID retrieves a user by their ID
	FindByID(id string) (*entities.User, error)

	// FindByEmail retrieves a user by their email
	FindByEmail(email string) (*entities.User, error)

	// FindByEmployeeID retrieves a user by their employee ID
	FindByEmployeeID(employeeID string) (*entities.User, error)

	// FindAll retrieves all users with optional filtering
	FindAll(role *entities.UserRole, isActive *bool) ([]*entities.User, error)

	// Update updates an existing user in the database
	Update(user *entities.User) error

	// Delete removes a user from the database (soft delete by deactivating)
	Delete(id string) error

	// Exists checks if a user exists with the given email or employee ID
	Exists(email, employeeID string) (bool, error)
}