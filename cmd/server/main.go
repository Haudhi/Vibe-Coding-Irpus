package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net/http"
	"time"
	"encoding/json"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/company/ga-ticketing/src/config"
)

func main() {
	// Initialize logger
	logger := initLogger()
	defer logger.Sync()

	logger.Info("Starting GA Ticketing System")

	// Load configuration
	cfg, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Database connection skipped (using mock data)")

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{
			"status":  "healthy",
			"service": "ga-ticketing",
		})
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handleLogin)
			r.Get("/me", handleGetCurrentUser)
		})

		// Ticket routes (no authentication required)
		r.Route("/tickets", func(r chi.Router) {
			r.Get("/", handleGetTickets)
			r.Post("/", handleCreateTicket)
			r.Get("/{id}", handleGetTicket)
			r.Put("/{id}/status", handleUpdateTicketStatus)
			r.Put("/{id}/approval", handleTicketApproval)
			r.Post("/{id}/comments", handleAddComment)
			r.Get("/{id}/comments", handleGetComments)
		})

		// Asset routes (no authentication required)
		r.Route("/assets", func(r chi.Router) {
			r.Get("/", handleGetAssets)
			r.Post("/", handleCreateAsset)
			r.Put("/{id}/stock", handleUpdateAssetStock)
		})
	})

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.GetServerAddress(),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("GA Ticketing System started successfully",
			zap.String("version", cfg.App.Version),
			zap.String("environment", cfg.App.Env),
			zap.String("server_address", cfg.GetServerAddress()),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-shutdownChan

	logger.Info("Shutting down GA Ticketing System...")

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}

	logger.Info("GA Ticketing System stopped")
}

// Handler functions

// Ticket handlers
func handleGetTickets(w http.ResponseWriter, r *http.Request) {
	// Mock tickets data for testing
	tickets := []map[string]interface{}{
		{
			"id":            "ticket-1",
			"ticket_number": "GA-2025-0001",
			"title":         "Office Supplies Request",
			"description":   "Need to order new pens and notebooks",
			"status":        "pending",
			"priority":      "medium",
			"category":      "office_supplies",
			"estimated_cost": 15000,
			"created_at":    "2025-11-06T10:00:00Z",
			"created_by":    "John Doe",
		},
		{
			"id":            "ticket-2",
			"ticket_number": "GA-2025-0002",
			"title":         "Computer Monitor",
			"description":   "Request for new 24-inch monitor",
			"status":        "approved",
			"priority":      "high",
			"category":      "hardware",
			"estimated_cost": 850000,
			"created_at":    "2025-11-05T14:30:00Z",
			"created_by":    "Jane Smith",
		},
	}

	render.JSON(w, r, map[string]interface{}{
		"tickets": tickets,
		"total":   len(tickets),
	})
}

func handleCreateTicket(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Priority      string `json:"priority"`
		Category      string `json:"category"`
		EstimatedCost int64  `json:"estimated_cost"`
		RequesterName string `json:"requester_name"`
		Department    string `json:"department"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Validation
	if req.Title == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "title is required"})
		return
	}

	if req.RequesterName == "" {
		req.RequesterName = "Anonymous" // Default when not provided
	}

	if req.Department == "" {
		req.Department = "General"
	}

	// Generate ticket number
	ticketNumber := fmt.Sprintf("GA-%d-%04d", time.Now().Year(), 1) // Simplified

	// Determine if approval is required
	requiresApproval := req.EstimatedCost >= 500000 || req.Category == "office_furniture"

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"id":                fmt.Sprintf("ticket-%d", time.Now().Unix()),
		"ticket_number":     ticketNumber,
		"title":             req.Title,
		"description":       req.Description,
		"status":            "pending",
		"priority":          req.Priority,
		"category":          req.Category,
		"estimated_cost":    req.EstimatedCost,
		"requires_approval": requiresApproval,
		"requester_name":    req.RequesterName,
		"department":        req.Department,
		"created_at":        time.Now().Format(time.RFC3339),
		"message":           "Ticket created successfully",
	})
}

func handleGetTicket(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	// For now, return a mock response
	render.JSON(w, r, map[string]interface{}{
		"id":            ticketID,
		"ticket_number": "GA-2025-0001",
		"title":         "Mock Ticket",
		"description":   "This is a mock ticket",
		"status":        "pending",
		"priority":      "medium",
		"category":      "office_supplies",
	})
}

func handleUpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	var req struct {
		Status     string `json:"status"`
		AssignedTo string `json:"assigned_to"`
		ActualCost int64  `json:"actual_cost"`
		Notes      string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "approved", "rejected", "in_progress", "completed", "cancelled"}
	statusValid := false
	for _, status := range validStatuses {
		if req.Status == status {
			statusValid = true
			break
		}
	}

	if !statusValid {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid status. Valid statuses: pending, approved, rejected, in_progress, completed, cancelled"})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"id":           ticketID,
		"status":       req.Status,
		"assigned_to":  req.AssignedTo,
		"actual_cost":  req.ActualCost,
		"notes":        req.Notes,
		"updated_at":   time.Now().Format(time.RFC3339),
		"message":      "Ticket status updated successfully",
	})
}

func handleTicketApproval(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	var req struct {
		Action       string `json:"action"`
		ApproverName string `json:"approver_name"`
		Notes        string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	if req.Action != "approve" && req.Action != "reject" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "action must be 'approve' or 'reject'"})
		return
	}

	if req.ApproverName == "" {
		req.ApproverName = "Anonymous Approver"
	}

	status := "approved"
	if req.Action == "reject" {
		status = "rejected"
	}

	render.JSON(w, r, map[string]interface{}{
		"id":      ticketID,
		"status":  status,
		"message": fmt.Sprintf("Ticket %s successfully", req.Action),
		"approval_info": map[string]interface{}{
			"status":         status,
			"approver_name":  req.ApproverName,
			"approval_notes": req.Notes,
			"approved_at":    time.Now().Format(time.RFC3339),
		},
	})
}

func handleAddComment(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	var req struct {
		Comment string `json:"comment"`
		AuthorName string `json:"author_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	if req.Comment == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "comment is required"})
		return
	}

	if req.AuthorName == "" {
		req.AuthorName = "Anonymous"
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"id":         fmt.Sprintf("comment-%d", time.Now().Unix()),
		"ticket_id":  ticketID,
		"author_name": req.AuthorName,
		"comment":    req.Comment,
		"created_at": time.Now().Format(time.RFC3339),
		"message":    "Comment added successfully",
	})
}

