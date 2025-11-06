package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// UpdateInventoryUseCase implements the use case for updating asset inventory
type UpdateInventoryUseCase struct {
	assetRepo services.AssetRepository
}

// NewUpdateInventoryUseCase creates a new UpdateInventoryUseCase
func NewUpdateInventoryUseCase(assetRepo services.AssetRepository) *UpdateInventoryUseCase {
	return &UpdateInventoryUseCase{
		assetRepo: assetRepo,
	}
}

// Execute updates asset inventory
func (uc *UpdateInventoryUseCase) Execute(ctx context.Context, assetID, userID string, req *dto.UpdateInventoryRequest) (*dto.AssetResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	if assetID == "" {
		return nil, fmt.Errorf("asset ID is required")
	}

	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Get existing asset
	asset, err := uc.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Parse change type
	var changeType entities.ChangeType
	switch req.ChangeType {
	case "add":
		changeType = entities.ChangeTypeAdd
	case "remove":
		changeType = entities.ChangeTypeRemove
	case "adjust":
		changeType = entities.ChangeTypeAdjust
	default:
		return nil, fmt.Errorf("invalid change type: %s", req.ChangeType)
	}

	// Update inventory
	if err := asset.UpdateInventory(changeType, req.Quantity, req.Reason, userID); err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	// Save updated asset
	if err := uc.assetRepo.Update(ctx, asset); err != nil {
		return nil, fmt.Errorf("failed to save asset: %w", err)
	}

	// Save inventory logs
	logs := asset.GetInventoryLogs()
	if len(logs) > 0 {
		// Save the most recent log
		if err := uc.assetRepo.AddInventoryLog(ctx, logs[len(logs)-1]); err != nil {
			return nil, fmt.Errorf("failed to save inventory log: %w", err)
		}
	}

	// Convert to response DTO
	response := dto.AssetFromEntity(asset)
	return response, nil
}
