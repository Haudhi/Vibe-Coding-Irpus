package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// CreateAssetUseCase implements the use case for creating an asset
type CreateAssetUseCase struct {
	assetRepo services.AssetRepository
}

// NewCreateAssetUseCase creates a new CreateAssetUseCase
func NewCreateAssetUseCase(assetRepo services.AssetRepository) *CreateAssetUseCase {
	return &CreateAssetUseCase{
		assetRepo: assetRepo,
	}
}

// Execute creates a new asset
func (uc *CreateAssetUseCase) Execute(ctx context.Context, req *dto.CreateAssetRequest) (*dto.AssetResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Parse category
	category, err := entities.ValidateCategory(req.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Create asset entity
	asset, err := entities.NewAsset(
		req.Name,
		req.Description,
		category,
		req.Quantity,
		req.Location,
		valueobjects.NewMoney(req.UnitCost),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset: %w", err)
	}

	// Save asset
	if err := uc.assetRepo.Create(ctx, asset); err != nil {
		return nil, fmt.Errorf("failed to save asset: %w", err)
	}

	// Convert to response DTO
	response := dto.AssetFromEntity(asset)
	return response, nil
}
