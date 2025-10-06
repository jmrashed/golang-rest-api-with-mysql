package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"jmrashed/apps/userApp/middleware"
	"jmrashed/apps/userApp/model"
	"jmrashed/apps/userApp/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	authResponse, err := h.authService.Register(req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusCreated, "User registered successfully", authResponse)
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	authResponse, err := h.authService.Login(req)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Login successful", authResponse)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshTokenRequest
	
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	authResponse, err := h.authService.RefreshToken(req)
	if err != nil {
		writeErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Token refreshed successfully", authResponse)
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	var req model.RefreshTokenRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := h.authService.Logout(claims.UserID, req.RefreshToken); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Logout successful", nil)
}

// LogoutAll handles logout from all devices
func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	if err := h.authService.LogoutAll(claims.UserID); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to logout from all devices")
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Logged out from all devices", nil)
}

// GetProfile returns user profile
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	user, err := h.authService.GetUserProfile(claims.UserID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Profile retrieved successfully", user)
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.authService.UpdateUserProfile(claims.UserID, req.Username, req.Email)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Profile updated successfully", user)
}

// ChangePassword changes user password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if err := h.authService.ChangePassword(claims.UserID, req.CurrentPassword, req.NewPassword); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Password changed successfully", nil)
}

// Helper functions
func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func writeSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.SuccessResponse{
		Message: message,
		Data:    data,
	})
}