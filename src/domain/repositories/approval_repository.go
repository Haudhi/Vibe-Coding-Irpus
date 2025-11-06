package repositories

import (
	"github.com/company/ga-ticketing/src/domain/entities"
)

// ApprovalRepository defines the interface for approval persistence operations
type ApprovalRepository interface {
	// Create saves a new approval to the database
	Create(approval *entities.Approval) error

	// FindByID retrieves an approval by their ID
	FindByID(id string) (*entities.Approval, error)

	// FindByTicketID retrieves all approvals for a specific ticket
	FindByTicketID(ticketID string) ([]*entities.Approval, error)

	// FindByApproverID retrieves all approvals by a specific approver
	FindByApproverID(approverID string) ([]*entities.Approval, error)

	// FindPendingByApprover retrieves all pending approvals for a specific approver
	FindPendingByApprover(approverID string) ([]*entities.Approval, error)

	// FindPending retrieves all pending approvals in the system
	FindPending() ([]*entities.Approval, error)

	// FindByTicketAndApprover retrieves an approval for a specific ticket and approver
	FindByTicketAndApprover(ticketID, approverID string) (*entities.Approval, error)

	// Update updates an existing approval in the database
	Update(approval *entities.Approval) error

	// Delete removes an approval from the database
	Delete(id string) error

	// CheckExists checks if an approval exists for a specific ticket and approver
	CheckExists(ticketID, approverID string) (bool, error)
}