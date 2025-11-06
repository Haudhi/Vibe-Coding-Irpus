package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// AddCommentUseCase implements the use case for adding a comment to a ticket
type AddCommentUseCase struct {
	ticketService *services.TicketService
}

// NewAddCommentUseCase creates a new AddCommentUseCase
func NewAddCommentUseCase(ticketService *services.TicketService) *AddCommentUseCase {
	return &AddCommentUseCase{
		ticketService: ticketService,
	}
}

// Execute adds a comment to a ticket
func (uc *AddCommentUseCase) Execute(ctx context.Context, ticketID, userID string, userRole entities.UserRole, req *dto.CommentRequest) (*dto.CommentResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Verify user has access to the ticket
	_, err := uc.ticketService.GetTicket(ctx, ticketID, userID, userRole)
	if err != nil {
		return nil, fmt.Errorf("access denied or ticket not found: %w", err)
	}

	// Add comment to ticket
	comment, err := uc.ticketService.AddCommentToTicket(ctx, ticketID, req.Content, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	// Convert to response DTO
	response := dto.CommentFromEntity(comment)
	return response, nil
}