func handleGetComments(w http.ResponseWriter, r *http.Request) {
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Ticket ID is required"})
		return
	}

	// Mock comments for testing
	comments := []map[string]interface{}{
		{
			"id":          "comment-1",
			"ticket_id":   ticketID,
			"author_name": "John Doe",
			"comment":     "This request is urgent. Please prioritize.",
			"created_at":  "2025-11-06T10:30:00Z",
		},
		{
			"id":          "comment-2",
			"ticket_id":   ticketID,
			"author_name": "Jane Smith",
			"comment":     "I have reviewed this request and it looks good to proceed.",
			"created_at":  "2025-11-06T11:15:00Z",
		},
	}

	render.JSON(w, r, map[string]interface{}{
		"comments": comments,
		"total":    len(comments),
	})
}

// Asset handlers
func handleGetAssets(w http.ResponseWriter, r *http.Request) {
	// Mock assets data for testing
	assets := []map[string]interface{}{
		{
			"id":          "asset-1",
			"asset_code":  "GA-SUP-001",
			"name":        "Office Chair",
			"description": "Ergonomic office chair with lumbar support",
			"category":    "furniture",
			"quantity":    15,
			"available":   12,
			"location":    "Storage Room A",
			"condition":   "good",
			"created_at":  "2025-11-01T10:00:00Z",
		},
		{
			"id":          "asset-2",
			"asset_code":  "GA-EQP-002",
			"name":        "Laptop Dell XPS 15",
			"description": "High-performance laptop for development work",
			"category":    "hardware",
			"quantity":    5,
			"available":   2,
			"location":    "IT Department",
			"condition":   "excellent",
			"created_at":  "2025-11-02T14:30:00Z",
		},
		{
			"id":          "asset-3",
			"asset_code":  "GA-SUP-003",
			"name":        "Printer Paper A4",
			"description": "Standard A4 printing paper, 500 sheets per ream",
			"category":    "office_supplies",
			"quantity":    100,
			"available":   87,
			"location":    "Supply Cabinet",
			"condition":   "new",
			"created_at":  "2025-11-03T09:15:00Z",
		},
	}

	render.JSON(w, r, map[string]interface{}{
		"assets": assets,
		"total":  len(assets),
	})
}

