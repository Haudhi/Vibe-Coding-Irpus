package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// GetTicketsUseCase implements the use case for retrieving tickets
type GetTicketsUseCase struct {
	ticketService *services.TicketService
}

// NewGetTicketsUseCase creates a new GetTicketsUseCase
func NewGetTicketsUseCase(ticketService *services.TicketService) *GetTicketsUseCase {
	return &GetTicketsUseCase{
		ticketService: ticketService,
	}
}

// Execute retrieves tickets based on user role and filters
func (uc *GetTicketsUseCase) Execute(ctx context.Context, req *dto.GetTicketsRequest) (*dto.GetTicketsResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	var tickets []*entities.Ticket
	var err error

	// Calculate offset from page and limit
	offset := (req.Page - 1) * req.Limit

	// Get tickets based on user role
	switch req.UserRole {
	case "admin":
		// Admin can see all tickets
		tickets, err = uc.ticketService.GetAllTickets(ctx, req.Limit, offset)
	case "requester":
		// Requester can only see their own tickets
		tickets, err = uc.ticketService.GetUserTickets(ctx, req.UserID, req.Limit, offset)
	case "approver":
		// Approver can see tickets requiring approval
		tickets, err = uc.ticketService.GetTicketsNeedingApproval(ctx, req.Limit, offset)
	default:
		return nil, fmt.Errorf("invalid user role: %s", req.UserRole)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}

	// Convert to response DTOs
	ticketDTOs := make([]*dto.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketDTOs[i] = dto.TicketFromEntity(ticket)
	}

	response := &dto.GetTicketsResponse{
		Tickets:   ticketDTOs,
		Page:      req.Page,
		Limit:     req.Limit,
		Total:     len(ticketDTOs), // In real implementation, get total from repository
	}

	return response, nil
}