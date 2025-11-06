package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	authmiddleware "github.com/company/ga-ticketing/src/interface/http/middleware"
	"github.com/company/ga-ticketing/src/interface/http/handlers"
)

// SetupRouter configures and returns the main application router
func SetupRouter(
	ticketHandler *handlers.TicketHandler,
	assetHandler *handlers.AssetHandler,
	approvalHandler *handlers.ApprovalHandler,
	commentHandler *handlers.CommentHandler,
) *chi.Mux {
	r := chi.NewRouter()

	// Standard middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// CORS middleware
	r.Use(corsMiddleware)

	// Health check endpoint
	r.Get("/health", healthCheck)

	// API version 1 routes
	r.Route("/v1", func(r chi.Router) {
		// Authentication middleware for all v1 routes
		r.Use(middleware.AuthMiddleware)

		// Ticket routes
		r.Route("/tickets", func(r chi.Router) {
			r.Get("/", ticketHandler.GetTickets)
			r.Post("/", ticketHandler.CreateTicket)

			r.Route("/{ticketID}", func(r chi.Router) {
				r.Get("/", ticketHandler.GetTicket)
				r.Put("/", ticketHandler.UpdateTicket)
				r.Post("/assign", ticketHandler.AssignTicket)

				// Comment routes
				r.Route("/comments", func(r chi.Router) {
					r.Get("/", commentHandler.GetComments)
					r.Post("/", commentHandler.AddComment)
				})

				// Approval routes
				r.Route("/approval", func(r chi.Router) {
					r.Use(approvalMiddleware) // Only approvers/admins can access
					r.Post("/approve", approvalHandler.ApproveTicket)
					r.Post("/reject", approvalHandler.RejectTicket)
				})
			})
		})

		// Asset routes (admin only)
		r.Route("/assets", func(r chi.Router) {
			r.Use(adminMiddleware) // Only admins can access assets

			r.Get("/", assetHandler.GetAssets)
			r.Post("/", assetHandler.CreateAsset)

			r.Route("/{assetID}", func(r chi.Router) {
				r.Get("/", assetHandler.GetAsset)
				r.Put("/", assetHandler.UpdateAsset)
				r.Post("/inventory", assetHandler.UpdateInventory)
			})
		})

		// Approval routes for approvers
		r.Route("/approvals", func(r chi.Router) {
			r.Use(approvalMiddleware) // Only approvers/admins can access
			r.Get("/pending", approvalHandler.GetPendingApprovals)
		})
	})

	return r
}

// healthCheck provides a simple health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, http.StatusOK, map[string]string{
		"status": "healthy",
		"service": "ga-ticketing",
	})
}

// corsMiddleware provides CORS configuration
func corsMiddleware(next http.Handler) http.Handler {
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
}

// adminMiddleware restricts access to admin users only
func adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := authmiddleware.GetUserFromContext(r.Context())
		if err != nil {
			render.Render(w, r, authmiddleware.ErrUnauthorized("user not found in context"))
			return
		}

		if user.Role != "admin" {
			render.Render(w, r, authmiddleware.ErrForbidden("admin access required"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// approvalMiddleware restricts access to approvers and admins
func approvalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := authmiddleware.GetUserFromContext(r.Context())
		if err != nil {
			render.Render(w, r, authmiddleware.ErrUnauthorized("user not found in context"))
			return
		}

		if user.Role != "approver" && user.Role != "admin" {
			render.Render(w, r, authmiddleware.ErrForbidden("approver access required"))
			return
		}

		next.ServeHTTP(w, r)
	})
}