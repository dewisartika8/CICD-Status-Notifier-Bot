package services_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/handler/project"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository/postgres"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/service"
	"github.com/dewisartika8/cicd-status-notifier-bot/pkg/logger"
	"github.com/dewisartika8/cicd-status-notifier-bot/tests/testutils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Constants for API endpoints and content types
const (
	APIProjectsEndpoint     = "/api/v1/projects"
	APIProjectByIDEndpoint  = "/api/v1/projects/%s"
	ContentTypeHeader       = "Content-Type"
	ApplicationJSONMimeType = "application/json"
)

type ProjectIntegrationTestSuite struct {
	suite.Suite
	app     *fiber.App
	cleanup func()
}

func (suite *ProjectIntegrationTestSuite) SetupSuite() {
	db := testutils.SetupTestDB(suite.T())
	suite.cleanup = func() {
		testutils.TeardownTestDB(suite.T(), db)
	}

	// Setup repositories
	projectRepo := postgres.NewProjectRepository(db)

	// Setup services
	projectService := service.NewProjectService(service.Dep{
		ProjectRepo: projectRepo,
	})

	// Setup logger
	testLogger := logger.NewLogger()

	// Setup handlers
	projectHandler := project.NewProjectHandler(project.ProjectHandlerDep{
		ProjectService: projectService,
		Logger:         testLogger,
	})

	// Setup Fiber app
	suite.app = fiber.New()
	api := suite.app.Group("/api/v1")
	projectHandler.RegisterRoutes(api)
}

func (suite *ProjectIntegrationTestSuite) TearDownSuite() {
	if suite.cleanup != nil {
		suite.cleanup()
	}
}

func (suite *ProjectIntegrationTestSuite) TestCreateProject() {
	// Test data
	createReq := dto.CreateProjectRequest{
		Name:          "Test Project",
		RepositoryURL: "https://github.com/test/repo",
		WebhookSecret: "test-secret-123",
	}

	reqBody, _ := json.Marshal(createReq)

	// Create request
	req := httptest.NewRequest("POST", APIProjectsEndpoint, bytes.NewReader(reqBody))
	req.Header.Set(ContentTypeHeader, ApplicationJSONMimeType)

	// Execute request
	resp, err := suite.app.Test(req)
	suite.NoError(err)

	// Assert response
	assert.Equal(suite.T(), 201, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	assert.Equal(suite.T(), "Project created successfully", response["message"])
	assert.NotNil(suite.T(), response["data"])

	// Verify project data
	projectData, ok := response["data"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), createReq.Name, projectData["name"])
	assert.Equal(suite.T(), createReq.RepositoryURL, projectData["repository_url"])
	assert.Equal(suite.T(), "active", projectData["status"])
}

