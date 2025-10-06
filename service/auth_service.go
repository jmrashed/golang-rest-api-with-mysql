package service

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"jmrashed/apps/userApp/auth"
	"jmrashed/apps/userApp/model"
	"jmrashed/apps/userApp/repository"

	"github.com/go-playground/validator/v10"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	validator *validator.Validate
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		validator: validator.New(),
	}
}

// Register creates a new user account
func (s *AuthService) Register(req model.RegisterRequest) (*model.AuthResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if user already exists
	if _, err := s.userRepo.GetUserByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}
	if _, err := s.userRepo.GetUserByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Assign default role (user)
	if err := s.userRepo.AssignRoleToUser(user.ID, 2); err != nil {
		return nil, fmt.Errorf("failed to assign default role: %w", err)
	}

	// Load user with roles and permissions
	userWithRoles, err := s.userRepo.GetUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user roles: %w", err)
	}

	// Generate tokens
	authResponse, err := auth.GenerateTokens(*userWithRoles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token
	refreshTokenHash := s.hashToken(authResponse.RefreshToken)
	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.userRepo.StoreRefreshToken(refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Remove password hash from response
	authResponse.User.PasswordHash = ""
	return authResponse, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req model.LoginRequest) (*model.AuthResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get user by username
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check password
	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Generate tokens
	authResponse, err := auth.GenerateTokens(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store refresh token
	refreshTokenHash := s.hashToken(authResponse.RefreshToken)
	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.userRepo.StoreRefreshToken(refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Remove password hash from response
	authResponse.User.PasswordHash = ""
	return authResponse, nil
}

// RefreshToken generates new access token using refresh token
func (s *AuthService) RefreshToken(req model.RefreshTokenRequest) (*model.AuthResponse, error) {
	// Validate request
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Validate refresh token
	_, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Check if refresh token exists in database
	refreshTokenHash := s.hashToken(req.RefreshToken)
	storedToken, err := s.userRepo.GetRefreshToken(refreshTokenHash)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	// Get user with roles and permissions
	user, err := s.userRepo.GetUserByID(storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Generate new tokens
	authResponse, err := auth.GenerateTokens(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Delete old refresh token and store new one
	if err := s.userRepo.DeleteRefreshToken(refreshTokenHash); err != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	newRefreshTokenHash := s.hashToken(authResponse.RefreshToken)
	newRefreshToken := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.userRepo.StoreRefreshToken(newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	// Remove password hash from response
	authResponse.User.PasswordHash = ""
	return authResponse, nil
}

// Logout invalidates refresh token
func (s *AuthService) Logout(userID int, refreshToken string) error {
	refreshTokenHash := s.hashToken(refreshToken)
	return s.userRepo.DeleteRefreshToken(refreshTokenHash)
}

// LogoutAll invalidates all refresh tokens for a user
func (s *AuthService) LogoutAll(userID int) error {
	return s.userRepo.DeleteUserRefreshTokens(userID)
}

// GetUserProfile returns user profile information
func (s *AuthService) GetUserProfile(userID int) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// Remove password hash
	user.PasswordHash = ""
	return user, nil
}

// UpdateUserProfile updates user profile information
func (s *AuthService) UpdateUserProfile(userID int, username, email string) (*model.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if username is taken by another user
	if username != user.Username {
		if existingUser, err := s.userRepo.GetUserByUsername(username); err == nil && existingUser.ID != userID {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email is taken by another user
	if email != user.Email {
		if existingUser, err := s.userRepo.GetUserByEmail(email); err == nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
		}
	}

	user.Username = username
	user.Email = email

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Remove password hash
	user.PasswordHash = ""
	return user, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID int, currentPassword, newPassword string) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if !auth.CheckPassword(currentPassword, user.PasswordHash) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = hashedPassword
	if err := s.userRepo.UpdateUser(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate all refresh tokens to force re-login
	return s.userRepo.DeleteUserRefreshTokens(userID)
}

// hashToken creates a SHA256 hash of the token for storage
func (s *AuthService) hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}