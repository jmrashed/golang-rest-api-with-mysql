package auth

import (
	"testing"
	"time"

	"jmrashed/apps/userApp/model"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"
	
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"
	wrongPassword := "wrongpassword"
	
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	
	// Test correct password
	assert.True(t, CheckPassword(password, hash))
	
	// Test wrong password
	assert.False(t, CheckPassword(wrongPassword, hash))
}

func TestGenerateTokens(t *testing.T) {
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
	
	authResponse, err := GenerateTokens(user)
	assert.NoError(t, err)
	assert.NotNil(t, authResponse)
	assert.NotEmpty(t, authResponse.AccessToken)
	assert.NotEmpty(t, authResponse.RefreshToken)
	assert.Equal(t, user.ID, authResponse.User.ID)
	assert.Greater(t, authResponse.ExpiresIn, int64(0))
}

func TestValidateAccessToken(t *testing.T) {
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
	
	authResponse, err := GenerateTokens(user)
	assert.NoError(t, err)
	
	claims, err := ValidateAccessToken(authResponse.AccessToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.Email, claims.Email)
	assert.Contains(t, claims.Roles, "user")
	assert.Contains(t, claims.Permissions, "read_todos")
}

func TestValidateRefreshToken(t *testing.T) {
	user := model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}
	
	authResponse, err := GenerateTokens(user)
	assert.NoError(t, err)
	
	claims, err := ValidateRefreshToken(authResponse.RefreshToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.NotEmpty(t, claims.JTI)
}

func TestValidateInvalidToken(t *testing.T) {
	invalidToken := "invalid.token.here"
	
	_, err := ValidateAccessToken(invalidToken)
	assert.Error(t, err)
	
	_, err = ValidateRefreshToken(invalidToken)
	assert.Error(t, err)
}

func TestHasPermission(t *testing.T) {
	permissions := []string{"read_todos", "write_todos"}
	
	assert.True(t, HasPermission(permissions, "read_todos"))
	assert.True(t, HasPermission(permissions, "write_todos"))
	assert.False(t, HasPermission(permissions, "delete_todos"))
}

func TestHasRole(t *testing.T) {
	roles := []string{"user", "moderator"}
	
	assert.True(t, HasRole(roles, "user"))
	assert.True(t, HasRole(roles, "moderator"))
	assert.False(t, HasRole(roles, "admin"))
}

func TestGenerateSecureToken(t *testing.T) {
	token1, err := GenerateSecureToken(32)
	assert.NoError(t, err)
	assert.Len(t, token1, 64) // hex encoding doubles the length
	
	token2, err := GenerateSecureToken(32)
	assert.NoError(t, err)
	assert.NotEqual(t, token1, token2) // Should be different each time
}