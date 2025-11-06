package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// TicketRepository defines the interface for ticket persistence
type TicketRepository interface {
	Create(ctx context.Context, ticket *entities.Ticket) error
	GetByID(ctx context.Context, id string) (*entities.Ticket, error)
	GetByTicketNumber(ctx context.Context, ticketNumber string) (*entities.Ticket, error)
	GetByRequesterID(ctx context.Context, requesterID string, limit, offset int) ([]*entities.Ticket, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Ticket, error)
	Update(ctx context.Context, ticket *entities.Ticket) error
	Delete(ctx context.Context, id string) error
	GetNextSequenceNumber(ctx context.Context, year int) (int, error)
}

// UserRepository defines the interface for user persistence
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entities.User, error)
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	Update(ctx context.Context, user *entities.User) error
}

// TicketService provides business logic for ticket operations
type TicketService struct {
	ticketRepo TicketRepository
	userRepo   UserRepository
}

// NewTicketService creates a new TicketService
func NewTicketService(ticketRepo TicketRepository, userRepo UserRepository) *TicketService {
	return &TicketService{
		ticketRepo: ticketRepo,
		userRepo:   userRepo,
	}
}

// CreateTicket creates a new ticket
func (s *TicketService) CreateTicket(
	ctx context.Context,
	title, description string,
	category entities.TicketCategory,
	priority entities.TicketPriority,
	estimatedCost *valueobjects.Money,
	requesterID string,
) (*entities.Ticket, error) {
	// Validate requester exists
	requester, err := s.userRepo.GetByID(ctx, requesterID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user can create tickets
	if !requester.HasPermission("create_ticket") {
		return nil, errors.New("user does not have permission to create tickets")
	}

	// Create ticket
	ticket, err := entities.NewTicket(
		title,
		description,
		category,
		priority,
		estimatedCost,
		requesterID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Note: Ticket number generation is handled within the entity itself
	// In a real implementation, we might want to ensure uniqueness here

	// Save ticket
	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to save ticket: %w", err)
	}

	return ticket, nil
}

// GetTicket retrieves a ticket by ID with access control
func (s *TicketService) GetTicket(ctx context.Context, ticketID, userID string, userRole entities.UserRole) (*entities.Ticket, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	// Check access permissions
	if !s.canUserViewTicket(ticket, userID, userRole) {
		return nil, errors.New("access denied")
	}

	return ticket, nil
}

// GetUserTickets retrieves tickets for a specific user
func (s *TicketService) GetUserTickets(ctx context.Context, userID string, limit, offset int) ([]*entities.Ticket, error) {
	tickets, err := s.ticketRepo.GetByRequesterID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tickets: %w", err)
	}

	return tickets, nil
}

// GetAllTickets retrieves all tickets (admin only)
func (s *TicketService) GetAllTickets(ctx context.Context, limit, offset int) ([]*entities.Ticket, error) {
	tickets, err := s.ticketRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all tickets: %w", err)
	}

	return tickets, nil
}

// AssignTicket assigns a ticket to an admin
func (s *TicketService) AssignTicket(ctx context.Context, ticketID, adminID string) error {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Get admin user
	admin, err := s.userRepo.GetByID(ctx, adminID)
	if err != nil {
		return fmt.Errorf("admin not found: %w", err)
	}

	// Check if assignee is admin
	if admin.GetRole() != entities.RoleAdmin {
		return errors.New("assignee must be an admin")
	}

	// Check if ticket is already assigned
	if ticket.GetAssignedAdminID() != nil {
		return errors.New("ticket is already assigned")
	}

	// Assign ticket
	if err := ticket.AssignToAdmin(adminID); err != nil {
		return fmt.Errorf("failed to assign ticket: %w", err)
	}

	// Save changes
	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return fmt.Errorf("failed to save assigned ticket: %w", err)
	}

	return nil
}

// ReassignTicket reassigns a ticket to a different admin
func (s *TicketService) ReassignTicket(ctx context.Context, ticketID, newAdminID string) error {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Get new admin user
	newAdmin, err := s.userRepo.GetByID(ctx, newAdminID)
	if err != nil {
		return fmt.Errorf("admin not found: %w", err)
	}

	// Check if new assignee is admin
	if newAdmin.GetRole() != entities.RoleAdmin {
		return errors.New("new assignee must be an admin")
	}

	// Reassign ticket
	if err := ticket.ReassignToAdmin(newAdminID); err != nil {
		return fmt.Errorf("failed to reassign ticket: %w", err)
	}

	// Save changes
	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return fmt.Errorf("failed to save reassigned ticket: %w", err)
	}

	return nil
}

