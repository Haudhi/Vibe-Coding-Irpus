package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// AssetCategory represents the category of an asset
type AssetCategory string

const (
	AssetCategoryOfficeFurniture      AssetCategory = "office_furniture"
	AssetCategoryOfficeSupplies       AssetCategory = "office_supplies"
	AssetCategoryPantrySupplies       AssetCategory = "pantry_supplies"
	AssetCategoryFacilityEquipment    AssetCategory = "facility_equipment"
	AssetCategoryMeetingRoomEquipment AssetCategory = "meeting_room_equipment"
	AssetCategoryCleaningSupplies     AssetCategory = "cleaning_supplies"
)

// AssetCondition represents the condition of an asset
type AssetCondition string

const (
	ConditionGood              AssetCondition = "good"
	ConditionNeedsMaintenance  AssetCondition = "needs_maintenance"
	ConditionBroken            AssetCondition = "broken"
)

// ChangeType represents the type of inventory change
type ChangeType string

const (
	ChangeTypeAdd    ChangeType = "add"
	ChangeTypeRemove ChangeType = "remove"
	ChangeTypeAdjust ChangeType = "adjust"
)

// InventoryLog represents a change in asset inventory
type InventoryLog struct {
	ID         string
	AssetID    string
	ChangeType ChangeType
	Quantity   int
	Reason     string
	CreatedBy  string
	CreatedAt  time.Time
}

// NewInventoryLog creates a new inventory log entry
func NewInventoryLog(assetID string, changeType ChangeType, quantity int, reason, createdBy string) *InventoryLog {
	return &InventoryLog{
		ID:         uuid.New().String(),
		AssetID:    assetID,
		ChangeType: changeType,
		Quantity:   quantity,
		Reason:     reason,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
	}
}

// Asset represents a physical inventory item
type Asset struct {
	id                string
	assetCode         string
	name              string
	description       string
	category          AssetCategory
	quantity          int
	availableQuantity int
	location          string
	condition         AssetCondition
	unitCost          *valueobjects.Money
	lastMaintenanceAt *time.Time
	nextMaintenanceAt *time.Time
	createdAt         time.Time
	updatedAt         time.Time
	inventoryLogs     []*InventoryLog
}

// NewAsset creates a new asset
func NewAsset(
	name, description string,
	category AssetCategory,
	quantity int,
	location string,
	unitCost *valueobjects.Money,
) (*Asset, error) {
	// Validate input
	if name == "" {
		return nil, errors.New("asset name is required")
	}
	if quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}
	if unitCost == nil {
		return nil, errors.New("unit cost is required")
	}
	if unitCost.Amount < 0 {
		return nil, errors.New("unit cost cannot be negative")
	}
	if location == "" {
		return nil, errors.New("location is required")
	}

	now := time.Now()
	asset := &Asset{
		id:                uuid.New().String(),
		assetCode:         generateAssetCode(category, now),
		name:              name,
		description:       description,
		category:          category,
		quantity:          quantity,
		availableQuantity: quantity,
		location:          location,
		condition:         ConditionGood,
		unitCost:          unitCost,
		createdAt:         now,
		updatedAt:         now,
		inventoryLogs:     make([]*InventoryLog, 0),
	}

	return asset, nil
}

// Getters
func (a *Asset) GetID() string                         { return a.id }
func (a *Asset) GetAssetCode() string                  { return a.assetCode }
func (a *Asset) GetName() string                       { return a.name }
func (a *Asset) GetDescription() string                { return a.description }
func (a *Asset) GetCategory() AssetCategory            { return a.category }
func (a *Asset) GetQuantity() int                      { return a.quantity }
func (a *Asset) GetAvailableQuantity() int             { return a.availableQuantity }
func (a *Asset) GetLocation() string                   { return a.location }
func (a *Asset) GetCondition() AssetCondition          { return a.condition }
func (a *Asset) GetUnitCost() *valueobjects.Money      { return a.unitCost }
func (a *Asset) GetLastMaintenanceAt() *time.Time      { return a.lastMaintenanceAt }
func (a *Asset) GetNextMaintenanceAt() *time.Time      { return a.nextMaintenanceAt }
func (a *Asset) GetCreatedAt() time.Time               { return a.createdAt }
func (a *Asset) GetUpdatedAt() time.Time               { return a.updatedAt }
func (a *Asset) GetInventoryLogs() []*InventoryLog     { return a.inventoryLogs }

// Setters (business logic)
func (a *Asset) SetName(name string) error {
	if name == "" {
		return errors.New("asset name is required")
	}
	a.name = name
	a.updatedAt = time.Now()
	return nil
}

func (a *Asset) SetDescription(description string) {
	a.description = description
	a.updatedAt = time.Now()
}

func (a *Asset) SetLocation(location string) error {
	if location == "" {
		return errors.New("location is required")
	}
	a.location = location
	a.updatedAt = time.Now()
	return nil
}

func (a *Asset) SetCondition(condition AssetCondition) {
	a.condition = condition
	a.updatedAt = time.Now()
}

func (a *Asset) SetUnitCost(cost *valueobjects.Money) error {
	if cost == nil {
		return errors.New("unit cost is required")
	}
	if cost.Amount < 0 {
		return errors.New("unit cost cannot be negative")
	}
	a.unitCost = cost
	a.updatedAt = time.Now()
	return nil
}

