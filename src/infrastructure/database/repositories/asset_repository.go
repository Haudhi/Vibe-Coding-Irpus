package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/company/ga-ticketing/src/domain/entities"
	"github.com/company/ga-ticketing/src/domain/services"
	"github.com/company/ga-ticketing/src/domain/valueobjects"
)

// assetRecord represents the database model for an asset
type assetRecord struct {
	ID                string
	AssetCode         string
	Name              string
	Description       string
	Category          string
	Quantity          int
	AvailableQuantity int
	Location          string
	Condition         string
	UnitCost          int64
	LastMaintenanceAt sql.NullTime
	NextMaintenanceAt sql.NullTime
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// AssetRepository implements the AssetRepository interface using PostgreSQL
type AssetRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

// NewAssetRepository creates a new AssetRepository
func NewAssetRepository(pool *pgxpool.Pool, logger *zap.Logger) services.AssetRepository {
	return &AssetRepository{
		pool:   pool,
		logger: logger,
	}
}

// Create saves a new asset to the database
func (r *AssetRepository) Create(ctx context.Context, asset *entities.Asset) error {
	query := `
		INSERT INTO assets (
			id, asset_code, name, description, category, quantity,
			available_quantity, location, condition, unit_cost,
			last_maintenance_at, next_maintenance_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	var lastMaintenance, nextMaintenance sql.NullTime
	if asset.GetLastMaintenanceAt() != nil {
		lastMaintenance = sql.NullTime{Time: *asset.GetLastMaintenanceAt(), Valid: true}
	}
	if asset.GetNextMaintenanceAt() != nil {
		nextMaintenance = sql.NullTime{Time: *asset.GetNextMaintenanceAt(), Valid: true}
	}

	_, err := r.pool.Exec(ctx, query,
		asset.GetID(),
		asset.GetAssetCode(),
		asset.GetName(),
		asset.GetDescription(),
		string(asset.GetCategory()),
		asset.GetQuantity(),
		asset.GetAvailableQuantity(),
		asset.GetLocation(),
		string(asset.GetCondition()),
		asset.GetUnitCost().Amount,
		lastMaintenance,
		nextMaintenance,
		asset.GetCreatedAt(),
		asset.GetUpdatedAt(),
	)

	if err != nil {
		r.logger.Error("Failed to create asset",
			zap.String("asset_id", asset.GetID()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create asset: %w", err)
	}

	r.logger.Info("Asset created successfully",
		zap.String("asset_id", asset.GetID()),
		zap.String("asset_code", asset.GetAssetCode()),
	)

	return nil
}

// GetByID retrieves an asset by ID
func (r *AssetRepository) GetByID(ctx context.Context, id string) (*entities.Asset, error) {
	query := `
		SELECT
			id, asset_code, name, description, category, quantity,
			available_quantity, location, condition, unit_cost,
			last_maintenance_at, next_maintenance_at, created_at, updated_at
		FROM assets
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var asset assetRecord

	err := row.Scan(
		&asset.ID,
		&asset.AssetCode,
		&asset.Name,
		&asset.Description,
		&asset.Category,
		&asset.Quantity,
		&asset.AvailableQuantity,
		&asset.Location,
		&asset.Condition,
		&asset.UnitCost,
		&asset.LastMaintenanceAt,
		&asset.NextMaintenanceAt,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		r.logger.Error("Failed to get asset by ID",
			zap.String("asset_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Convert to domain entity
	domainAsset, err := r.mapToDomainEntity(&asset)
	if err != nil {
		return nil, fmt.Errorf("failed to map asset to domain entity: %w", err)
	}

	return domainAsset, nil
}

// GetByAssetCode retrieves an asset by asset code
func (r *AssetRepository) GetByAssetCode(ctx context.Context, assetCode string) (*entities.Asset, error) {
	query := `
		SELECT
			id, asset_code, name, description, category, quantity,
			available_quantity, location, condition, unit_cost,
			last_maintenance_at, next_maintenance_at, created_at, updated_at
		FROM assets
		WHERE asset_code = $1
	`

	row := r.pool.QueryRow(ctx, query, assetCode)

	var asset assetRecord

	err := row.Scan(
		&asset.ID,
		&asset.AssetCode,
		&asset.Name,
		&asset.Description,
		&asset.Category,
		&asset.Quantity,
		&asset.AvailableQuantity,
		&asset.Location,
		&asset.Condition,
		&asset.UnitCost,
		&asset.LastMaintenanceAt,
		&asset.NextMaintenanceAt,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("asset not found")
		}
		r.logger.Error("Failed to get asset by code",
			zap.String("asset_code", assetCode),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}

	// Convert to domain entity
	domainAsset, err := r.mapToDomainEntity(&asset)
	if err != nil {
		return nil, fmt.Errorf("failed to map asset to domain entity: %w", err)
	}

	return domainAsset, nil
}

// GetAll retrieves all assets with pagination
func (r *AssetRepository) GetAll(ctx context.Context, limit, offset int) ([]*entities.Asset, error) {
	query := `
		SELECT
			id, asset_code, name, description, category, quantity,
			available_quantity, location, condition, unit_cost,
			last_maintenance_at, next_maintenance_at, created_at, updated_at
		FROM assets
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get all assets", zap.Error(err))
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	defer rows.Close()

	var assets []*entities.Asset
	for rows.Next() {
		var asset assetRecord

		if err := rows.Scan(
			&asset.ID,
			&asset.AssetCode,
			&asset.Name,
			&asset.Description,
			&asset.Category,
			&asset.Quantity,
			&asset.AvailableQuantity,
			&asset.Location,
			&asset.Condition,
			&asset.UnitCost,
			&asset.LastMaintenanceAt,
			&asset.NextMaintenanceAt,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan asset row: %w", err)
		}

		domainAsset, err := r.mapToDomainEntity(&asset)
		if err != nil {
			r.logger.Warn("Failed to map asset",
				zap.String("asset_id", asset.ID),
				zap.Error(err),
			)
			continue
		}
		assets = append(assets, domainAsset)
	}

	return assets, nil
}

// GetByCategory retrieves assets by category with pagination
func (r *AssetRepository) GetByCategory(ctx context.Context, category entities.AssetCategory, limit, offset int) ([]*entities.Asset, error) {
	query := `
		SELECT
			id, asset_code, name, description, category, quantity,
			available_quantity, location, condition, unit_cost,
			last_maintenance_at, next_maintenance_at, created_at, updated_at
		FROM assets
		WHERE category = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, string(category), limit, offset)
	if err != nil {
		r.logger.Error("Failed to get assets by category",
			zap.String("category", string(category)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	defer rows.Close()

	var assets []*entities.Asset
	for rows.Next() {
		var asset assetRecord

		if err := rows.Scan(
			&asset.ID,
			&asset.AssetCode,
			&asset.Name,
			&asset.Description,
			&asset.Category,
			&asset.Quantity,
			&asset.AvailableQuantity,
			&asset.Location,
			&asset.Condition,
			&asset.UnitCost,
			&asset.LastMaintenanceAt,
			&asset.NextMaintenanceAt,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan asset row: %w", err)
		}

		domainAsset, err := r.mapToDomainEntity(&asset)
		if err != nil {
			r.logger.Warn("Failed to map asset",
				zap.String("asset_id", asset.ID),
				zap.Error(err),
			)
			continue
		}
		assets = append(assets, domainAsset)
	}

	return assets, nil
}

// GetByCondition retrieves assets by condition with pagination
func (r *AssetRepository) GetByCondition(ctx context.Context, condition entities.AssetCondition, limit, offset int) ([]*entities.Asset, error) {
	query := `
		SELECT
			id, asset_code, name, description, category, quantity,
			available_quantity, location, condition, unit_cost,
			last_maintenance_at, next_maintenance_at, created_at, updated_at
		FROM assets
		WHERE condition = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, string(condition), limit, offset)
	if err != nil {
		r.logger.Error("Failed to get assets by condition",
			zap.String("condition", string(condition)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	defer rows.Close()

	var assets []*entities.Asset
	for rows.Next() {
		var asset assetRecord

		if err := rows.Scan(
			&asset.ID,
			&asset.AssetCode,
			&asset.Name,
			&asset.Description,
			&asset.Category,
			&asset.Quantity,
			&asset.AvailableQuantity,
			&asset.Location,
			&asset.Condition,
			&asset.UnitCost,
			&asset.LastMaintenanceAt,
			&asset.NextMaintenanceAt,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan asset row: %w", err)
		}

		domainAsset, err := r.mapToDomainEntity(&asset)
		if err != nil {
			r.logger.Warn("Failed to map asset",
				zap.String("asset_id", asset.ID),
				zap.Error(err),
			)
			continue
		}
		assets = append(assets, domainAsset)
	}

	return assets, nil
}

// Update updates an existing asset
func (r *AssetRepository) Update(ctx context.Context, asset *entities.Asset) error {
	query := `
		UPDATE assets SET
			name = $2,
			description = $3,
			quantity = $4,
			available_quantity = $5,
			location = $6,
			condition = $7,
			unit_cost = $8,
			last_maintenance_at = $9,
			next_maintenance_at = $10,
			updated_at = $11
		WHERE id = $1
	`

	var lastMaintenance, nextMaintenance sql.NullTime
	if asset.GetLastMaintenanceAt() != nil {
		lastMaintenance = sql.NullTime{Time: *asset.GetLastMaintenanceAt(), Valid: true}
	}
	if asset.GetNextMaintenanceAt() != nil {
		nextMaintenance = sql.NullTime{Time: *asset.GetNextMaintenanceAt(), Valid: true}
	}

	_, err := r.pool.Exec(ctx, query,
		asset.GetID(),
		asset.GetName(),
		asset.GetDescription(),
		asset.GetQuantity(),
		asset.GetAvailableQuantity(),
		asset.GetLocation(),
		string(asset.GetCondition()),
		asset.GetUnitCost().Amount,
		lastMaintenance,
		nextMaintenance,
		asset.GetUpdatedAt(),
	)

	if err != nil {
		r.logger.Error("Failed to update asset",
			zap.String("asset_id", asset.GetID()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update asset: %w", err)
	}

	r.logger.Info("Asset updated successfully",
		zap.String("asset_id", asset.GetID()),
	)

	return nil
}

// Delete soft deletes an asset (not implemented - assets should not be deleted)
func (r *AssetRepository) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("asset deletion is not supported")
}

// GetTotalCount returns the total count of assets
func (r *AssetRepository) GetTotalCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM assets`

	var count int
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to get total asset count", zap.Error(err))
		return 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return count, nil
}

// AddInventoryLog adds an inventory log entry
func (r *AssetRepository) AddInventoryLog(ctx context.Context, log *entities.InventoryLog) error {
	query := `
		INSERT INTO inventory_logs (id, asset_id, change_type, quantity, reason, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		log.ID,
		log.AssetID,
		string(log.ChangeType),
		log.Quantity,
		log.Reason,
		log.CreatedBy,
		log.CreatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to add inventory log",
			zap.String("asset_id", log.AssetID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to add inventory log: %w", err)
	}

	return nil
}

// GetInventoryLogs retrieves inventory logs for an asset
func (r *AssetRepository) GetInventoryLogs(ctx context.Context, assetID string, limit, offset int) ([]*entities.InventoryLog, error) {
	query := `
		SELECT id, asset_id, change_type, quantity, reason, created_by, created_at
		FROM inventory_logs
		WHERE asset_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, assetID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get inventory logs",
			zap.String("asset_id", assetID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get inventory logs: %w", err)
	}
	defer rows.Close()

	var logs []*entities.InventoryLog
	for rows.Next() {
		var log entities.InventoryLog
		var changeType string

		if err := rows.Scan(
			&log.ID,
			&log.AssetID,
			&changeType,
			&log.Quantity,
			&log.Reason,
			&log.CreatedBy,
			&log.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan inventory log row: %w", err)
		}

		log.ChangeType = entities.ChangeType(changeType)
		logs = append(logs, &log)
	}

	return logs, nil
}

// Helper methods

func (r *AssetRepository) mapToDomainEntity(record *assetRecord) (*entities.Asset, error) {
	// Parse category
	category, err := entities.ValidateCategory(record.Category)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Parse condition
	condition, err := entities.ValidateCondition(record.Condition)
	if err != nil {
		return nil, fmt.Errorf("invalid condition: %w", err)
	}

	// Create asset using NewAsset
	asset, err := entities.NewAsset(
		record.Name,
		record.Description,
		category,
		record.Quantity,
		record.Location,
		valueobjects.NewMoney(record.UnitCost),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create asset entity: %w", err)
	}

	// Update fields that can't be set in constructor
	asset.SetCondition(condition)

	// Set maintenance dates if available
	var lastMaintenance, nextMaintenance *time.Time
	if record.LastMaintenanceAt.Valid {
		lastMaintenance = &record.LastMaintenanceAt.Time
	}
	if record.NextMaintenanceAt.Valid {
		nextMaintenance = &record.NextMaintenanceAt.Time
	}
	asset.SetMaintenanceDates(lastMaintenance, nextMaintenance)

	return asset, nil
}
