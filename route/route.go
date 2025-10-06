package route

import (
	"log"
	"net/http"
	"os"
	"time"

	"jmrashed/apps/userApp/database"
	"jmrashed/apps/userApp/handlers"
	"jmrashed/apps/userApp/middleware"
	"jmrashed/apps/userApp/repository"
	"jmrashed/apps/userApp/seeder"
	"jmrashed/apps/userApp/service"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"golang.org/x/time/rate"
)

// SetupRoutes configures all application routes
func SetupRoutes() {
	// Initialize database
	dbConfig := database.GetDefaultConfig()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitializeSchema(); err != nil {
		log.Printf("Warning: Failed to initialize schema: %v", err)
	}

	// Run seeder
	seederInstance := seeder.NewSeeder(db.DB)
	if err := seederInstance.Run(); err != nil {
		log.Printf("Warning: Failed to run seeder: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	todoRepo := repository.NewTodoRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo)
	todoService := service.NewTodoService(todoRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	todoHandler := handlers.NewTodoHandler(todoService)
	healthHandler := handlers.NewHealthHandler(db.DB)

	// Initialize middleware
	rateLimiter := middleware.NewRateLimiter(rate.Every(time.Minute), 60) // 60 requests per minute
	cache := middleware.NewCache(5 * time.Minute) // 5 minute cache

	// Setup router
	router := mux.NewRouter().StrictSlash(true)

	// Apply global middleware
	router.Use(middleware.CORS)
	router.Use(middleware.Logging)
	router.Use(middleware.RateLimit(rateLimiter))

	// Static files with caching
	staticRouter := router.PathPrefix("/static/").Subrouter()
	staticRouter.Use(middleware.CacheMiddleware(cache))
	staticRouter.PathPrefix("/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Health check
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// Serve swagger.yaml file (must be before Swagger UI route)
	router.HandleFunc("/docs/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/yaml")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeFile(w, r, "docs/swagger.yaml")
	}).Methods("GET")

	// Swagger UI
	router.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/static/swagger.yaml"),
	))

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public routes (no authentication required)
	public := api.PathPrefix("").Subrouter()
	public.HandleFunc("/register", authHandler.Register).Methods("POST")
	public.HandleFunc("/login", authHandler.Login).Methods("POST")
	public.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

	// Protected routes (authentication required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// User profile routes
	protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	protected.HandleFunc("/profile", authHandler.UpdateProfile).Methods("PUT")
	protected.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")
	protected.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	protected.HandleFunc("/logout-all", authHandler.LogoutAll).Methods("POST")

	// Todo routes with permission-based access
	todos := protected.PathPrefix("/todos").Subrouter()
	todos.Use(middleware.RequirePermission("read_todos"))
	todos.HandleFunc("", todoHandler.GetUserTodos).Methods("GET")
	todos.HandleFunc("/{id:[0-9]+}", todoHandler.GetTodo).Methods("GET")
	
	// Todo creation/modification requires write permission
	todosWrite := todos.PathPrefix("").Subrouter()
	todosWrite.Use(middleware.RequirePermission("write_todos"))
	todosWrite.HandleFunc("", todoHandler.CreateTodo).Methods("POST")
	todosWrite.HandleFunc("/{id:[0-9]+}", todoHandler.UpdateTodo).Methods("PUT")
	
	// Todo deletion requires delete permission
	todosDelete := todos.PathPrefix("").Subrouter()
	todosDelete.Use(middleware.RequirePermission("delete_todos"))
	todosDelete.HandleFunc("/{id:[0-9]+}", todoHandler.DeleteTodo).Methods("DELETE")

	// Admin routes (admin role required)
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireRole("admin"))
	admin.HandleFunc("/todos", todoHandler.GetAllTodos).Methods("GET")

	// Moderator routes (moderator or admin role required)
	moderator := protected.PathPrefix("/moderator").Subrouter()
	moderator.Use(middleware.RequireAnyRole("moderator", "admin"))
	// Add moderator-specific routes here

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	log.Println("Features enabled: Authentication, Authorization, Rate Limiting, Caching, Logging")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
