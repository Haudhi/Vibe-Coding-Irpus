package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// GetTicketUseCase implements the use case for retrieving a single ticket
type GetTicketUseCase struct {
	ticketService *services.TicketService
}

// NewGetTicketUseCase creates a new GetTicketUseCase
func NewGetTicketUseCase(ticketService *services.TicketService) *GetTicketUseCase {
	return &GetTicketUseCase{
		ticketService: ticketService,
	}
}

// Execute retrieves a ticket by ID with access control
func (uc *GetTicketUseCase) Execute(ctx context.Context, req *dto.GetTicketRequest) (*dto.TicketResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Convert role string to enum
	userRole, err := parseUserRole(req.UserRole)
	if err != nil {
		return nil, fmt.Errorf("invalid user role: %w", err)
	}

	// Get ticket
	ticket, err := uc.ticketService.GetTicket(ctx, req.TicketID, req.UserID, userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Convert to response DTO
	response := dto.TicketFromEntity(ticket)
	return response, nil
}

// parseUserRole converts string to UserRole
func parseUserRole(role string) (entities.UserRole, error) {
	switch role {
	case "requester":
		return entities.RoleRequester, nil
	case "approver":
		return entities.RoleApprover, nil
	case "admin":
		return entities.RoleAdmin, nil
	default:
		return "", fmt.Errorf("unknown role: %s", role)
	}
}