func handleCreateAsset(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Quantity    int    `json:"quantity"`
		Location    string `json:"location"`
		Condition   string `json:"condition"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Validation
	if req.Name == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "name is required"})
		return
	}

	if req.Quantity <= 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "quantity must be greater than 0"})
		return
	}

	if req.Category == "" {
		req.Category = "general"
	}

	if req.Location == "" {
		req.Location = "Not specified"
	}

	if req.Condition == "" {
		req.Condition = "good"
	}

	// Generate asset code
	categoryCode := "GEN"
	switch req.Category {
	case "furniture":
		categoryCode = "FUR"
	case "hardware":
		categoryCode = "EQP"
	case "office_supplies":
		categoryCode = "SUP"
	}

	assetID := fmt.Sprintf("asset-%d", time.Now().Unix())
	assetCode := fmt.Sprintf("GA-%s-%04d", categoryCode, 1) // Simplified

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"id":          assetID,
		"asset_code":  assetCode,
		"name":        req.Name,
		"description": req.Description,
		"category":    req.Category,
		"quantity":    req.Quantity,
		"available":   req.Quantity, // All available when created
		"location":    req.Location,
		"condition":   req.Condition,
		"created_at":  time.Now().Format(time.RFC3339),
		"message":     "Asset added successfully",
	})
}

func handleUpdateAssetStock(w http.ResponseWriter, r *http.Request) {
	assetID := chi.URLParam(r, "id")
	if assetID == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Asset ID is required"})
		return
	}

	var req struct {
		Quantity    int    `json:"quantity"`
		Operation   string `json:"operation"` // "add", "subtract", "set"
		Notes       string `json:"notes"`
		UpdatedBy   string `json:"updated_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	if req.Quantity < 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "quantity cannot be negative"})
		return
	}

	if req.Operation == "" {
		req.Operation = "set" // Default operation
	}

	if req.UpdatedBy == "" {
		req.UpdatedBy = "Anonymous"
	}

	// Validate operation
	validOps := []string{"add", "subtract", "set"}
	opValid := false
	for _, op := range validOps {
		if req.Operation == op {
			opValid = true
			break
		}
	}

	if !opValid {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "operation must be 'add', 'subtract', or 'set'"})
		return
	}

	// Mock calculation
	var newQuantity, previousQuantity int
	switch req.Operation {
	case "add":
		previousQuantity = 10 // Mock previous
		newQuantity = previousQuantity + req.Quantity
	case "subtract":
		previousQuantity = 10 // Mock previous
		newQuantity = previousQuantity - req.Quantity
		if newQuantity < 0 {
			newQuantity = 0
		}
	case "set":
		previousQuantity = 10 // Mock previous
		newQuantity = req.Quantity
	}

	render.JSON(w, r, map[string]interface{}{
		"id":                assetID,
		"quantity":          newQuantity,
		"previous_quantity": previousQuantity,
		"operation":         req.Operation,
		"notes":             req.Notes,
		"updated_by":        req.UpdatedBy,
		"updated_at":        time.Now().Format(time.RFC3339),
		"message":           "Stock updated successfully",
	})
}

// Auth handlers
func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	// Validate input
	if req.Email == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "email is required"})
		return
	}
	if req.Password == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "password is required"})
		return
	}

	// Mock user authentication based on test credentials from documentation
	validUsers := map[string]map[string]interface{}{
		"requester@company.com": {
			"id":         "1",
			"name":       "John Doe",
			"role":       "requester",
			"department": "Finance",
			"password":   "password123",
		},
		"approver@company.com": {
			"id":         "2",
			"name":       "Jane Approver",
			"role":       "approver",
			"department": "Management",
			"password":   "password123",
		},
		"admin@company.com": {
			"id":         "3",
			"name":       "Admin GA",
			"role":       "admin",
			"department": "General Affairs",
			"password":   "password123",
		},
	}

	userData, exists := validUsers[req.Email]
	if !exists || userData["password"] != req.Password {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Invalid credentials"})
		return
	}

	// Generate mock JWT token (in production, use proper JWT library)
	token := fmt.Sprintf("mock-jwt-token-%d-%s", time.Now().Unix(), req.Email)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":         userData["id"],
			"email":      req.Email,
			"name":       userData["name"],
			"role":       userData["role"],
			"department": userData["department"],
			"created_at": "2025-01-15T10:30:00Z",
		},
	})
}

func handleGetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, extract and validate JWT token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Unauthorized"})
		return
	}

	// Mock user response based on token
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Extract email from mock token
	if strings.Contains(token, "requester@company.com") {
		render.JSON(w, r, map[string]interface{}{
			"id":         "1",
			"email":      "requester@company.com",
			"name":       "John Doe",
			"role":       "requester",
			"department": "Finance",
			"created_at": "2025-01-15T10:30:00Z",
		})
	} else if strings.Contains(token, "approver@company.com") {
		render.JSON(w, r, map[string]interface{}{
			"id":         "2",
			"email":      "approver@company.com",
			"name":       "Jane Approver",
			"role":       "approver",
			"department": "Management",
			"created_at": "2025-01-15T10:30:00Z",
		})
	} else if strings.Contains(token, "admin@company.com") {
		render.JSON(w, r, map[string]interface{}{
			"id":         "3",
			"email":      "admin@company.com",
			"name":       "Admin GA",
			"role":       "admin",
			"department": "General Affairs",
			"created_at": "2025-01-15T10:30:00Z",
		})
	} else {
		render.Status(r, http.StatusUnauthorized)
		render.JSON(w, r, map[string]string{"error": "Invalid token"})
	}
}

// initLogger initializes the application logger
func initLogger() *zap.Logger {
	var zapConfig zap.Config

	// Set log level based on environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "json"
	}

	if logFormat == "console" {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}

	// Set log level
	switch logLevel {
	case "debug":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configure encoder
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.StacktraceKey = "stacktrace"

	logger, err := zapConfig.Build()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	return logger
}