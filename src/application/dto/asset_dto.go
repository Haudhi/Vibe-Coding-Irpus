package dto

import (
	"fmt"

	"github.com/company/ga-ticketing/src/domain/entities"
)

// CreateAssetRequest represents a request to create an asset
type CreateAssetRequest struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
	Category    string `json:"category" validate:"required,oneof=office_furniture office_supplies pantry_supplies facility_equipment meeting_room_equipment cleaning_supplies"`
	Quantity    int    `json:"quantity" validate:"required,min=0"`
	Location    string `json:"location" validate:"required,max=255"`
	UnitCost    int64  `json:"unit_cost" validate:"required,min=0"`
}

// Validate validates the CreateAssetRequest
func (req *CreateAssetRequest) Validate() error {
	if req.Name == "" {
		return fmt.Errorf("asset name is required")
	}
	if len(req.Name) > 255 {
		return fmt.Errorf("asset name must be 255 characters or less")
	}
	if req.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}
	if req.Location == "" {
		return fmt.Errorf("location is required")
	}
	if len(req.Location) > 255 {
		return fmt.Errorf("location must be 255 characters or less")
	}
	if req.UnitCost < 0 {
		return fmt.Errorf("unit cost cannot be negative")
	}
	return nil
}

// UpdateAssetRequest represents a request to update an asset
type UpdateAssetRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
	Location    *string `json:"location,omitempty" validate:"omitempty,max=255"`
	Condition   *string `json:"condition,omitempty" validate:"omitempty,oneof=good needs_maintenance broken"`
	UnitCost    *int64  `json:"unit_cost,omitempty" validate:"omitempty,min=0"`
}

// Validate validates the UpdateAssetRequest
func (req *UpdateAssetRequest) Validate() error {
	if req.Name != nil && *req.Name == "" {
		return fmt.Errorf("asset name cannot be empty")
	}
	if req.Name != nil && len(*req.Name) > 255 {
		return fmt.Errorf("asset name must be 255 characters or less")
	}
	if req.Location != nil && *req.Location == "" {
		return fmt.Errorf("location cannot be empty")
	}
	if req.Location != nil && len(*req.Location) > 255 {
		return fmt.Errorf("location must be 255 characters or less")
	}
	if req.UnitCost != nil && *req.UnitCost < 0 {
		return fmt.Errorf("unit cost cannot be negative")
	}
	return nil
}

// UpdateInventoryRequest represents a request to update asset inventory
type UpdateInventoryRequest struct {
	ChangeType string `json:"change_type" validate:"required,oneof=add remove adjust"`
	Quantity   int    `json:"quantity" validate:"required,min=1"`
	Reason     string `json:"reason" validate:"required,max=500"`
}

// Validate validates the UpdateInventoryRequest
func (req *UpdateInventoryRequest) Validate() error {
	if req.ChangeType == "" {
		return fmt.Errorf("change type is required")
	}
	if req.ChangeType != "add" && req.ChangeType != "remove" && req.ChangeType != "adjust" {
		return fmt.Errorf("change type must be one of: add, remove, adjust")
	}
	if req.Quantity < 1 {
		return fmt.Errorf("quantity must be at least 1")
	}
	if req.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	if len(req.Reason) > 500 {
		return fmt.Errorf("reason must be 500 characters or less")
	}
	return nil
}

// GetAssetsRequest represents a request to get assets with filtering
type GetAssetsRequest struct {
	Page      int    `json:"page" validate:"min=1"`
	Limit     int    `json:"limit" validate:"min=1,max=100"`
	Category  string `json:"category,omitempty"`
	Condition string `json:"condition,omitempty"`
}

// Validate validates the GetAssetsRequest
func (req *GetAssetsRequest) Validate() error {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	return nil
}

