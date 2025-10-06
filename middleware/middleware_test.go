package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"jmrashed/apps/userApp/auth"
	"jmrashed/apps/userApp/model"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Create a test user and generate tokens
	user := model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Roles: []model.Role{
			{
				ID:   1,
				Name: "user",
				Permissions: []model.Permission{
					{
						ID:   1,
						Name: "read_todos",
					},
				},
			},
		},
	}

	authResponse, err := auth.GenerateTokens(user)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + authResponse.AccessToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid authorization format",
			authHeader:     "InvalidFormat",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create middleware
			middleware := AuthMiddleware(testHandler)

			// Create request
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute middleware
			middleware.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRequirePermission(t *testing.T) {
	// Create test claims with permissions
	claims := &auth.Claims{
		UserID:      1,
		Username:    "testuser",
		Permissions: []string{"read_todos", "write_todos"},
	}

	tests := []struct {
		name               string
		requiredPermission string
		userClaims         *auth.Claims
		expectedStatus     int
	}{
		{
			name:               "User has required permission",
			requiredPermission: "read_todos",
			userClaims:         claims,
			expectedStatus:     http.StatusOK,
		},
		{
			name:               "User lacks required permission",
			requiredPermission: "delete_todos",
			userClaims:         claims,
			expectedStatus:     http.StatusForbidden,
		},
		{
			name:               "No user context",
			requiredPermission: "read_todos",
			userClaims:         nil,
			expectedStatus:     http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create middleware
			middleware := RequirePermission(tt.requiredPermission)(testHandler)

			// Create request with context
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.userClaims != nil {
				ctx := context.WithValue(req.Context(), UserContextKey, tt.userClaims)
				req = req.WithContext(ctx)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute middleware
			middleware.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRequireRole(t *testing.T) {
	// Create test claims with roles
	claims := &auth.Claims{
		UserID:   1,
		Username: "testuser",
		Roles:    []string{"user", "moderator"},
	}

	tests := []struct {
		name           string
		requiredRole   string
		userClaims     *auth.Claims
		expectedStatus int
	}{
		{
			name:           "User has required role",
			requiredRole:   "user",
			userClaims:     claims,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User lacks required role",
			requiredRole:   "admin",
			userClaims:     claims,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "No user context",
			requiredRole:   "user",
			userClaims:     nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create middleware
			middleware := RequireRole(tt.requiredRole)(testHandler)

			// Create request with context
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.userClaims != nil {
				ctx := context.WithValue(req.Context(), UserContextKey, tt.userClaims)
				req = req.WithContext(ctx)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute middleware
			middleware.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRequireAnyRole(t *testing.T) {
	// Create test claims with roles
	claims := &auth.Claims{
		UserID:   1,
		Username: "testuser",
		Roles:    []string{"user"},
	}

	tests := []struct {
		name           string
		requiredRoles  []string
		userClaims     *auth.Claims
		expectedStatus int
	}{
		{
			name:           "User has one of required roles",
			requiredRoles:  []string{"admin", "user"},
			userClaims:     claims,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User lacks all required roles",
			requiredRoles:  []string{"admin", "moderator"},
			userClaims:     claims,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create middleware
			middleware := RequireAnyRole(tt.requiredRoles...)(testHandler)

			// Create request with context
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.userClaims != nil {
				ctx := context.WithValue(req.Context(), UserContextKey, tt.userClaims)
				req = req.WithContext(ctx)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute middleware
			middleware.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	claims := &auth.Claims{
		UserID:   1,
		Username: "testuser",
	}

	// Test with valid context
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, claims)
	req = req.WithContext(ctx)

	retrievedClaims, ok := GetUserFromContext(req)
	assert.True(t, ok)
	assert.Equal(t, claims.UserID, retrievedClaims.UserID)
	assert.Equal(t, claims.Username, retrievedClaims.Username)

	// Test with empty context
	req2 := httptest.NewRequest("GET", "/test", nil)
	_, ok2 := GetUserFromContext(req2)
	assert.False(t, ok2)
}