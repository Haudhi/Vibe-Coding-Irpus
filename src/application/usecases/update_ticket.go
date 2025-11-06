package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// UpdateTicketUseCase implements the use case for updating a ticket
type UpdateTicketUseCase struct {
	ticketService *services.TicketService
}

// NewUpdateTicketUseCase creates a new UpdateTicketUseCase
func NewUpdateTicketUseCase(ticketService *services.TicketService) *UpdateTicketUseCase {
	return &UpdateTicketUseCase{
		ticketService: ticketService,
	}
}

// Execute updates a ticket
func (uc *UpdateTicketUseCase) Execute(ctx context.Context, ticketID, userID string, userRole entities.UserRole, req *dto.UpdateTicketRequest) (*dto.TicketResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Get the ticket
	ticket, err := uc.ticketService.GetTicket(ctx, ticketID, userID, userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Check permissions - only admin can update tickets
	if userRole != entities.RoleAdmin {
		return nil, fmt.Errorf("only admins can update tickets")
	}

	// Update title if provided
	if req.Title != nil {
		if err := ticket.SetTitle(*req.Title); err != nil {
			return nil, fmt.Errorf("failed to update title: %w", err)
		}
	}

	// Update description if provided
	if req.Description != nil {
		if err := ticket.SetDescription(*req.Description); err != nil {
			return nil, fmt.Errorf("failed to update description: %w", err)
		}
	}

	// Update priority if provided
	if req.Priority != nil {
		priority, err := parseTicketPriority(*req.Priority)
		if err != nil {
			return nil, fmt.Errorf("invalid priority: %w", err)
		}
		ticket.SetPriority(priority)
	}

	// Update actual cost if provided
	if req.ActualCost != nil {
		if err := ticket.SetActualCost(valueobjects.NewMoney(*req.ActualCost)); err != nil {
			return nil, fmt.Errorf("failed to update actual cost: %w", err)
		}
	}

	// Update status if provided
	if req.Status != nil {
		status, err := parseTicketStatus(*req.Status)
		if err != nil {
			return nil, fmt.Errorf("invalid status: %w", err)
		}

		reason := ""
		if req.Reason != nil {
			reason = *req.Reason
		}

		updatedBy := userID
		if req.UpdatedBy != nil {
			updatedBy = *req.UpdatedBy
		}

		if err := uc.ticketService.UpdateTicketStatus(ctx, ticketID, status, reason, updatedBy); err != nil {
			return nil, fmt.Errorf("failed to update status: %w", err)
		}

		// Refresh the ticket to get updated status
		ticket, err = uc.ticketService.GetTicket(ctx, ticketID, userID, userRole)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh ticket: %w", err)
		}
	} else {
		// Save changes if status wasn't updated (status update saves automatically)
		if err := uc.ticketService.UpdateTicketStatus(ctx, ticketID, ticket.GetStatus(), "Ticket updated", userID); err != nil {
			return nil, fmt.Errorf("failed to save ticket changes: %w", err)
		}

		// Refresh the ticket
		ticket, err = uc.ticketService.GetTicket(ctx, ticketID, userID, userRole)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh ticket: %w", err)
		}
	}

	// Convert to response DTO
	response := dto.TicketFromEntity(ticket)
	return response, nil
}

// parseTicketStatus converts string to TicketStatus
func parseTicketStatus(status string) (entities.TicketStatus, error) {
	switch status {
	case "pending":
		return entities.StatusPending, nil
	case "waiting_approval":
		return entities.StatusWaitingApproval, nil
	case "approved":
		return entities.StatusApproved, nil
	case "rejected":
		return entities.StatusRejected, nil
	case "in_progress":
		return entities.StatusInProgress, nil
	case "completed":
		return entities.StatusCompleted, nil
	case "closed":
		return entities.StatusClosed, nil
	default:
		return "", fmt.Errorf("unknown status: %s", status)
	}
}