func (a *Asset) SetMaintenanceDates(lastMaintenance, nextMaintenance *time.Time) {
	a.lastMaintenanceAt = lastMaintenance
	a.nextMaintenanceAt = nextMaintenance
	a.updatedAt = time.Now()
}

// UpdateInventory updates the asset inventory
func (a *Asset) UpdateInventory(changeType ChangeType, quantity int, reason, createdBy string) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if reason == "" {
		return errors.New("reason is required")
	}
	if createdBy == "" {
		return errors.New("created by is required")
	}

	switch changeType {
	case ChangeTypeAdd:
		a.quantity += quantity
		a.availableQuantity += quantity
	case ChangeTypeRemove:
		if a.availableQuantity < quantity {
			return fmt.Errorf("insufficient available quantity (available: %d, requested: %d)", a.availableQuantity, quantity)
		}
		a.quantity -= quantity
		a.availableQuantity -= quantity
	case ChangeTypeAdjust:
		// Adjust allows setting to a specific quantity
		if quantity < 0 {
			return errors.New("adjusted quantity cannot be negative")
		}
		diff := quantity - a.quantity
		a.quantity = quantity
		a.availableQuantity += diff
		if a.availableQuantity < 0 {
			a.availableQuantity = 0
		}
	default:
		return errors.New("invalid change type")
	}

	// Create inventory log
	log := NewInventoryLog(a.id, changeType, quantity, reason, createdBy)
	a.inventoryLogs = append(a.inventoryLogs, log)

	a.updatedAt = time.Now()
	return nil
}

// AllocateQuantity allocates a quantity for use (reduces available)
func (a *Asset) AllocateQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	if a.availableQuantity < quantity {
		return fmt.Errorf("insufficient available quantity (available: %d, requested: %d)", a.availableQuantity, quantity)
	}

	a.availableQuantity -= quantity
	a.updatedAt = time.Now()
	return nil
}

// ReleaseQuantity releases a previously allocated quantity
func (a *Asset) ReleaseQuantity(quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	newAvailable := a.availableQuantity + quantity
	if newAvailable > a.quantity {
		return errors.New("released quantity would exceed total quantity")
	}

	a.availableQuantity = newAvailable
	a.updatedAt = time.Now()
	return nil
}

// IsAvailable checks if the asset has available quantity
func (a *Asset) IsAvailable() bool {
	return a.availableQuantity > 0 && a.condition == ConditionGood
}

// RequiresMaintenance checks if the asset needs maintenance
func (a *Asset) RequiresMaintenance() bool {
	if a.condition == ConditionNeedsMaintenance || a.condition == ConditionBroken {
		return true
	}

	if a.nextMaintenanceAt != nil && time.Now().After(*a.nextMaintenanceAt) {
		return true
	}

	return false
}

// generateAssetCode generates a unique asset code
func generateAssetCode(category AssetCategory, createdAt time.Time) string {
	// Generate a code based on category and timestamp
	// Format: CATEGORY-YYYYMMDD-XXXX (where XXXX is a sequence or random)
	year := createdAt.Year()
	month := int(createdAt.Month())
	day := createdAt.Day()

	var categoryCode string
	switch category {
	case AssetCategoryOfficeFurniture:
		categoryCode = "OF"
	case AssetCategoryOfficeSupplies:
		categoryCode = "OS"
	case AssetCategoryPantrySupplies:
		categoryCode = "PS"
	case AssetCategoryFacilityEquipment:
		categoryCode = "FE"
	case AssetCategoryMeetingRoomEquipment:
		categoryCode = "MR"
	case AssetCategoryCleaningSupplies:
		categoryCode = "CS"
	default:
		categoryCode = "GA"
	}

	// Use a simple sequence based on time
	sequence := (createdAt.Hour()*3600 + createdAt.Minute()*60 + createdAt.Second()) % 10000

	return fmt.Sprintf("%s-%04d%02d%02d-%04d", categoryCode, year, month, day, sequence)
}

// ValidateCategory validates if the category is valid
func ValidateCategory(category string) (AssetCategory, error) {
	switch category {
	case string(AssetCategoryOfficeFurniture):
		return AssetCategoryOfficeFurniture, nil
	case string(AssetCategoryOfficeSupplies):
		return AssetCategoryOfficeSupplies, nil
	case string(AssetCategoryPantrySupplies):
		return AssetCategoryPantrySupplies, nil
	case string(AssetCategoryFacilityEquipment):
		return AssetCategoryFacilityEquipment, nil
	case string(AssetCategoryMeetingRoomEquipment):
		return AssetCategoryMeetingRoomEquipment, nil
	case string(AssetCategoryCleaningSupplies):
		return AssetCategoryCleaningSupplies, nil
	default:
		return "", fmt.Errorf("invalid asset category: %s", category)
	}
}

// ValidateCondition validates if the condition is valid
func ValidateCondition(condition string) (AssetCondition, error) {
	switch condition {
	case string(ConditionGood):
		return ConditionGood, nil
	case string(ConditionNeedsMaintenance):
		return ConditionNeedsMaintenance, nil
	case string(ConditionBroken):
		return ConditionBroken, nil
	default:
		return "", fmt.Errorf("invalid asset condition: %s", condition)
	}
}
