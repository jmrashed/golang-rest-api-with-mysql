package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"jmrashed/apps/userApp/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (r *UserRepository) CreateUser(user *model.User) error {
	query := `INSERT INTO users (username, email, password_hash, is_active) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}
	
	user.ID = int(id)
	return nil
}

// GetUserByID retrieves a user by ID with roles and permissions
func (r *UserRepository) GetUserByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at 
			  FROM users WHERE id = ? AND is_active = true`
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Load roles and permissions
	if err := r.loadUserRoles(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at 
			  FROM users WHERE username = ? AND is_active = true`
	
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Load roles and permissions
	if err := r.loadUserRoles(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, username, email, password_hash, is_active, created_at, updated_at 
			  FROM users WHERE email = ? AND is_active = true`
	
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Load roles and permissions
	if err := r.loadUserRoles(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user information
func (r *UserRepository) UpdateUser(user *model.User) error {
	query := `UPDATE users SET username = ?, email = ?, password_hash = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ?`
	
	_, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

// DeleteUser soft deletes a user
func (r *UserRepository) DeleteUser(id int) error {
	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// AssignRoleToUser assigns a role to a user
func (r *UserRepository) AssignRoleToUser(userID, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?) 
			  ON DUPLICATE KEY UPDATE assigned_at = CURRENT_TIMESTAMP`
	
	_, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}
	return nil
}

// RemoveRoleFromUser removes a role from a user
func (r *UserRepository) RemoveRoleFromUser(userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = ? AND role_id = ?`
	_, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}
	return nil
}

// StoreRefreshToken stores a refresh token
func (r *UserRepository) StoreRefreshToken(token *model.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, token.UserID, token.TokenHash, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get token ID: %w", err)
	}
	
	token.ID = int(id)
	return nil
}

// GetRefreshToken retrieves a refresh token by hash
func (r *UserRepository) GetRefreshToken(tokenHash string) (*model.RefreshToken, error) {
	token := &model.RefreshToken{}
	query := `SELECT id, user_id, token_hash, expires_at, created_at 
			  FROM refresh_tokens WHERE token_hash = ? AND expires_at > NOW()`
	
	err := r.db.QueryRow(query, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	
	return token, nil
}

// DeleteRefreshToken deletes a refresh token
func (r *UserRepository) DeleteRefreshToken(tokenHash string) error {
	query := `DELETE FROM refresh_tokens WHERE token_hash = ?`
	_, err := r.db.Exec(query, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

// DeleteUserRefreshTokens deletes all refresh tokens for a user
func (r *UserRepository) DeleteUserRefreshTokens(userID int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user refresh tokens: %w", err)
	}
	return nil
}

// CleanupExpiredTokens removes expired refresh tokens
func (r *UserRepository) CleanupExpiredTokens() error {
	query := `DELETE FROM refresh_tokens WHERE expires_at <= NOW()`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	return nil
}

// loadUserRoles loads roles and permissions for a user
func (r *UserRepository) loadUserRoles(user *model.User) error {
	query := `SELECT r.id, r.name, r.description, r.created_at,
					 p.id, p.name, p.description, p.resource, p.action, p.created_at
			  FROM roles r
			  JOIN user_roles ur ON r.id = ur.role_id
			  LEFT JOIN role_permissions rp ON r.id = rp.role_id
			  LEFT JOIN permissions p ON rp.permission_id = p.id
			  WHERE ur.user_id = ?
			  ORDER BY r.id, p.id`

	rows, err := r.db.Query(query, user.ID)
	if err != nil {
		return fmt.Errorf("failed to load user roles: %w", err)
	}
	defer rows.Close()

	roleMap := make(map[int]*model.Role)
	
	for rows.Next() {
		var roleID, permID sql.NullInt64
		var roleName, roleDesc, permName, permDesc, resource, action sql.NullString
		var roleCreatedAt, permCreatedAt sql.NullTime

		err := rows.Scan(
			&roleID, &roleName, &roleDesc, &roleCreatedAt,
			&permID, &permName, &permDesc, &resource, &action, &permCreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning role/permission: %v", err)
			continue
		}

		if !roleID.Valid {
			continue
		}

		// Get or create role
		role, exists := roleMap[int(roleID.Int64)]
		if !exists {
			role = &model.Role{
				ID:          int(roleID.Int64),
				Name:        roleName.String,
				Description: roleDesc.String,
				CreatedAt:   roleCreatedAt.Time,
				Permissions: []model.Permission{},
			}
			roleMap[int(roleID.Int64)] = role
		}

		// Add permission if exists
		if permID.Valid {
			permission := model.Permission{
				ID:          int(permID.Int64),
				Name:        permName.String,
				Description: permDesc.String,
				Resource:    resource.String,
				Action:      action.String,
				CreatedAt:   permCreatedAt.Time,
			}
			role.Permissions = append(role.Permissions, permission)
		}
	}

	// Convert map to slice
	for _, role := range roleMap {
		user.Roles = append(user.Roles, *role)
	}

	return nil
}