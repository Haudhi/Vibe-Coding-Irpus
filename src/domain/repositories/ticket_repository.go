package repositories

import (
	"github.com/company/ga-ticketing/src/domain/entities"
)

// TicketRepository defines the interface for ticket persistence operations
type TicketRepository interface {
	// Create saves a new ticket to the database
	Create(ticket *entities.Ticket) error

	// FindByID retrieves a ticket by their ID
	FindByID(id string) (*entities.Ticket, error)

	// FindByTicketNumber retrieves a ticket by their ticket number
	FindByTicketNumber(ticketNumber string) (*entities.Ticket, error)

	// FindByRequester retrieves all tickets created by a specific user
	FindByRequester(requesterID string) ([]*entities.Ticket, error)

	// FindByAssignedAdmin retrieves all tickets assigned to a specific admin
	FindByAssignedAdmin(adminID string) ([]*entities.Ticket, error)

	// FindByStatus retrieves all tickets with a specific status
	FindByStatus(status entities.TicketStatus) ([]*entities.Ticket, error)

	// FindByCategory retrieves all tickets in a specific category
	FindByCategory(category entities.TicketCategory) ([]*entities.Ticket, error)

	// FindRequiringApproval retrieves all tickets that require approval
	FindRequiringApproval() ([]*entities.Ticket, error)

	// FindAll retrieves all tickets with optional filtering
	FindAll(filters TicketFilters) ([]*entities.Ticket, error)

	// Update updates an existing ticket in the database
	Update(ticket *entities.Ticket) error

	// Delete removes a ticket from the database
	Delete(id string) error

	// GetNextSequence returns the next ticket number sequence for the given year
	GetNextSequence(year int) (int, error)
}

// TicketFilters defines filters for ticket queries
type TicketFilters struct {
	RequesterID     *string
	AssignedAdminID *string
	Status          *entities.TicketStatus
	Category        *entities.TicketCategory
	Priority        *entities.TicketPriority
	RequiresApproval *bool
	Limit           *int
	Offset          *int
}