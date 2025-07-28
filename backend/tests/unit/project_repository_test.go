package database_test

import (
	"context"
	"testing"

	"github.com/dewisartika8/cicd-status-notifier-bot/internal/adapters/database"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/entities"
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/domain/ports"
	"github.com/dewisartika8/cicd-status-notifier-bot/tests/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := database.NewProjectRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		project       *entities.Project
		expectedError error
	}{
		{
			name:          "should create project successfully",
			project:       testutils.CreateTestProject("test-project", "https://github.com/user/repo"),
			expectedError: nil,
		},
		{
			name:          "should fail with duplicate name",
			project:       testutils.CreateTestProject("test-project", "https://github.com/user/another-repo"), // Same name as above
			expectedError: ports.ErrProjectAlreadyExists,
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
				assert.NotEqual(t, uuid.Nil, tt.project.ID)
				assert.NotZero(t, tt.project.CreatedAt)
				assert.NotZero(t, tt.project.UpdatedAt)
			}
		})
	}
}

func TestProjectRepository_GetByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := database.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("test-project", "https://github.com/user/repo")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	tests := []struct {
		name          string
		projectID     uuid.UUID
		expectedError error
	}{
		{
			name:          "should get project by ID",
			projectID:     project.ID,
			expectedError: nil,
		},
		{
			name:          "should return error for non-existent project",
			projectID:     uuid.New(),
			expectedError: ports.ErrProjectNotFound,
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
				assert.Equal(t, project.ID, foundProject.ID)
				assert.Equal(t, project.Name, foundProject.Name)
				assert.Equal(t, project.RepositoryURL, foundProject.RepositoryURL)
				assert.Equal(t, project.IsActive, foundProject.IsActive)
			}
		})
	}
}

func TestProjectRepository_GetByName(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := database.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("unique-project-name", "https://github.com/user/repo")
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
			expectedError: ports.ErrProjectNotFound,
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
				assert.Equal(t, project.ID, foundProject.ID)
				assert.Equal(t, project.Name, foundProject.Name)
			}
		})
	}
}

func TestProjectRepository_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := database.NewProjectRepository(db)
	ctx := context.Background()

	// Create multiple test projects
	project1 := testutils.CreateTestProject("project-1", "https://github.com/user/repo1")
	project2 := testutils.CreateTestProject("project-2", "https://github.com/user/repo2")

	err := repo.Create(ctx, project1)
	require.NoError(t, err)
	err = repo.Create(ctx, project2)
	require.NoError(t, err)

	projects, err := repo.List(ctx)

	require.NoError(t, err)
	assert.Len(t, projects, 2)

	// Check if both projects are in the list
	projectIDs := make(map[uuid.UUID]bool)
	for _, p := range projects {
		projectIDs[p.ID] = true
	}

	assert.True(t, projectIDs[project1.ID])
	assert.True(t, projectIDs[project2.ID])
}

func TestProjectRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := database.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("original-name", "https://github.com/user/repo")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	// Update project
	telegramChatID := int64(123456789)
	err = project.Update("updated-name", "https://github.com/updated/repo", &telegramChatID, false)
	require.NoError(t, err)

	err = repo.Update(ctx, project)
	require.NoError(t, err)

	// Verify update
	updatedProject, err := repo.GetByID(ctx, project.ID)
	require.NoError(t, err)
	assert.Equal(t, "updated-name", updatedProject.Name)
	assert.Equal(t, "https://github.com/updated/repo", updatedProject.RepositoryURL)
	assert.False(t, updatedProject.IsActive)
	assert.Equal(t, &telegramChatID, updatedProject.TelegramChatID)
}

func TestProjectRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	repo := database.NewProjectRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("to-be-deleted", "https://github.com/user/repo")
	err := repo.Create(ctx, project)
	require.NoError(t, err)

	tests := []struct {
		name          string
		projectID     uuid.UUID
		expectedError error
	}{
		{
			name:          "should delete existing project",
			projectID:     project.ID,
			expectedError: nil,
		},
		{
			name:          "should return error for non-existent project",
			projectID:     uuid.New(),
			expectedError: ports.ErrProjectNotFound,
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
				assert.Equal(t, ports.ErrProjectNotFound, err)
			}
		})
	}
}

func TestProjectRepository_GetWithBuildEvents(t *testing.T) {
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(t, db)

	projectRepo := database.NewProjectRepository(db)
	buildEventRepo := database.NewBuildEventRepository(db)
	ctx := context.Background()

	// Create test project
	project := testutils.CreateTestProject("project-with-builds", "https://github.com/user/repo")
	err := projectRepo.Create(ctx, project)
	require.NoError(t, err)

	// Create test build events
	buildEvent1 := testutils.CreateTestBuildEvent(project.ID, entities.EventTypeBuildSuccess, entities.BuildStatusSuccess, "main")
	buildEvent2 := testutils.CreateTestBuildEvent(project.ID, entities.EventTypeBuildFailed, entities.BuildStatusFailed, "develop")

	err = buildEventRepo.Create(ctx, buildEvent1)
	require.NoError(t, err)
	err = buildEventRepo.Create(ctx, buildEvent2)
	require.NoError(t, err)

	// Test getting project with build events
	projectWithEvents, buildEvents, err := projectRepo.GetWithBuildEvents(ctx, project.ID, 10)
	require.NoError(t, err)
	assert.Equal(t, project.ID, projectWithEvents.ID)
	assert.Len(t, buildEvents, 2)

	// Test with non-existent project
	_, _, err = projectRepo.GetWithBuildEvents(ctx, uuid.New(), 10)
	assert.Equal(t, ports.ErrProjectNotFound, err)
}
