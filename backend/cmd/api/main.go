package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dishdice/backend/internal/ai"
	"github.com/dishdice/backend/internal/config"
	"github.com/dishdice/backend/internal/database"
	"github.com/dishdice/backend/internal/handlers"
	"github.com/dishdice/backend/internal/middleware"
	"github.com/dishdice/backend/internal/repository"
	"github.com/dishdice/backend/internal/services"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Run migrations
	migrationsPath := filepath.Join("migrations")
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	proposalRepo := repository.NewProposalRepository(db)
	shoppingRepo := repository.NewShoppingRepository(db)

	// Initialize AI client
	aiClient := ai.NewClient(cfg.OpenAIAPIKey)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	proposalService := services.NewProposalService(proposalRepo, userRepo, aiClient)
	mealService := services.NewMealService(proposalRepo, userRepo, aiClient)
	shoppingService := services.NewShoppingService(shoppingRepo, proposalRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, userRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	proposalHandler := handlers.NewProposalHandler(proposalService, shoppingService)
	mealHandler := handlers.NewMealHandler(mealService, shoppingService)
	shoppingHandler := handlers.NewShoppingHandler(shoppingService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(middleware.Logging)
	r.Use(chimiddleware.Recoverer)

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.AllowedOrigins, "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Auth routes (public)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(authService))

			// Auth
			r.Get("/auth/me", authHandler.GetMe)

			// User preferences
			r.Get("/user/preferences", userHandler.GetPreferences)
			r.Put("/user/preferences", userHandler.UpdatePreferences)

			// Proposals
			r.Get("/proposals", proposalHandler.ListProposals)
			r.Post("/proposals", proposalHandler.CreateProposal)
			r.Get("/proposals/{id}", proposalHandler.GetProposal)
			r.Delete("/proposals/{id}", proposalHandler.DeleteProposal)
			r.Post("/proposals/{id}/save-to-shopping", proposalHandler.AddToShoppingList)

			// Meals
			r.Post("/meals/{id}/regenerate", mealHandler.RegenerateMeal)
			r.Put("/meals/{id}/select", mealHandler.SelectMealOption)
			r.Post("/meals/{id}/save-to-shopping", mealHandler.AddToShoppingList)

			// Shopping list
			r.Get("/shopping-list", shoppingHandler.GetShoppingList)
			r.Post("/shopping-list", shoppingHandler.AddItem)
			r.Put("/shopping-list/{id}/toggle", shoppingHandler.ToggleItem)
			r.Delete("/shopping-list/checked", shoppingHandler.DeleteChecked)
			r.Delete("/shopping-list/{id}", shoppingHandler.DeleteItem)
		})
	})

	// Serve static files (for production deployment)
	staticDir := "static"
	if _, err := os.Stat(staticDir); err == nil {
		log.Println("Serving static files from /static")

		// File server for static assets
		fileServer := http.FileServer(http.Dir(staticDir))

		// Serve static files, with SPA fallback
		r.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			path := req.URL.Path

			// Check if file exists
			fullPath := filepath.Join(staticDir, path)
			if _, err := os.Stat(fullPath); err == nil {
				// File exists, serve it
				http.StripPrefix("/", fileServer).ServeHTTP(w, req)
				return
			}

			// Check if it's a file request (has extension)
			if strings.Contains(path, ".") {
				http.NotFound(w, req)
				return
			}

			// SPA fallback - serve index.html
			http.ServeFile(w, req, filepath.Join(staticDir, "index.html"))
		})
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
