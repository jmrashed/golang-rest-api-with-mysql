package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"jmrashed/apps/userApp/model"
)

type TodoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// CreateTodo creates a new todo
func (r *TodoRepository) CreateTodo(todo *model.Todo) error {
	query := `INSERT INTO todos (user_id, title, content, completed) VALUES (?, ?, ?, ?)`
	result, err := r.db.Exec(query, todo.UserID, todo.Title, todo.Content, todo.Completed)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get todo ID: %w", err)
	}
	
	todo.ID = int(id)
	return nil
}

// GetTodoByID retrieves a todo by ID
func (r *TodoRepository) GetTodoByID(id int) (*model.Todo, error) {
	todo := &model.Todo{}
	query := `SELECT id, user_id, title, content, completed, created_at, updated_at 
			  FROM todos WHERE id = ?`
	
	err := r.db.QueryRow(query, id).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Content,
		&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	
	return todo, nil
}

// GetTodosByUser retrieves todos for a user with pagination and filtering
func (r *TodoRepository) GetTodosByUser(userID int, req model.PaginationRequest) ([]model.Todo, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE user_id = ?"
	args := []interface{}{userID}
	
	// Add search filter
	if req.Search != "" {
		whereClause += " AND (title LIKE ? OR content LIKE ?)"
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}
	
	// Add completion filter
	if req.Filter == "completed" {
		whereClause += " AND completed = true"
	} else if req.Filter == "pending" {
		whereClause += " AND completed = false"
	}
	
	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM todos %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "ORDER BY created_at DESC"
	if req.Sort != "" {
		order := "ASC"
		if req.Order == "desc" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", req.Sort, order)
	}
	
	// Build main query with pagination
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, user_id, title, content, completed, created_at, updated_at 
		FROM todos %s %s LIMIT ? OFFSET ?
	`, whereClause, orderBy)
	
	args = append(args, req.Limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()
	
	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(
			&todo.ID, &todo.UserID, &todo.Title, &todo.Content,
			&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}
	
	return todos, total, nil
}

// UpdateTodo updates a todo
func (r *TodoRepository) UpdateTodo(todo *model.Todo) error {
	query := `UPDATE todos SET title = ?, content = ?, completed = ?, updated_at = CURRENT_TIMESTAMP 
			  WHERE id = ? AND user_id = ?`
	
	result, err := r.db.Exec(query, todo.Title, todo.Content, todo.Completed, todo.ID, todo.UserID)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("todo not found or access denied")
	}
	
	return nil
}

// DeleteTodo deletes a todo
func (r *TodoRepository) DeleteTodo(id, userID int) error {
	query := `DELETE FROM todos WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("todo not found or access denied")
	}
	
	return nil
}

// GetAllTodos retrieves all todos (admin only) with pagination
func (r *TodoRepository) GetAllTodos(req model.PaginationRequest) ([]model.Todo, int64, error) {
	// Build WHERE clause for search
	whereClause := ""
	args := []interface{}{}
	
	if req.Search != "" {
		whereClause = "WHERE title LIKE ? OR content LIKE ?"
		searchTerm := "%" + req.Search + "%"
		args = append(args, searchTerm, searchTerm)
	}
	
	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM todos %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}
	
	// Build ORDER BY clause
	orderBy := "ORDER BY created_at DESC"
	if req.Sort != "" {
		order := "ASC"
		if req.Order == "desc" {
			order = "DESC"
		}
		orderBy = fmt.Sprintf("ORDER BY %s %s", req.Sort, order)
	}
	
	// Build main query with pagination
	offset := (req.Page - 1) * req.Limit
	query := fmt.Sprintf(`
		SELECT id, user_id, title, content, completed, created_at, updated_at 
		FROM todos %s %s LIMIT ? OFFSET ?
	`, whereClause, orderBy)
	
	args = append(args, req.Limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()
	
	var todos []model.Todo
	for rows.Next() {
		var todo model.Todo
		err := rows.Scan(
			&todo.ID, &todo.UserID, &todo.Title, &todo.Content,
			&todo.Completed, &todo.CreatedAt, &todo.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}
	
	return todos, total, nil
}