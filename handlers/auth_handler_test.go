package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jmrashed/apps/userApp/model"
	"jmrashed/apps/userApp/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(req model.RegisterRequest) (*model.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*model.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Login(req model.LoginRequest) (*model.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*model.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(req model.RefreshTokenRequest) (*model.AuthResponse, error) {
	args := m.Called(req)
	return args.Get(0).(*model.AuthResponse), args.Error(1)
}

func (m *MockAuthService) Logout(userID int, refreshToken string) error {
	args := m.Called(userID, refreshToken)
	return args.Error(0)
}

func (m *MockAuthService) LogoutAll(userID int) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockAuthService) GetUserProfile(userID int) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) UpdateUserProfile(userID int, username, email string) (*model.User, error) {
	args := m.Called(userID, username, email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockAuthService) ChangePassword(userID int, currentPassword, newPassword string) error {
	args := m.Called(userID, currentPassword, newPassword)
	return args.Error(0)
}

func TestAuthHandler_Register(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Successful registration",
			requestBody: model.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func() {
				mockService.On("Register", mock.AnythingOfType("model.RegisterRequest")).Return(
					&model.AuthResponse{
						User: model.User{
							ID:       1,
							Username: "testuser",
							Email:    "test@example.com",
						},
						AccessToken:  "access_token",
						RefreshToken: "refresh_token",
						ExpiresIn:    900,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Register(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Successful login",
			requestBody: model.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			mockSetup: func() {
				mockService.On("Login", mock.AnythingOfType("model.LoginRequest")).Return(
					&model.AuthResponse{
						User: model.User{
							ID:       1,
							Username: "testuser",
							Email:    "test@example.com",
						},
						AccessToken:  "access_token",
						RefreshToken: "refresh_token",
						ExpiresIn:    900,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Login(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	mockService := new(MockAuthService)
	handler := NewAuthHandler(mockService)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Successful token refresh",
			requestBody: model.RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			mockSetup: func() {
				mockService.On("RefreshToken", mock.AnythingOfType("model.RefreshTokenRequest")).Return(
					&model.AuthResponse{
						User: model.User{
							ID:       1,
							Username: "testuser",
							Email:    "test@example.com",
						},
						AccessToken:  "new_access_token",
						RefreshToken: "new_refresh_token",
						ExpiresIn:    900,
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			requestBody:    "invalid json",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			var body []byte
			var err error
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.RefreshToken(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}