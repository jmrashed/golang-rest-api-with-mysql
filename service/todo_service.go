package service

import (
	"fmt"

	"jmrashed/apps/userApp/model"
	"jmrashed/apps/userApp/repository"

	"github.com/go-playground/validator/v10"
)

type TodoService struct {
	todoRepo  *repository.TodoRepository
	validator *validator.Validate
}

func NewTodoService(todoRepo *repository.TodoRepository) *TodoService {
	return &TodoService{
		todoRepo:  todoRepo,
		validator: validator.New(),
	}
}

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(userID int, req model.CreateTodoRequest) (*model.Todo, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	todo := &model.Todo{
		UserID:    userID,
		Title:     req.Title,
		Content:   req.Content,
		Completed: false,
	}

	if err := s.todoRepo.CreateTodo(todo); err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return todo, nil
}

// GetTodoByID retrieves a todo by ID
func (s *TodoService) GetTodoByID(id int) (*model.Todo, error) {
	return s.todoRepo.GetTodoByID(id)
}

// GetUserTodos retrieves todos for a user with pagination
func (s *TodoService) GetUserTodos(userID int, req model.PaginationRequest) (*model.PaginatedResponse, error) {
	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	todos, total, err := s.todoRepo.GetTodosByUser(userID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	
	pagination := model.Pagination{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	return &model.PaginatedResponse{
		Data:       todos,
		Pagination: pagination,
	}, nil
}

// UpdateTodo updates a todo
func (s *TodoService) UpdateTodo(id, userID int, req model.UpdateTodoRequest) (*model.Todo, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Get existing todo
	todo, err := s.todoRepo.GetTodoByID(id)
	if err != nil {
		return nil, fmt.Errorf("todo not found: %w", err)
	}

	// Check ownership
	if todo.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	// Update fields
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Content != nil {
		todo.Content = *req.Content
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}

	if err := s.todoRepo.UpdateTodo(todo); err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return todo, nil
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(id, userID int) error {
	return s.todoRepo.DeleteTodo(id, userID)
}

// GetAllTodos retrieves all todos (admin only)
func (s *TodoService) GetAllTodos(req model.PaginationRequest) (*model.PaginatedResponse, error) {
	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	todos, total, err := s.todoRepo.GetAllTodos(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))
	
	pagination := model.Pagination{
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}

	return &model.PaginatedResponse{
		Data:       todos,
		Pagination: pagination,
	}, nil
}