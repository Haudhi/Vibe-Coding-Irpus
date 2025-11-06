package repositories

import (
	"github.com/company/ga-ticketing/src/domain/entities"
)

// AssetRepository defines the interface for asset persistence operations
type AssetRepository interface {
	// Create saves a new asset to the database
	Create(asset *entities.Asset) error

	// FindByID retrieves an asset by their ID
	FindByID(id string) (*entities.Asset, error)

	// FindByAssetCode retrieves an asset by their asset code
	FindByAssetCode(assetCode string) (*entities.Asset, error)

	// FindByCategory retrieves all assets in a specific category
	FindByCategory(category entities.AssetCategory) ([]*entities.Asset, error)

	// FindByLocation retrieves all assets at a specific location
	FindByLocation(location string) ([]*entities.Asset, error)

	// FindByCondition retrieves all assets with a specific condition
	FindByCondition(condition entities.AssetCondition) ([]*entities.Asset, error)

	// FindAvailable retrieves all assets that are available (quantity > 0 and condition = good)
	FindAvailable() ([]*entities.Asset, error)

	// FindRequiringMaintenance retrieves all assets that need maintenance
	FindRequiringMaintenance() ([]*entities.Asset, error)

	// FindAll retrieves all assets with optional filtering
	FindAll(filters AssetFilters) ([]*entities.Asset, error)

	// Update updates an existing asset in the database
	Update(asset *entities.Asset) error

	// Delete removes an asset from the database
	Delete(id string) error

	// CheckAvailability checks if an asset has sufficient available quantity
	CheckAvailability(assetID string, quantity int) (bool, error)
}

// AssetFilters defines filters for asset queries
type AssetFilters struct {
	Category      *entities.AssetCategory
	Condition     *entities.AssetCondition
	Location      *string
	AvailableOnly *bool
	Limit         *int
	Offset        *int
}