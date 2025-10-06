package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"jmrashed/apps/userApp/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	server      *httptest.Server
	client      *http.Client
	accessToken string
}

func (suite *E2ETestSuite) SetupSuite() {
	// Setup test server
	// Note: In a real implementation, you would setup your actual server here
	suite.client = &http.Client{}
}

func (suite *E2ETestSuite) TearDownSuite() {
	if suite.server != nil {
		suite.server.Close()
	}
}

func (suite *E2ETestSuite) TestUserRegistrationAndLoginFlow() {
	// Test user registration
	registerReq := model.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	registerBody, _ := json.Marshal(registerReq)
	resp, err := suite.client.Post(
		fmt.Sprintf("%s/api/v1/register", suite.server.URL),
		"application/json",
		bytes.NewBuffer(registerBody),
	)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var registerResponse model.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&registerResponse)
	resp.Body.Close()

	// Extract tokens from registration response
	authData := registerResponse.Data.(map[string]interface{})
	suite.accessToken = authData["access_token"].(string)

	// Test user login
	loginReq := model.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginBody, _ := json.Marshal(loginReq)
	resp, err = suite.client.Post(
		fmt.Sprintf("%s/api/v1/login", suite.server.URL),
		"application/json",
		bytes.NewBuffer(loginBody),
	)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestTodoCRUDFlow() {
	// Create todo
	createReq := model.CreateTodoRequest{
		Title:   "Test Todo",
		Content: "This is a test todo",
	}

	createBody, _ := json.Marshal(createReq)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/todos", suite.server.URL), bytes.NewBuffer(createBody))
	req.Header.Set("Authorization", "Bearer "+suite.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	var createResponse model.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&createResponse)
	resp.Body.Close()

	todoData := createResponse.Data.(map[string]interface{})
	todoID := int(todoData["id"].(float64))

	// Get todo
	req, _ = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/todos/%d", suite.server.URL, todoID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.accessToken)

	resp, err = suite.client.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Update todo
	completed := true
	updateReq := model.UpdateTodoRequest{
		Title:     &createReq.Title,
		Content:   &createReq.Content,
		Completed: &completed,
	}

	updateBody, _ := json.Marshal(updateReq)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/todos/%d", suite.server.URL, todoID), bytes.NewBuffer(updateBody))
	req.Header.Set("Authorization", "Bearer "+suite.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = suite.client.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Delete todo
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/todos/%d", suite.server.URL, todoID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.accessToken)

	resp, err = suite.client.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func (suite *E2ETestSuite) TestPaginationAndFiltering() {
	// Create multiple todos
	for i := 0; i < 5; i++ {
		createReq := model.CreateTodoRequest{
			Title:   fmt.Sprintf("Todo %d", i+1),
			Content: fmt.Sprintf("Content for todo %d", i+1),
		}

		createBody, _ := json.Marshal(createReq)
		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/todos", suite.server.URL), bytes.NewBuffer(createBody))
		req.Header.Set("Authorization", "Bearer "+suite.accessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := suite.client.Do(req)
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
		resp.Body.Close()
	}

	// Test pagination
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/todos?page=1&limit=3", suite.server.URL), nil)
	req.Header.Set("Authorization", "Bearer "+suite.accessToken)

	resp, err := suite.client.Do(req)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	var paginatedResponse model.SuccessResponse
	json.NewDecoder(resp.Body).Decode(&paginatedResponse)
	resp.Body.Close()

	responseData := paginatedResponse.Data.(map[string]interface{})
	pagination := responseData["pagination"].(map[string]interface{})
	
	assert.Equal(suite.T(), float64(1), pagination["page"])
	assert.Equal(suite.T(), float64(3), pagination["limit"])
}

func (suite *E2ETestSuite) TestRateLimiting() {
	// Make multiple rapid requests to test rate limiting
	for i := 0; i < 70; i++ { // Exceed the 60 requests per minute limit
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/health", suite.server.URL), nil)
		resp, err := suite.client.Do(req)
		
		assert.NoError(suite.T(), err)
		
		if i >= 60 {
			// Should be rate limited
			assert.Equal(suite.T(), http.StatusTooManyRequests, resp.StatusCode)
		}
		
		resp.Body.Close()
	}
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}