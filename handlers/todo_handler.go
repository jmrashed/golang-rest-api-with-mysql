package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"jmrashed/apps/userApp/middleware"
	"jmrashed/apps/userApp/model"
	"jmrashed/apps/userApp/service"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	todoService *service.TodoService
}

func NewTodoHandler(todoService *service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// CreateTodo creates a new todo
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	var req model.CreateTodoRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	todo, err := h.todoService.CreateTodo(claims.UserID, req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusCreated, "Todo created successfully", todo)
}

// GetTodo retrieves a todo by ID
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := h.todoService.GetTodoByID(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Todo not found")
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Todo retrieved successfully", todo)
}

// GetUserTodos retrieves todos for current user
func (h *TodoHandler) GetUserTodos(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	// Parse query parameters
	req := model.PaginationRequest{
		Page:   1,
		Limit:  10,
		Sort:   "created_at",
		Order:  "desc",
		Filter: "all",
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			req.Limit = l
		}
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		req.Sort = sort
	}

	if order := r.URL.Query().Get("order"); order != "" {
		req.Order = order
	}

	if search := r.URL.Query().Get("search"); search != "" {
		req.Search = search
	}

	if filter := r.URL.Query().Get("filter"); filter != "" {
		req.Filter = filter
	}

	result, err := h.todoService.GetUserTodos(claims.UserID, req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Todos retrieved successfully", result)
}

// UpdateTodo updates a todo
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req model.UpdateTodoRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	todo, err := h.todoService.UpdateTodo(id, claims.UserID, req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Todo updated successfully", todo)
}

// DeleteTodo deletes a todo
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		writeErrorResponse(w, http.StatusUnauthorized, "User context not found")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	if err := h.todoService.DeleteTodo(id, claims.UserID); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "Todo deleted successfully", nil)
}

// GetAllTodos retrieves all todos (admin only)
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	req := model.PaginationRequest{
		Page:  1,
		Limit: 10,
		Sort:  "created_at",
		Order: "desc",
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			req.Limit = l
		}
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		req.Sort = sort
	}

	if order := r.URL.Query().Get("order"); order != "" {
		req.Order = order
	}

	if search := r.URL.Query().Get("search"); search != "" {
		req.Search = search
	}

	result, err := h.todoService.GetAllTodos(req)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeSuccessResponse(w, http.StatusOK, "All todos retrieved successfully", result)
}