// UpdateTicketStatus updates the status of a ticket
func (s *TicketService) UpdateTicketStatus(
	ctx context.Context,
	ticketID string,
	newStatus entities.TicketStatus,
	reason, updatedBy string,
) error {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Update status
	if err := ticket.SetStatus(newStatus, reason, updatedBy); err != nil {
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	// Save changes
	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return fmt.Errorf("failed to save updated ticket: %w", err)
	}

	return nil
}

// AddCommentToTicket adds a comment to a ticket
func (s *TicketService) AddCommentToTicket(
	ctx context.Context,
	ticketID, content, userID string,
) (*entities.Comment, error) {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	// Add comment
	comment, err := ticket.AddComment(content, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	// Save changes
	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, fmt.Errorf("failed to save ticket with comment: %w", err)
	}

	return comment, nil
}

// DeleteTicket deletes a ticket (soft delete via status update)
func (s *TicketService) DeleteTicket(ctx context.Context, ticketID, userID string) error {
	// Get ticket
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	// Check if user can delete this ticket
	if ticket.GetRequesterID() != userID {
		return errors.New("only ticket creator can delete ticket")
	}

	// Check if ticket can be deleted (only in certain statuses)
	status := ticket.GetStatus()
	if status != entities.StatusPending && status != entities.StatusRejected {
		return errors.New("ticket can only be deleted when pending or rejected")
	}

	// Soft delete by setting status to closed
	if err := ticket.SetStatus(entities.StatusClosed, "Ticket deleted by requester", userID); err != nil {
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	// Save changes
	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return fmt.Errorf("failed to save deleted ticket: %w", err)
	}

	return nil
}

// GetTicketsByStatus retrieves tickets by status
func (s *TicketService) GetTicketsByStatus(
	ctx context.Context,
	status entities.TicketStatus,
	limit, offset int,
) ([]*entities.Ticket, error) {
	// This would require a new method in the repository
	// For now, we'll return an empty slice
	return []*entities.Ticket{}, nil
}

// GetTicketsByCategory retrieves tickets by category
func (s *TicketService) GetTicketsByCategory(
	ctx context.Context,
	category entities.TicketCategory,
	limit, offset int,
) ([]*entities.Ticket, error) {
	// This would require a new method in the repository
	// For now, we'll return an empty slice
	return []*entities.Ticket{}, nil
}

// Helper methods

func (s *TicketService) canUserViewTicket(
	ticket *entities.Ticket,
	userID string,
	userRole entities.UserRole,
) bool {
	// Admin can view any ticket
	if userRole == entities.RoleAdmin {
		return true
	}

	// User can view their own tickets
	if ticket.GetRequesterID() == userID {
		return true
	}

	// Assigned admin can view the ticket
	if ticket.GetAssignedAdminID() != nil && *ticket.GetAssignedAdminID() == userID {
		return true
	}

	// Approver can view tickets requiring approval
	if userRole == entities.RoleApprover && ticket.RequiresApproval() {
		return true
	}

	return false
}

// GetTicketsNeedingApproval retrieves tickets that need approval
func (s *TicketService) GetTicketsNeedingApproval(ctx context.Context, limit, offset int) ([]*entities.Ticket, error) {
	// This would require a new method in the repository
	// For now, we'll return an empty slice
	return []*entities.Ticket{}, nil
}

// GetTicketsAssignedToAdmin retrieves tickets assigned to a specific admin
func (s *TicketService) GetTicketsAssignedToAdmin(ctx context.Context, adminID string, limit, offset int) ([]*entities.Ticket, error) {
	// This would require a new method in the repository
	// For now, we'll return an empty slice
	return []*entities.Ticket{}, nil
}

// ValidateTicketForApproval checks if a ticket can be approved
func (s *TicketService) ValidateTicketForApproval(ctx context.Context, ticketID string) error {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("ticket not found: %w", err)
	}

	if !ticket.RequiresApproval() {
		return errors.New("ticket does not require approval")
	}

	status := ticket.GetStatus()
	if status != entities.StatusWaitingApproval {
		return fmt.Errorf("ticket is not waiting for approval (current status: %s)", status)
	}

	return nil
}

// GetTicketStatistics returns statistics about tickets
func (s *TicketService) GetTicketStatistics(ctx context.Context) (map[string]int, error) {
	// This would require aggregation queries
	// For now, return empty statistics
	stats := map[string]int{
		"total":       0,
		"pending":     0,
		"in_progress": 0,
		"completed":   0,
		"closed":      0,
	}
	return stats, nil
}