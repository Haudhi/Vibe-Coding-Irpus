package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// CreateTicketUseCase implements the use case for creating a ticket
type CreateTicketUseCase struct {
	ticketService *services.TicketService
}

// NewCreateTicketUseCase creates a new CreateTicketUseCase
func NewCreateTicketUseCase(ticketService *services.TicketService) *CreateTicketUseCase {
	return &CreateTicketUseCase{
		ticketService: ticketService,
	}
}

// Execute creates a new ticket
func (uc *CreateTicketUseCase) Execute(ctx context.Context, req *dto.CreateTicketRequest) (*dto.TicketResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Convert category string to enum
	category, err := parseTicketCategory(req.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Convert priority string to enum
	priority, err := parseTicketPriority(req.Priority)
	if err != nil {
		return nil, fmt.Errorf("invalid priority: %w", err)
	}

	// Create ticket
	ticket, err := uc.ticketService.CreateTicket(
		ctx,
		req.Title,
		req.Description,
		category,
		priority,
		valueobjects.NewMoney(req.EstimatedCost),
		req.RequesterID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Convert to response DTO
	response := dto.TicketFromEntity(ticket)
	return response, nil
}

// parseTicketCategory converts string to TicketCategory
func parseTicketCategory(category string) (entities.TicketCategory, error) {
	switch category {
	case "office_supplies":
		return entities.CategoryOfficeSupplies, nil
	case "facility_maintenance":
		return entities.CategoryFacilityMaintenance, nil
	case "pantry_supplies":
		return entities.CategoryPantrySupplies, nil
	case "meeting_room":
		return entities.CategoryMeetingRoom, nil
	case "office_furniture":
		return entities.CategoryOfficeFurniture, nil
	case "general_service":
		return entities.CategoryGeneralService, nil
	default:
		return "", fmt.Errorf("unknown category: %s", category)
	}
}

// parseTicketPriority converts string to TicketPriority
func parseTicketPriority(priority string) (entities.TicketPriority, error) {
	switch priority {
	case "low":
		return entities.PriorityLow, nil
	case "medium":
		return entities.PriorityMedium, nil
	case "high":
		return entities.PriorityHigh, nil
	default:
		return "", fmt.Errorf("unknown priority: %s", priority)
	}
}