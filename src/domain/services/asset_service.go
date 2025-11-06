package services

import (
	"context"

	"github.com/company/ga-ticketing/src/domain/entities"
)

// AssetRepository defines the interface for asset persistence
type AssetRepository interface {
	Create(ctx context.Context, asset *entities.Asset) error
	GetByID(ctx context.Context, id string) (*entities.Asset, error)
	GetByAssetCode(ctx context.Context, assetCode string) (*entities.Asset, error)
	GetAll(ctx context.Context, limit, offset int) ([]*entities.Asset, error)
	GetByCategory(ctx context.Context, category entities.AssetCategory, limit, offset int) ([]*entities.Asset, error)
	GetByCondition(ctx context.Context, condition entities.AssetCondition, limit, offset int) ([]*entities.Asset, error)
	Update(ctx context.Context, asset *entities.Asset) error
	Delete(ctx context.Context, id string) error
	GetTotalCount(ctx context.Context) (int, error)
	AddInventoryLog(ctx context.Context, log *entities.InventoryLog) error
	GetInventoryLogs(ctx context.Context, assetID string, limit, offset int) ([]*entities.InventoryLog, error)
}

// AssetService provides business logic for asset operations
type AssetService struct {
	assetRepo AssetRepository
	userRepo  UserRepository
}

// NewAssetService creates a new AssetService
func NewAssetService(assetRepo AssetRepository, userRepo UserRepository) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
		userRepo:  userRepo,
	}
}