func (suite *ProjectIntegrationTestSuite) TestCreateProjectValidationError() {
	// Test data with invalid URL
	createReq := dto.CreateProjectRequest{
		Name:          "Test Project",
		RepositoryURL: "invalid-url",
		WebhookSecret: "test-secret-123",
	}

	reqBody, _ := json.Marshal(createReq)

	// Create request
	req := httptest.NewRequest("POST", APIProjectsEndpoint, bytes.NewReader(reqBody))
	req.Header.Set(ContentTypeHeader, ApplicationJSONMimeType)

	// Execute request
	resp, err := suite.app.Test(req)
	suite.NoError(err)

	// Assert response
	assert.Equal(suite.T(), 400, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	assert.Equal(suite.T(), "Validation failed", response["error"])
	assert.NotNil(suite.T(), response["details"])
}

func (suite *ProjectIntegrationTestSuite) TestListProjects() {
	// First create a project
	createReq := dto.CreateProjectRequest{
		Name:          "List Test Project",
		RepositoryURL: "https://github.com/test/list-repo",
		WebhookSecret: "test-secret-list",
	}

	reqBody, _ := json.Marshal(createReq)
	createReqHttp := httptest.NewRequest("POST", APIProjectsEndpoint, bytes.NewReader(reqBody))
	createReqHttp.Header.Set(ContentTypeHeader, ApplicationJSONMimeType)

	_, err := suite.app.Test(createReqHttp)
	suite.NoError(err)

	// Now list projects
	req := httptest.NewRequest("GET", APIProjectsEndpoint, nil)

	// Execute request
	resp, err := suite.app.Test(req)
	suite.NoError(err)

	// Assert response
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	assert.Equal(suite.T(), "Projects retrieved successfully", response["message"])
	assert.NotNil(suite.T(), response["data"])

	// Verify list data
	listData, ok := response["data"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.NotNil(suite.T(), listData["projects"])
	assert.NotNil(suite.T(), listData["total"])
}

func (suite *ProjectIntegrationTestSuite) TestGetProject() {
	// First create a project
	createReq := dto.CreateProjectRequest{
		Name:          "Get Test Project",
		RepositoryURL: "https://github.com/test/get-repo",
		WebhookSecret: "test-secret-get",
	}

	reqBody, _ := json.Marshal(createReq)
	createReqHttp := httptest.NewRequest("POST", APIProjectsEndpoint, bytes.NewReader(reqBody))
	createReqHttp.Header.Set(ContentTypeHeader, ApplicationJSONMimeType)

	createResp, err := suite.app.Test(createReqHttp)
	suite.NoError(err)

	// Parse create response to get project ID
	var createResponse map[string]interface{}
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	suite.NoError(err)

	projectData := createResponse["data"].(map[string]interface{})
	projectID := projectData["id"].(string)

	// Now get the project
	req := httptest.NewRequest("GET", fmt.Sprintf(APIProjectByIDEndpoint, projectID), nil)

	// Execute request
	resp, err := suite.app.Test(req)
	suite.NoError(err)

	// Assert response
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	assert.Equal(suite.T(), "Project retrieved successfully", response["message"])
	assert.NotNil(suite.T(), response["data"])

	// Verify project data
	getData, ok := response["data"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), projectID, getData["id"])
	assert.Equal(suite.T(), createReq.Name, getData["name"])
}

func (suite *ProjectIntegrationTestSuite) TestGetProjectNotFound() {
	// Test with non-existent UUID
	req := httptest.NewRequest("GET", fmt.Sprintf(APIProjectByIDEndpoint, "550e8400-e29b-41d4-a716-446655440000"), nil)

	// Execute request
	resp, err := suite.app.Test(req)
	suite.NoError(err)

	// Assert response
	assert.Equal(suite.T(), 404, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	assert.Contains(suite.T(), response["error"], "not found")
}

func (suite *ProjectIntegrationTestSuite) TestUpdateProject() {
	// First create a project
	createReq := dto.CreateProjectRequest{
		Name:          "Update Test Project",
		RepositoryURL: "https://github.com/test/update-repo",
		WebhookSecret: "test-secret-update",
	}

	reqBody, _ := json.Marshal(createReq)
	createReqHttp := httptest.NewRequest("POST", APIProjectsEndpoint, bytes.NewReader(reqBody))
	createReqHttp.Header.Set(ContentTypeHeader, ApplicationJSONMimeType)

	createResp, err := suite.app.Test(createReqHttp)
	suite.NoError(err)

	// Parse create response to get project ID
	var createResponse map[string]interface{}
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	suite.NoError(err)

	projectData := createResponse["data"].(map[string]interface{})
	projectID := projectData["id"].(string)

	// Prepare update data
	newName := "Updated Project Name"
	updateReq := dto.UpdateProjectRequest{
		Name: &newName,
	}

	updateBody, _ := json.Marshal(updateReq)

	// Update the project
	req := httptest.NewRequest("PUT", fmt.Sprintf(APIProjectByIDEndpoint, projectID), bytes.NewReader(updateBody))
	req.Header.Set(ContentTypeHeader, ApplicationJSONMimeType)

	// Execute request
	resp, err := suite.app.Test(req)
	suite.NoError(err)

	// Assert response
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	assert.Equal(suite.T(), "Project updated successfully", response["message"])
	assert.NotNil(suite.T(), response["data"])

	// Verify updated data
	updatedData, ok := response["data"].(map[string]interface{})
	assert.True(suite.T(), ok)
	assert.Equal(suite.T(), newName, updatedData["name"])
	assert.Equal(suite.T(), createReq.RepositoryURL, updatedData["repository_url"])
}

func (suite *ProjectIntegrationTestSuite) TestDeleteProject() {
	// First create a project
	createReq := dto.CreateProjectRequest{
		Name:          "Delete Test Project",
		RepositoryURL: "https://github.com/test/delete-repo",
		WebhookSecret: "test-secret-delete",
	}

	reqBody, _ := json.Marshal(createReq)
	createReqHttp := httptest.NewRequest("POST", APIProjectsEndpoint, bytes.NewReader(reqBody))
	createReqHttp.Header.Set(ContentTypeHeader, ApplicationJSONMimeType)

	createResp, err := suite.app.Test(createReqHttp)
	suite.NoError(err)

	// Parse create response to get project ID
	var createResponse map[string]interface{}
	err = json.NewDecoder(createResp.Body).Decode(&createResponse)
	suite.NoError(err)

	projectData := createResponse["data"].(map[string]interface{})
	projectID := projectData["id"].(string)

	// Delete the project
	req := httptest.NewRequest("DELETE", fmt.Sprintf(APIProjectByIDEndpoint, projectID), nil)

	// Execute request
	resp, err := suite.app.Test(req)
	suite.NoError(err)

	// Assert response
	assert.Equal(suite.T(), 200, resp.StatusCode)

	// Parse response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err)

	assert.Equal(suite.T(), "Project deleted successfully", response["message"])

	// Verify project is deleted by trying to get it
	getReq := httptest.NewRequest("GET", fmt.Sprintf(APIProjectByIDEndpoint, projectID), nil)
	getResp, err := suite.app.Test(getReq)
	suite.NoError(err)
	assert.Equal(suite.T(), 404, getResp.StatusCode)
}

func TestProjectIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ProjectIntegrationTestSuite))
}
