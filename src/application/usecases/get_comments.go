package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// GetCommentsUseCase implements the use case for getting ticket comments
type GetCommentsUseCase struct {
	ticketService *services.TicketService
}

// NewGetCommentsUseCase creates a new GetCommentsUseCase
func NewGetCommentsUseCase(ticketService *services.TicketService) *GetCommentsUseCase {
	return &GetCommentsUseCase{
		ticketService: ticketService,
	}
}

// Execute retrieves comments for a ticket with pagination
func (uc *GetCommentsUseCase) Execute(ctx context.Context, req *dto.GetCommentsRequest) (*dto.GetCommentsResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Parse user role
	userRole, err := entities.RoleFromString(req.UserRole)
	if err != nil {
		return nil, fmt.Errorf("invalid user role: %w", err)
	}

	// Get the ticket to verify access
	ticket, err := uc.ticketService.GetTicket(ctx, req.TicketID, req.UserID, userRole)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Get comments from ticket
	comments := ticket.GetComments()

	// Apply pagination
	offset := (req.Page - 1) * req.Limit
	total := len(comments)

	// Ensure offset is within bounds
	if offset >= total {
		return &dto.GetCommentsResponse{
			Comments: []*dto.CommentResponse{},
			Page:     req.Page,
			Limit:    req.Limit,
			Total:    total,
		}, nil
	}

	// Calculate end index
	end := offset + req.Limit
	if end > total {
		end = total
	}

	// Get paginated comments
	paginatedComments := comments[offset:end]

	// Convert to response DTOs
	commentResponses := make([]*dto.CommentResponse, len(paginatedComments))
	for i, comment := range paginatedComments {
		commentResponses[i] = dto.CommentFromEntity(comment)
	}

	return &dto.GetCommentsResponse{
		Comments: commentResponses,
		Page:     req.Page,
		Limit:    req.Limit,
		Total:    total,
	}, nil
}
