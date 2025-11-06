package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/company/ga-ticketing/src/application/dto"
	"github.com/company/ga-ticketing/src/application/usecases"
	"github.com/company/ga-ticketing/src/domain/entities"
)

// AssetHandler handles asset-related HTTP requests
type AssetHandler struct {
	createAssetUC      *usecases.CreateAssetUseCase
	getAssetUC         *usecases.GetAssetUseCase
	getAssetsUC        *usecases.GetAssetsUseCase
	updateAssetUC      *usecases.UpdateAssetUseCase
	updateInventoryUC  *usecases.UpdateInventoryUseCase
	logger             *zap.Logger
}

// NewAssetHandler creates a new AssetHandler
func NewAssetHandler(
	createAssetUC *usecases.CreateAssetUseCase,
	getAssetUC *usecases.GetAssetUseCase,
	getAssetsUC *usecases.GetAssetsUseCase,
	updateAssetUC *usecases.UpdateAssetUseCase,
	updateInventoryUC *usecases.UpdateInventoryUseCase,
	logger *zap.Logger,
) *AssetHandler {
	return &AssetHandler{
		createAssetUC:     createAssetUC,
		getAssetUC:        getAssetUC,
		getAssetsUC:       getAssetsUC,
		updateAssetUC:     updateAssetUC,
		updateInventoryUC: updateInventoryUC,
		logger:            logger,
	}
}

// GetAssets handles GET /v1/assets
func (h *AssetHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	// Check permissions - only admin can view assets
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	if userRole != string(entities.RoleAdmin) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, map[string]string{"error": "Only admins can access assets"})
		return
	}

	// Get pagination parameters
	page := 1
	limit := 20

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	req := &dto.GetAssetsRequest{
		Page:      page,
		Limit:     limit,
		Category:  r.URL.Query().Get("category"),
		Condition: r.URL.Query().Get("condition"),
	}

	response, err := h.getAssetsUC.Execute(r.Context(), req)
	if err != nil {
		h.logger.Error("Failed to get assets", zap.Error(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}

// CreateAsset handles POST /v1/assets
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	// Check permissions - only admin can create assets
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	if userRole != string(entities.RoleAdmin) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, map[string]string{"error": "Only admins can create assets"})
		return
	}

	var req dto.CreateAssetRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		h.logger.Warn("Failed to decode asset request", zap.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	asset, err := h.createAssetUC.Execute(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create asset", zap.Error(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info("Asset created successfully",
		zap.String("asset_id", asset.ID),
		zap.String("asset_code", asset.AssetCode),
	)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, asset)
}

// GetAsset handles GET /v1/assets/{assetId}
func (h *AssetHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	// Check permissions - only admin can view assets
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	if userRole != string(entities.RoleAdmin) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, map[string]string{"error": "Only admins can access assets"})
		return
	}

	assetID := chi.URLParam(r, "assetId")
	if assetID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Asset ID is required"})
		return
	}

	asset, err := h.getAssetUC.Execute(r.Context(), assetID)
	if err != nil {
		h.logger.Error("Failed to get asset",
			zap.String("asset_id", assetID),
			zap.Error(err),
		)

		if err.Error() == "asset not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Asset not found"})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, asset)
}

// UpdateAsset handles PUT /v1/assets/{assetId}
func (h *AssetHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	// Check permissions - only admin can update assets
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	if userRole != string(entities.RoleAdmin) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, map[string]string{"error": "Only admins can update assets"})
		return
	}

	assetID := chi.URLParam(r, "assetId")
	if assetID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Asset ID is required"})
		return
	}

	var req dto.UpdateAssetRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	asset, err := h.updateAssetUC.Execute(r.Context(), assetID, &req)
	if err != nil {
		h.logger.Error("Failed to update asset",
			zap.String("asset_id", assetID),
			zap.Error(err),
		)

		if err.Error() == "asset not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Asset not found"})
			return
		}

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info("Asset updated successfully", zap.String("asset_id", assetID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, asset)
}

// UpdateInventory handles POST /v1/assets/{assetId}/inventory
func (h *AssetHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	// Check permissions - only admin can update inventory
	userRole, ok := r.Context().Value("user_role").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User role not found"})
		return
	}

	if userRole != string(entities.RoleAdmin) {
		render.Status(r, http.StatusForbidden)
		render.JSON(w, r, map[string]string{"error": "Only admins can update inventory"})
		return
	}

	assetID := chi.URLParam(r, "assetId")
	if assetID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Asset ID is required"})
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "User not authenticated"})
		return
	}

	var req dto.UpdateInventoryRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	asset, err := h.updateInventoryUC.Execute(r.Context(), assetID, userID, &req)
	if err != nil {
		h.logger.Error("Failed to update inventory",
			zap.String("asset_id", assetID),
			zap.String("user_id", userID),
			zap.Error(err),
		)

		if err.Error() == "asset not found" {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Asset not found"})
			return
		}

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	h.logger.Info("Inventory updated successfully",
		zap.String("asset_id", assetID),
		zap.String("user_id", userID),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, asset)
}
