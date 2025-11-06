package repositories

import (
	"github.com/company/ga-ticketing/src/domain/entities"
)

// CommentRepository defines the interface for comment persistence operations
type CommentRepository interface {
	// Create saves a new comment to the database
	Create(comment *entities.Comment) error

	// FindByID retrieves a comment by their ID
	FindByID(id string) (*entities.Comment, error)

	// FindByTicketID retrieves all comments for a specific ticket
	FindByTicketID(ticketID string) ([]*entities.Comment, error)

	// FindByUserID retrieves all comments by a specific user
	FindByUserID(userID string) ([]*entities.Comment, error)

	// Update updates an existing comment in the database
	Update(comment *entities.Comment) error

	// Delete removes a comment from the database
	Delete(id string) error

	// DeleteByTicketID removes all comments for a specific ticket
	DeleteByTicketID(ticketID string) error
}