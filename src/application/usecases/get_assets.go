package usecases

import (
	"context"
	"fmt"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
)

// GetAssetsUseCase implements the use case for getting assets with filtering
type GetAssetsUseCase struct {
	assetRepo services.AssetRepository
}

// NewGetAssetsUseCase creates a new GetAssetsUseCase
func NewGetAssetsUseCase(assetRepo services.AssetRepository) *GetAssetsUseCase {
	return &GetAssetsUseCase{
		assetRepo: assetRepo,
	}
}

// Execute retrieves assets with filtering and pagination
func (uc *GetAssetsUseCase) Execute(ctx context.Context, req *dto.GetAssetsRequest) (*dto.GetAssetsResponse, error) {
	// Validate input
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Calculate offset
	offset := (req.Page - 1) * req.Limit

	var assets []*entities.Asset
	var err error

	// Apply filters
	if req.Category != "" {
		// Get by category
		category, err := entities.ValidateCategory(req.Category)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
		assets, err = uc.assetRepo.GetByCategory(ctx, category, req.Limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get assets by category: %w", err)
		}
	} else if req.Condition != "" {
		// Get by condition
		condition, err := entities.ValidateCondition(req.Condition)
		if err != nil {
			return nil, fmt.Errorf("invalid condition: %w", err)
		}
		assets, err = uc.assetRepo.GetByCondition(ctx, condition, req.Limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get assets by condition: %w", err)
		}
	} else {
		// Get all
		assets, err = uc.assetRepo.GetAll(ctx, req.Limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get assets: %w", err)
		}
	}

	// Get total count
	total, err := uc.assetRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Convert to response DTOs
	assetResponses := make([]*dto.AssetResponse, len(assets))
	for i, asset := range assets {
		assetResponses[i] = dto.AssetFromEntity(asset)
	}

	return &dto.GetAssetsResponse{
		Assets: assetResponses,
		Page:   req.Page,
		Limit:  req.Limit,
		Total:  total,
	}, nil
}
