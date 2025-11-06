package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/services"
)

// GetAssetUseCase implements the use case for getting an asset by ID
type GetAssetUseCase struct {
	assetRepo services.AssetRepository
}

// NewGetAssetUseCase creates a new GetAssetUseCase
func NewGetAssetUseCase(assetRepo services.AssetRepository) *GetAssetUseCase {
	return &GetAssetUseCase{
		assetRepo: assetRepo,
	}
}

// Execute retrieves an asset by ID
func (uc *GetAssetUseCase) Execute(ctx context.Context, assetID string) (*dto.AssetResponse, error) {
	if assetID == "" {
		return nil, fmt.Errorf("asset ID is required")
	}

	// Get asset
	asset, err := uc.assetRepo.GetByID(ctx, assetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Convert to response DTO
	response := dto.AssetFromEntity(asset)
	return response, nil
}
