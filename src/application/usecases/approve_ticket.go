package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// ApproveTicketUseCase implements the use case for approving a ticket
type ApproveTicketUseCase struct {
	ticketService *services.TicketService
}

// NewApproveTicketUseCase creates a new ApproveTicketUseCase
func NewApproveTicketUseCase(ticketService *services.TicketService) *ApproveTicketUseCase {
	return &ApproveTicketUseCase{
		ticketService: ticketService,
	}
}

// Execute approves a ticket
func (uc *ApproveTicketUseCase) Execute(ctx context.Context, ticketID, userID string, userRole entities.UserRole, req *dto.ApproveTicketRequest) (*dto.TicketResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check permissions - only approver or admin can approve tickets
	if userRole != entities.RoleApprover && userRole != entities.RoleAdmin {
		return nil, fmt.Errorf("only approvers or admins can approve tickets")
	}

	// Validate ticket is eligible for approval
	if err := uc.ticketService.ValidateTicketForApproval(ctx, ticketID); err != nil {
		return nil, fmt.Errorf("ticket cannot be approved: %w", err)
	}

	// Update ticket status to approved
	comments := "Ticket approved"
	if req.Comments != "" {
		comments = req.Comments
	}

	if err := uc.ticketService.UpdateTicketStatus(ctx, ticketID, entities.StatusApproved, comments, userID); err != nil {
		return nil, fmt.Errorf("failed to approve ticket: %w", err)
	}

	// Get the updated ticket
	ticket, err := uc.ticketService.GetTicket(ctx, ticketID, userID, userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved ticket: %w", err)
	}

	// Convert to response DTO
	response := dto.TicketFromEntity(ticket)
	return response, nil
}
