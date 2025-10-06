package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"jmrashed/apps/userApp/model"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck performs comprehensive health check
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	services := make(map[string]string)
	
	// Check database connection
	if err := h.db.Ping(); err != nil {
		services["database"] = "unhealthy"
	} else {
		services["database"] = "healthy"
	}
	
	// Determine overall status
	status := "healthy"
	for _, serviceStatus := range services {
		if serviceStatus == "unhealthy" {
			status = "unhealthy"
			break
		}
	}
	
	response := model.HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Version:   "1.0.0",
		Services:  services,
	}
	
	statusCode := http.StatusOK
	if status == "unhealthy" {
		statusCode = http.StatusServiceUnavailable
	}
	
	writeSuccessResponse(w, statusCode, "Health check completed", response)
}