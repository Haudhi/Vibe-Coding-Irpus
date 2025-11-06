package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// AssignTicketUseCase implements the use case for assigning a ticket to an admin
type AssignTicketUseCase struct {
	ticketService *services.TicketService
}

// NewAssignTicketUseCase creates a new AssignTicketUseCase
func NewAssignTicketUseCase(ticketService *services.TicketService) *AssignTicketUseCase {
	return &AssignTicketUseCase{
		ticketService: ticketService,
	}
}

// Execute assigns a ticket to an admin
func (uc *AssignTicketUseCase) Execute(ctx context.Context, ticketID, userID string, userRole entities.UserRole, req *dto.AssignTicketRequest) (*dto.TicketResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check permissions - only admin can assign tickets
	if userRole != entities.RoleAdmin {
		return nil, fmt.Errorf("only admins can assign tickets")
	}

	// Assign the ticket
	if err := uc.ticketService.AssignTicket(ctx, ticketID, req.AdminID); err != nil {
		return nil, fmt.Errorf("failed to assign ticket: %w", err)
	}

	// Get the updated ticket
	ticket, err := uc.ticketService.GetTicket(ctx, ticketID, userID, userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get assigned ticket: %w", err)
	}

	// Convert to response DTO
	response := dto.TicketFromEntity(ticket)
	return response, nil
}
