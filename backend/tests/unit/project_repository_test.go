package repository_test

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapter/repository"
	projectdomain "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
	projectdto "github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/dto"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/shared/domain/value_objects"
	"github.com/dewisartika8/cicd-status-notifier-bot/tests/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testProjectName = "test-project"
	testRepoURL     = "https://github.com/user/repo"
)

func TestProjectRepositoryCreate(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := repository.NewProjectRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		project       *projectdomain.Project
		expectedError error
	}{
		{
			name:          "should create project successfully",
			project:       testutils.CreateTestProject(testProjectName, testRepoURL),
			expectedError: nil,
		},
		{
			name:          "should fail with duplicate name",
			project:       testutils.CreateTestProject(testProjectName, "https://github.com/user/another-repo"), // Same name as above
			expectedError: projectdomain.ErrProjectAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.project)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.False(t, tt.project.ID().IsNil())
				assert.False(t, tt.project.CreatedAt().IsZero())
				assert.False(t, tt.project.UpdatedAt().IsZero())
			}
		})
	}
}

func TestProjectRepositoryGetByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := repository.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject(testProjectName, testRepoURL)
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	tests := []struct {
		name          string
		projectID     value_objects.ID
		expectedError error
	}{
		{
			name:          "should get project by ID",
			projectID:     project.ID(),
			expectedError: nil,
		},
		{
			name:          "should return error for non-existent project",
			projectID:     value_objects.NewID(),
			expectedError: projectdomain.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundProject, err := repo.GetByID(ctx, tt.projectID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, foundProject)
			} else {
				require.NoError(t, err)
				require.NotNil(t, foundProject)
				assert.Equal(t, project.ID(), foundProject.ID())
				assert.Equal(t, project.Name(), foundProject.Name())
				assert.Equal(t, project.RepositoryURL(), foundProject.RepositoryURL())
				assert.Equal(t, project.IsActive(), foundProject.IsActive())
			}
		})
	}
}

func TestProjectRepositoryGetByName(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := repository.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("unique-project-name", testRepoURL)
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	tests := []struct {
		name          string
		projectName   string
		expectedError error
	}{
		{
			name:          "should get project by name",
			projectName:   "unique-project-name",
			expectedError: nil,
		},
		{
			name:          "should return error for non-existent project",
			projectName:   "non-existent-project",
			expectedError: projectdomain.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			foundProject, err := repo.GetByName(ctx, tt.projectName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, foundProject)
			} else {
				require.NoError(t, err)
				require.NotNil(t, foundProject)
				assert.Equal(t, project.ID(), foundProject.ID())
				assert.Equal(t, project.Name(), foundProject.Name())
			}
		})
	}
}

func TestProjectRepositoryList(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := repository.NewProjectRepository(db)
	ctx := context.Background()

	// Create multiple test projects
	project1 := testutils.CreateTestProject("project-1", "https://github.com/user/repo1")
	project2 := testutils.CreateTestProject("project-2", "https://github.com/user/repo2")

	err := repo.Create(ctx, project1)
	require.NoError(t, err)
	err = repo.Create(ctx, project2)
	require.NoError(t, err)

	projects, err := repo.List(ctx, projectdto.ListProjectFilters{})

	require.NoError(t, err)
	assert.Len(t, projects, 2)

	// Check if both projects are in the list
	projectIDs := make(map[string]bool)
	for _, p := range projects {
		projectIDs[p.ID().String()] = true
	}

	assert.True(t, projectIDs[project1.ID().String()])
	assert.True(t, projectIDs[project2.ID().String()])
}

func TestProjectRepositoryUpdate(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := repository.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("original-name", testRepoURL)
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// Update project
	telegramChatID := int64(123456789)
	err = project.UpdateName("updated-name")
	require.NoError(t, err)
	err = project.UpdateRepositoryURL("https://github.com/updated/repo")
	require.NoError(t, err)
	err = project.UpdateTelegramChatID(&telegramChatID)
	require.NoError(t, err)
	project.Deactivate()

	err = repo.Update(ctx, project)
	require.NoError(t, err)

	// Verify update
	updatedProject, err := repo.GetByID(ctx, project.ID())
	require.NoError(t, err)
	assert.Equal(t, "updated-name", updatedProject.Name())
	assert.Equal(t, "https://github.com/updated/repo", updatedProject.RepositoryURL())
	assert.False(t, updatedProject.IsActive())
	assert.Equal(t, &telegramChatID, updatedProject.TelegramChatID())
}

func TestProjectRepositoryDelete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := repository.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("to-be-deleted", testRepoURL)
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	tests := []struct {
		name          string
		projectID     value_objects.ID
		expectedError error
	}{
		{
			name:          "should delete existing project",
			projectID:     project.ID(),
			expectedError: nil,
		},
		{
			name:          "should return error for non-existent project",
			projectID:     value_objects.NewID(),
			expectedError: projectdomain.ErrProjectNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(ctx, tt.projectID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)

				// Verify project is deleted
				_, err := repo.GetByID(ctx, tt.projectID)
				assert.Equal(t, projectdomain.ErrProjectNotFound, err)
			}
		})
	}
}

func TestProjectRepositoryGetActiveProjects(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := repository.NewProjectRepository(db)
	ctx := context.Background()

	// Create test projects
	project1 := testutils.CreateTestProject("active-project", "https://github.com/user/repo1")
	project2 := testutils.CreateTestProject("inactive-project", testRepoURL)

	err := repo.Create(ctx, project1)
	require.NoError(t, err)
	err = repo.Create(ctx, project2)
	require.NoError(t, err)

	// Deactivate second project
	project2.Deactivate()
	err = repo.Update(ctx, project2)
	require.NoError(t, err)

	// Test getting active projects
	activeProjects, err := repo.GetActiveProjects(ctx)
	require.NoError(t, err)
	assert.Len(t, activeProjects, 1)
	assert.Equal(t, project1.ID(), activeProjects[0].ID())
	assert.True(t, activeProjects[0].IsActive())
}
