package fixtures

import (
	"github.com/dewisartika8/cicd-status-notifier-bot/internal/core/project/domain"
)

const (
	// TestWebhookSecret is a common webhook secret for test projects
	TestWebhookSecret = "test-webhook-secret"
)

// ProjectFixtures provides common test data for projects
type ProjectFixtures struct{}

// NewProjectFixtures creates a new instance of ProjectFixtures
func NewProjectFixtures() *ProjectFixtures {
	return &ProjectFixtures{}
}

// CreateActiveProject creates a test project with active status
func (f *ProjectFixtures) CreateActiveProject(name, repoURL string, chatID *int64) (*domain.Project, error) {
	project, err := domain.NewProject(name, repoURL, TestWebhookSecret, chatID)
	if err != nil {
		return nil, err
	}

	// Project is active by default
	return project, nil
}

// CreateInactiveProject creates a test project with inactive status
func (f *ProjectFixtures) CreateInactiveProject(name, repoURL string, chatID *int64) (*domain.Project, error) {
	project, err := domain.NewProject(name, repoURL, TestWebhookSecret, chatID)
	if err != nil {
		return nil, err
	}

	project.Deactivate()
	return project, nil
}

// CreateArchivedProject creates a test project with archived status
func (f *ProjectFixtures) CreateArchivedProject(name, repoURL string, chatID *int64) (*domain.Project, error) {
	project, err := domain.NewProject(name, repoURL, TestWebhookSecret, chatID)
	if err != nil {
		return nil, err
	}

	project.Archive()
	return project, nil
}

// CommonChatIDs provides common chat IDs for testing
func (f *ProjectFixtures) CommonChatIDs() map[string]*int64 {
	chatID1 := int64(123456789)
	chatID2 := int64(987654321)
	chatID3 := int64(555666777)

	return map[string]*int64{
		"main":    &chatID1,
		"dev":     &chatID2,
		"staging": &chatID3,
		"no_chat": nil,
	}
}

// SampleProjects creates a set of sample projects for testing
func (f *ProjectFixtures) SampleProjects() ([]*domain.Project, error) {
	chatIDs := f.CommonChatIDs()
	var projects []*domain.Project

	// Active projects
	activeProject1, err := f.CreateActiveProject("api-service", "https://github.com/company/api-service", chatIDs["main"])
	if err != nil {
		return nil, err
	}

	activeProject2, err := f.CreateActiveProject("web-frontend", "https://github.com/company/web-frontend", chatIDs["dev"])
	if err != nil {
		return nil, err
	}

	// Inactive project
	inactiveProject, err := f.CreateInactiveProject("legacy-service", "https://github.com/company/legacy-service", chatIDs["staging"])
	if err != nil {
		return nil, err
	}

	// Archived project
	archivedProject, err := f.CreateArchivedProject("old-api", "https://github.com/company/old-api", nil)
	if err != nil {
		return nil, err
	}

	projects = append(projects, activeProject1, activeProject2, inactiveProject, archivedProject)
	return projects, nil
}
