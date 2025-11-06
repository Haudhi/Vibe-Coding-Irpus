package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// UpdateAssetUseCase implements the use case for updating an asset
type UpdateAssetUseCase struct {
	assetRepo services.AssetRepository
}

// NewUpdateAssetUseCase creates a new UpdateAssetUseCase
func NewUpdateAssetUseCase(assetRepo services.AssetRepository) *UpdateAssetUseCase {
	return &UpdateAssetUseCase{
		assetRepo: assetRepo,
	}
}

// Execute updates an asset
func (uc *UpdateAssetUseCase) Execute(ctx context.Context, assetID string, req *dto.UpdateAssetRequest) (*dto.AssetResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	if assetID == "" {
		return nil, fmt.Errorf("asset ID is required")
	}

	// Get existing asset
	asset, err := uc.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		if err := asset.SetName(*req.Name); err != nil {
			return nil, fmt.Errorf("failed to update name: %w", err)
		}
	}

	if req.Description != nil {
		asset.SetDescription(*req.Description)
	}

	if req.Location != nil {
		if err := asset.SetLocation(*req.Location); err != nil {
			return nil, fmt.Errorf("failed to update location: %w", err)
		}
	}

	if req.Condition != nil {
		condition, err := entities.ValidateCondition(*req.Condition)
		if err != nil {
			return nil, fmt.Errorf("invalid condition: %w", err)
		}
		asset.SetCondition(condition)
	}

	if req.UnitCost != nil {
		if err := asset.SetUnitCost(valueobjects.NewMoney(*req.UnitCost)); err != nil {
			return nil, fmt.Errorf("failed to update unit cost: %w", err)
		}
	}

	// Save updated asset
	if err := uc.assetRepo.Update(ctx, asset); err != nil {
		return nil, fmt.Errorf("failed to update asset: %w", err)
	}

	// Convert to response DTO
	response := dto.AssetFromEntity(asset)
	return response, nil
}