// AssetResponse represents an asset response
type AssetResponse struct {
	ID                string  `json:"id"`
	AssetCode         string  `json:"asset_code"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Category          string  `json:"category"`
	Quantity          int     `json:"quantity"`
	AvailableQuantity int     `json:"available_quantity"`
	Location          string  `json:"location"`
	Condition         string  `json:"condition"`
	UnitCost          int64   `json:"unit_cost"`
	LastMaintenanceAt *string `json:"last_maintenance_at"`
	NextMaintenanceAt *string `json:"next_maintenance_at"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// GetAssetsResponse represents a response for getting multiple assets
type GetAssetsResponse struct {
	Assets []*AssetResponse `json:"assets"`
	Page   int              `json:"page"`
	Limit  int              `json:"limit"`
	Total  int              `json:"total"`
}

// InventoryLogResponse represents an inventory log response
type InventoryLogResponse struct {
	ID         string `json:"id"`
	AssetID    string `json:"asset_id"`
	ChangeType string `json:"change_type"`
	Quantity   int    `json:"quantity"`
	Reason     string `json:"reason"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
}

// AssetFromEntity converts a domain entity to DTO
func AssetFromEntity(asset *entities.Asset) *AssetResponse {
	response := &AssetResponse{
		ID:                asset.GetID(),
		AssetCode:         asset.GetAssetCode(),
		Name:              asset.GetName(),
		Description:       asset.GetDescription(),
		Category:          string(asset.GetCategory()),
		Quantity:          asset.GetQuantity(),
		AvailableQuantity: asset.GetAvailableQuantity(),
		Location:          asset.GetLocation(),
		Condition:         string(asset.GetCondition()),
		UnitCost:          asset.GetUnitCost().Amount,
		CreatedAt:         asset.GetCreatedAt().Format("2006-01-02T15:04:05Z"),
		UpdatedAt:         asset.GetUpdatedAt().Format("2006-01-02T15:04:05Z"),
	}

	// Handle optional fields
	if asset.GetLastMaintenanceAt() != nil {
		lastMaintenance := asset.GetLastMaintenanceAt().Format("2006-01-02T15:04:05Z")
		response.LastMaintenanceAt = &lastMaintenance
	}

	if asset.GetNextMaintenanceAt() != nil {
		nextMaintenance := asset.GetNextMaintenanceAt().Format("2006-01-02T15:04:05Z")
		response.NextMaintenanceAt = &nextMaintenance
	}

	return response
}

// InventoryLogFromEntity converts a domain inventory log to DTO
func InventoryLogFromEntity(log *entities.InventoryLog) *InventoryLogResponse {
	return &InventoryLogResponse{
		ID:         log.ID,
		AssetID:    log.AssetID,
		ChangeType: string(log.ChangeType),
		Quantity:   log.Quantity,
		Reason:     log.Reason,
		CreatedBy:  log.CreatedBy,
		CreatedAt:  log.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// ApproveTicketRequest represents a request to approve a ticket
type ApproveTicketRequest struct {
	Comments string `json:"comments" validate:"max=500"`
}

// Validate validates the ApproveTicketRequest
func (req *ApproveTicketRequest) Validate() error {
	if len(req.Comments) > 500 {
		return fmt.Errorf("comments must be 500 characters or less")
	}
	return nil
}

// RejectTicketRequest represents a request to reject a ticket
type RejectTicketRequest struct {
	Reason string `json:"reason" validate:"required,max=500"`
}

// Validate validates the RejectTicketRequest
func (req *RejectTicketRequest) Validate() error {
	if req.Reason == "" {
		return fmt.Errorf("rejection reason is required")
	}
	if len(req.Reason) > 500 {
		return fmt.Errorf("reason must be 500 characters or less")
	}
	return nil
}

// GetCommentsRequest represents a request to get comments with pagination
type GetCommentsRequest struct {
	TicketID string `json:"ticket_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	UserRole string `json:"user_role" validate:"required"`
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
}

// Validate validates the GetCommentsRequest
func (req *GetCommentsRequest) Validate() error {
	if req.TicketID == "" {
		return fmt.Errorf("ticket ID is required")
	}
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if req.UserRole == "" {
		return fmt.Errorf("user role is required")
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	return nil
}

// GetCommentsResponse represents a response for getting comments
type GetCommentsResponse struct {
	Comments []*CommentResponse `json:"comments"`
	Page     int                `json:"page"`
	Limit    int                `json:"limit"`
	Total    int                `json:"total"`
}
