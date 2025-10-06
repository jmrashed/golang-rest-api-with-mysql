package seeder

import (
	"database/sql"
	"log"

	"jmrashed/apps/userApp/auth"
	"jmrashed/apps/userApp/model"
)

type Seeder struct {
	db *sql.DB
}

func NewSeeder(db *sql.DB) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) Run() error {
	log.Println("Running database seeder...")

	if err := s.seedRoles(); err != nil {
		return err
	}

	if err := s.seedPermissions(); err != nil {
		return err
	}

	if err := s.seedRolePermissions(); err != nil {
		return err
	}

	if err := s.seedAdminUser(); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func (s *Seeder) seedRoles() error {
	roles := []model.Role{
		{ID: 1, Name: "admin", Description: "Administrator with full access"},
		{ID: 2, Name: "user", Description: "Regular user with limited access"},
		{ID: 3, Name: "moderator", Description: "Moderator with intermediate access"},
	}

	for _, role := range roles {
		_, err := s.db.Exec(`INSERT IGNORE INTO roles (id, name, description) VALUES (?, ?, ?)`,
			role.ID, role.Name, role.Description)
		if err != nil {
			return err
		}
	}

	log.Println("Roles seeded successfully")
	return nil
}

func (s *Seeder) seedPermissions() error {
	permissions := []model.Permission{
		{ID: 1, Name: "read_users", Description: "Read user information", Resource: "users", Action: "read"},
		{ID: 2, Name: "write_users", Description: "Create and update users", Resource: "users", Action: "write"},
		{ID: 3, Name: "delete_users", Description: "Delete users", Resource: "users", Action: "delete"},
		{ID: 4, Name: "read_todos", Description: "Read todos", Resource: "todos", Action: "read"},
		{ID: 5, Name: "write_todos", Description: "Create and update todos", Resource: "todos", Action: "write"},
		{ID: 6, Name: "delete_todos", Description: "Delete todos", Resource: "todos", Action: "delete"},
		{ID: 7, Name: "manage_roles", Description: "Manage user roles", Resource: "roles", Action: "manage"},
	}

	for _, perm := range permissions {
		_, err := s.db.Exec(`INSERT IGNORE INTO permissions (id, name, description, resource, action) VALUES (?, ?, ?, ?, ?)`,
			perm.ID, perm.Name, perm.Description, perm.Resource, perm.Action)
		if err != nil {
			return err
		}
	}

	log.Println("Permissions seeded successfully")
	return nil
}

func (s *Seeder) seedRolePermissions() error {
	rolePermissions := map[int][]int{
		1: {1, 2, 3, 4, 5, 6, 7}, // admin - all permissions
		2: {1, 4, 5},              // user - read users, read/write todos
		3: {1, 2, 4, 5, 6},        // moderator - users + todos management
	}

	for roleID, permIDs := range rolePermissions {
		for _, permID := range permIDs {
			_, err := s.db.Exec(`INSERT IGNORE INTO role_permissions (role_id, permission_id) VALUES (?, ?)`,
				roleID, permID)
			if err != nil {
				return err
			}
		}
	}

	log.Println("Role permissions seeded successfully")
	return nil
}

func (s *Seeder) seedAdminUser() error {
	// Check if admin user already exists
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		log.Println("Admin user already exists, skipping...")
		return nil
	}

	// Create admin user
	hashedPassword, err := auth.HashPassword("admin123")
	if err != nil {
		return err
	}

	result, err := s.db.Exec(`INSERT INTO users (username, email, password_hash, is_active) VALUES (?, ?, ?, ?)`,
		"admin", "admin@example.com", hashedPassword, true)
	if err != nil {
		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Assign admin role
	_, err = s.db.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`, userID, 1)
	if err != nil {
		return err
	}

	log.Printf("Admin user created successfully (ID: %d)", userID)
	return nil
}