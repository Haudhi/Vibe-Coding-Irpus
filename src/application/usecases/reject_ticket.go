package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// RejectTicketUseCase implements the use case for rejecting a ticket
type RejectTicketUseCase struct {
	ticketService *services.TicketService
}

// NewRejectTicketUseCase creates a new RejectTicketUseCase
func NewRejectTicketUseCase(ticketService *services.TicketService) *RejectTicketUseCase {
	return &RejectTicketUseCase{
		ticketService: ticketService,
	}
}

// Execute rejects a ticket
func (uc *RejectTicketUseCase) Execute(ctx context.Context, ticketID, userID string, userRole entities.UserRole, req *dto.RejectTicketRequest) (*dto.TicketResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Check permissions - only approver or admin can reject tickets
	if userRole != entities.RoleApprover && userRole != entities.RoleAdmin {
		return nil, fmt.Errorf("only approvers or admins can reject tickets")
	}

	// Validate ticket is eligible for rejection
	if err := uc.ticketService.ValidateTicketForApproval(ctx, ticketID); err != nil {
		return nil, fmt.Errorf("ticket cannot be rejected: %w", err)
	}

	// Update ticket status to rejected with reason
	if err := uc.ticketService.UpdateTicketStatus(ctx, ticketID, entities.StatusRejected, req.Reason, userID); err != nil {
		return nil, fmt.Errorf("failed to reject ticket: %w", err)
	}

	// Get the updated ticket
	ticket, err := uc.ticketService.GetTicket(ctx, ticketID, userID, userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejected ticket: %w", err)
	}

	// Convert to response DTO
	response := dto.TicketFromEntity(ticket)
	return response, nil
}
